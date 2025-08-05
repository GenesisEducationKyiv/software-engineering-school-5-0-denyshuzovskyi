//go:build integration

package test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	v1 "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-proto/gen/go/notification/v1"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/cache"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/client/weatherapi"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/client/weatherstack"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/database"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/lib/httputil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/lib/logger/noophandler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/lib/testutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/metrics"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/model"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/repository/postgresql"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/server/handler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/service/notification"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/service/subscription"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/service/weather"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/service/weatherprovider"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/service/weatherupd"
	validators "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/validator"
	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prometheus/client_golang/prometheus"
	ptestutil "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
)

type TestEnv struct {
	DB      *sql.DB
	Log     *slog.Logger
	Cleanup func()
}

func SetUpTestEnv(t *testing.T) *TestEnv {
	t.Helper()

	if testing.Short() {
		t.Skip()
	}

	log := slog.New(noophandler.NewNoOpHandler())
	ctx := t.Context()

	ctr, err := postgres.Run(ctx, "postgres:17-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	require.NoError(t, err)

	conn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	t.Log(conn)

	db, err := database.InitDB(ctx, conn, log)
	require.NoError(t, err)

	err = database.RunMigrations(db, ".", log)
	require.NoError(t, err)
	t.Log("migration completed successfully")

	return &TestEnv{
		DB:  db,
		Log: log,
		Cleanup: func() {
			require.NoError(t, db.Close())
			require.NoError(t, ctr.Terminate(ctx))
		},
	}
}

type TestCache struct {
	RedisCache *cache.RedisCache
	Cleanup    func()
}

func SetUpCache(t *testing.T) *TestCache {
	t.Helper()

	redisContainer, err := tcredis.Run(t.Context(), "redis:8-alpine")
	require.NoError(t, err)
	endpoint, err := redisContainer.Endpoint(t.Context(), "")
	require.NoError(t, err)

	redisClient := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	return &TestCache{
		RedisCache: cache.NewRedisCache(redisClient),
		Cleanup: func() {
			require.NoError(t, redisContainer.Terminate(t.Context()))
		},
	}
}

func TestGetWeatherIT(t *testing.T) {
	env := SetUpTestEnv(t)
	defer env.Cleanup()
	tCache := SetUpCache(t)
	defer tCache.Cleanup()

	// Client-interceptor for weatherapiClient
	waClient := httputil.NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
			Header:     make(http.Header),
		}, nil
	})
	weatherapiClient := weatherapi.NewClient("https://api.weatherapi.com/v1", "key", waClient, env.Log)
	weatherapiProvider := weatherprovider.NewWeatherapiProvider(weatherapiClient)

	// Client-interceptor for weatherstack
	currentWeatherData, err := os.ReadFile("./test_data/weatherstack_success_resp.json")
	require.NoError(t, err)
	wsClient := httputil.NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(currentWeatherData)),
			Header:     make(http.Header),
		}, nil
	})
	weatherstackClient := weatherstack.NewClient("https://api.weatherstack.com/current", "key", wsClient, env.Log)
	weatherstackProvider := weatherprovider.NewWeatherstackProvider(weatherstackClient)

	chainWeatherProvider := weatherprovider.NewChainWeatherProvider(env.Log, weatherapiProvider, weatherstackProvider)
	weatherCache := cache.NewJSONCache[dto.WeatherWithLocationDTO](tCache.RedisCache, 5*time.Minute)

	hitc := prometheus.NewCounter(prometheus.CounterOpts{Name: "test_hit", Help: ""})
	missc := prometheus.NewCounter(prometheus.CounterOpts{Name: "test_miss", Help: ""})
	errc := prometheus.NewCounter(prometheus.CounterOpts{Name: "test_err", Help: ""})
	cacheMetrics := metrics.NewPrometheusCacheMetrics(
		metrics.WithCacheHitsCounter(hitc),
		metrics.WithCacheMissesCounter(missc),
		metrics.WithCacheErrorsCounter(errc),
	)

	cachingWeatherProvider := weatherprovider.NewCachingWeatherProvider(weatherCache, chainWeatherProvider, cacheMetrics, env.Log)

	validate := validator.New()
	weatherService := weather.NewWeatherService(cachingWeatherProvider, env.Log)
	locationValidator := validators.NewLocationValidator(validate)
	weatherHandler := handler.NewWeatherHandler(weatherService, locationValidator, env.Log)

	city := "Kyiv"

	u := &url.URL{Path: "/weather"}
	q := u.Query()
	q.Set("city", city)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	require.NoError(t, err)

	// Expected
	expectedTemp := float32(16)
	expectedHum := float32(67)
	expectedDesc := "Clear "
	delta := 0.01

	// First req
	rr := httptest.NewRecorder()
	weatherHandler.GetCurrentWeather(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	assertWeatherResponse(t, rr.Body, expectedTemp, expectedHum, expectedDesc, delta)
	assertCacheMetrics(t, hitc, missc, errc, 0, 1, 0, delta)

	// Second req
	rr = httptest.NewRecorder()
	weatherHandler.GetCurrentWeather(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	assertWeatherResponse(t, rr.Body, expectedTemp, expectedHum, expectedDesc, delta)
	assertCacheMetrics(t, hitc, missc, errc, 1, 1, 0, delta)
}

func TestFullCycleIT(t *testing.T) {
	env := SetUpTestEnv(t)
	defer env.Cleanup()

	// Client-interceptor for weatherapiClient
	currentWeatherData, err := os.ReadFile("./test_data/weatherapi_success_resp.json")
	require.NoError(t, err)
	waClient := httputil.NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(currentWeatherData)),
			Header:     make(http.Header),
		}, nil
	})
	weatherapiClient := weatherapi.NewClient("https://api.weatherapi.com/v1", "key", waClient, env.Log)
	weatherapiProvider := weatherprovider.NewWeatherapiProvider(weatherapiClient)

	subsNotificationSenderMock := subscription.NewMockNotificationSender(t)
	notificationSenderMock := notification.NewMockNotificationSender(t)

	// Repositories, Services, Handlers
	subscriberRepository := postgresql.NewSubscriberRepository()
	subscriptionRepository := postgresql.NewSubscriptionRepository()
	tokenRepository := postgresql.NewTokenRepository()

	subscriptionService := subscription.NewSubscriptionService(
		env.DB,
		weatherapiProvider,
		subscriberRepository,
		subscriptionRepository,
		tokenRepository,
		subsNotificationSenderMock,
		env.Log)
	notificationService := notification.NewNotificationService(notificationSenderMock)
	weatherService := weather.NewWeatherService(weatherapiProvider, env.Log)
	weatherUpdateSendingService := weatherupd.NewWeatherUpdateSendingService(subscriptionService, weatherService, notificationService, env.Log)

	validate := validator.New()
	subscriptionValidator := validators.NewSubscriptionValidator(validate)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, subscriptionValidator, env.Log)

	// Multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("POST /subscribe", subscriptionHandler.Subscribe)
	mux.HandleFunc("GET /confirm/{token}", subscriptionHandler.Confirm)
	mux.HandleFunc("GET /unsubscribe/{token}", subscriptionHandler.Unsubscribe)

	// Subscribe
	subscriberEmail := "test@example.com"
	subsForm := url.Values{}
	subsForm.Set("email", subscriberEmail)
	subsForm.Set("city", "Kyiv")
	subsForm.Set("frequency", "daily")

	req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(subsForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	var capturedSendConfirmRequest *v1.SendConfirmationRequest
	subsNotificationSenderMock.EXPECT().
		SendConfirmation(mock.Anything, mock.AnythingOfType("*v1.SendConfirmationRequest")).
		Run(func(ctx context.Context, req *v1.SendConfirmationRequest, opts ...grpc.CallOption) {
			capturedSendConfirmRequest = req
		}).
		Return(nil, nil).Once()

	mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	lastToken := capturedSendConfirmRequest.GetNotificationWithToken().GetToken()
	require.NotEmpty(t, lastToken)

	// Confirm
	var capturedSendConfirmSuccessRequest *v1.SendConfirmationSuccessRequest
	subsNotificationSenderMock.EXPECT().
		SendConfirmationSuccess(mock.Anything, mock.AnythingOfType("*v1.SendConfirmationSuccessRequest")).
		Run(func(ctx context.Context, req *v1.SendConfirmationSuccessRequest, opts ...grpc.CallOption) {
			capturedSendConfirmSuccessRequest = req
		}).
		Return(nil, nil).Once()

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/confirm/%s", lastToken), nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	lastToken = capturedSendConfirmSuccessRequest.GetNotificationWithToken().GetToken()
	require.NotEmpty(t, lastToken)

	// Imitate notification job trigger
	var capturedSendWeatherUpdateRequest *v1.SendWeatherUpdateRequest
	notificationSenderMock.EXPECT().
		SendWeatherUpdate(mock.Anything, mock.AnythingOfType("*v1.SendWeatherUpdateRequest")).
		Run(func(ctx context.Context, req *v1.SendWeatherUpdateRequest, opts ...grpc.CallOption) {
			capturedSendWeatherUpdateRequest = req
		}).
		Return(nil, nil).
		Once()
	weatherUpdateSendingService.SendWeatherUpdates(t.Context(), model.Frequency_Daily)
	require.Equal(t, subscriberEmail, capturedSendWeatherUpdateRequest.GetWeatherUpdateNotification().GetNotificationWithToken().GetNotification().GetTo())
	lastToken = capturedSendWeatherUpdateRequest.GetWeatherUpdateNotification().GetNotificationWithToken().GetToken()
	require.NotEmpty(t, lastToken)

	// Unsubscribe
	var sendUnsubSuccessReq *v1.SendUnsubscribeSuccessRequest
	subsNotificationSenderMock.EXPECT().
		SendUnsubscribeSuccess(mock.Anything, mock.AnythingOfType("*v1.SendUnsubscribeSuccessRequest")).
		Run(func(ctx context.Context, req *v1.SendUnsubscribeSuccessRequest, opts ...grpc.CallOption) {
			sendUnsubSuccessReq = req
		}).
		Return(nil, nil).
		Once()
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/unsubscribe/%s", lastToken), nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, subscriberEmail, sendUnsubSuccessReq.GetNotification().GetTo())
}

func assertWeatherResponse(
	t *testing.T,
	body io.Reader,
	expectedTemp,
	expectedHum float32,
	expectedDesc string,
	delta float64,
) {
	t.Helper()
	actualWeatherDto, err := testutil.UnmarshalJSONFromReader[dto.WeatherDTO](body)
	require.NoError(t, err)

	require.InDelta(t, expectedTemp, actualWeatherDto.Temperature, delta)
	require.InDelta(t, expectedHum, actualWeatherDto.Humidity, delta)
	require.Equal(t, expectedDesc, actualWeatherDto.Description)
}

func assertCacheMetrics(
	t *testing.T,
	hitc prometheus.Counter,
	missc prometheus.Counter,
	errc prometheus.Counter,
	hit int,
	miss int,
	errCount int,
	delta float64,
) {
	t.Helper()

	require.InDelta(t, float64(hit), ptestutil.ToFloat64(hitc), delta)
	require.InDelta(t, float64(miss), ptestutil.ToFloat64(missc), delta)
	require.InDelta(t, float64(errCount), ptestutil.ToFloat64(errc), delta)
}
