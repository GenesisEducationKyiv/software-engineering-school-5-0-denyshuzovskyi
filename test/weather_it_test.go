//go:build integration

package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/emailclient"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config/email"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/testutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service/notification"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service/subscription"
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mailgun/mailgun-go/v4"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/weatherapi"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/database"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/httputil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/logger/noophandler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/repository/postgresql"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/server/handler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service/weather"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestEnv struct {
	DB      *sql.DB
	Log     *slog.Logger
	Cleanup func()
}

func SetupTestEnv(t *testing.T) *TestEnv {
	t.Helper()

	if testing.Short() {
		t.Skip()
	}

	if runtime.GOOS != "linux" {
		t.Skip("Works only on Linux (Testcontainers)")
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

func TestGetWeatherIT(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	currentWeatherData, err := os.ReadFile("./test_data/current_weather_success_resp.json")
	require.NoError(t, err)
	testClient := httputil.NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(currentWeatherData)),
			Header:     make(http.Header),
		}, nil
	})

	weatherApiClient := weatherapi.NewClient("https://api.weatherapi.com/v1", "key", testClient, env.Log)
	weatherRepository := postgresql.NewWeatherRepository()
	weatherService := weather.NewWeatherService(env.DB, weatherApiClient, weatherRepository, env.Log)
	weatherHandler := handler.NewWeatherHandler(weatherService, env.Log)

	city := "Kyiv"

	u := &url.URL{Path: "/weather"}
	q := u.Query()
	q.Set("city", city)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	weatherHandler.GetCurrentWeather(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var actualWeatherDto dto.WeatherDTO
	err = json.Unmarshal(rr.Body.Bytes(), &actualWeatherDto)
	require.NoError(t, err)

	expectedTemp := float32(6.6)
	expectedHum := float32(94)
	expectedDesc := "Light drizzle"
	delta := 0.01

	require.InDelta(t, expectedTemp, actualWeatherDto.Temperature, delta)
	require.InDelta(t, expectedHum, actualWeatherDto.Humidity, delta)
	require.Equal(t, expectedDesc, actualWeatherDto.Description)

	actualWeatherFromDB, err := weatherRepository.FindLastUpdatedByLocation(t.Context(), env.DB, city)
	require.NoError(t, err)
	require.NotNil(t, actualWeatherFromDB)
	require.InDelta(t, expectedTemp, actualWeatherFromDB.Temperature, delta)
	require.InDelta(t, expectedHum, actualWeatherFromDB.Humidity, delta)
	require.Equal(t, expectedDesc, actualWeatherFromDB.Description)
}

func TestWholeCycleIT(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// EmailData
	var cfg config.Config
	err := cleanenv.ReadConfig("./../config/config.yaml", &cfg)
	require.NoError(t, err)

	emailDataMap := email.PrepareEmailData(&cfg)
	confirmEmailData, confOk := emailDataMap["confirmation"]
	require.True(t, confOk)
	confirmSuccessEmailData, confSuccessOk := emailDataMap["confirmation-successful"]
	require.True(t, confSuccessOk)
	weatherEmailData, weatherOk := emailDataMap["weather"]
	require.True(t, weatherOk)
	unsubEmailData, unsubOk := emailDataMap["unsubscribe"]
	require.True(t, unsubOk)

	// Client-interceptor for weatherApiClient
	currentWeatherData, err := os.ReadFile("./test_data/current_weather_success_resp.json")
	require.NoError(t, err)
	waClient := httputil.NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(currentWeatherData)),
			Header:     make(http.Header),
		}, nil
	})
	weatherApiClient := weatherapi.NewClient("https://api.weatherapi.com/v1", "key", waClient, env.Log)

	// Client-interceptor for emailClient
	mailgunRespData, err := os.ReadFile("./test_data/mailgun_success_resp.json")
	require.NoError(t, err)
	var sentEmails []dto.SimpleEmail
	mgClient := httputil.NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path == "/v3/domain/messages" && req.Method == http.MethodPost {
			err := req.ParseMultipartForm(10 << 20) // 10 MB max memory
			require.NoError(t, err)

			f := req.MultipartForm.Value

			require.NotEmpty(t, f["from"])
			require.NotEmpty(t, f["to"])
			require.NotEmpty(t, f["subject"])
			require.NotEmpty(t, f["text"])

			e := dto.SimpleEmail{
				From:    f["from"][0],
				To:      f["to"][0],
				Subject: f["subject"][0],
				Text:    f["text"][0],
			}

			sentEmails = append(sentEmails, e)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(mailgunRespData)),
			Header:     make(http.Header),
		}, nil
	})
	mailgunClient := mailgun.NewMailgun("domain", "key")
	mailgunClient.SetClient(mgClient)
	emailClient := emailclient.NewEmailClient(mailgunClient)

	// Repositories, Services, Handlers
	weatherRepository := postgresql.NewWeatherRepository()
	subscriberRepository := postgresql.NewSubscriberRepository()
	subscriptionRepository := postgresql.NewSubscriptionRepository()
	tokenRepository := postgresql.NewTokenRepository()

	subscriptionService := subscription.NewSubscriptionService(
		env.DB, weatherApiClient,
		subscriberRepository,
		subscriptionRepository,
		tokenRepository,
		emailClient,
		confirmEmailData,
		confirmSuccessEmailData,
		unsubEmailData,
		env.Log)
	notificationService := notification.NewNotificationService(env.DB, weatherApiClient, weatherRepository, subscriberRepository, subscriptionRepository, tokenRepository, emailClient, env.Log)

	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, validator.New(), env.Log)

	// Multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("POST /subscribe", subscriptionHandler.Subscribe)
	mux.HandleFunc("GET /confirm/{token}", subscriptionHandler.Confirm)
	mux.HandleFunc("GET /unsubscribe/{token}", subscriptionHandler.Unsubscribe)

	// Subscribe
	subsForm := url.Values{}
	subsForm.Set("email", "test@example.com")
	subsForm.Set("city", "Kyiv")
	subsForm.Set("frequency", "daily")

	req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(subsForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	lastToken, err := testutil.ExtractFirstUUIDFromText(sentEmails[len(sentEmails)-1].Text)
	require.NoError(t, err)

	// Confirm
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/confirm/%s", lastToken), nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	lastToken, err = testutil.ExtractFirstUUIDFromText(sentEmails[len(sentEmails)-1].Text)
	require.NoError(t, err)

	// Imitate notification job trigger
	sentEmails = sentEmails[:0]
	require.Empty(t, sentEmails)
	notificationService.SendNotifications(t.Context(), model.Frequency_Hourly, weatherEmailData)
	require.Empty(t, sentEmails)
	notificationService.SendNotifications(t.Context(), model.Frequency_Daily, weatherEmailData)
	require.Len(t, sentEmails, 1)
	require.Equal(t, weatherEmailData.Subject, sentEmails[0].Subject)

	// Unsubscribe
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/unsubscribe/%s", lastToken), nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, unsubEmailData.Text, sentEmails[len(sentEmails)-1].Text)
}
