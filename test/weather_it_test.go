//go:build integration

package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/emailclient"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config"
	"github.com/go-playground/validator/v10"
	"github.com/mailgun/mailgun-go/v4"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
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

func TestSubConfirmUnsubIT(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	form := url.Values{}
	form.Set("email", "test@example.com")
	form.Set("city", "Kyiv")
	form.Set("frequency", "daily")

	req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

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

	mailgunRespData, err := os.ReadFile("./test_data/mailgun_success_resp.json")
	require.NoError(t, err)

	var sentText string
	clientForMailgun := httputil.NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path == "/v3/domain/messages" && req.Method == http.MethodPost {
			err := req.ParseMultipartForm(10 << 20) // 10 MB max memory
			require.NoError(t, err)
			sentTextArr := req.MultipartForm.Value["text"]
			require.NotEmpty(t, sentTextArr)
			sentText = sentTextArr[0]
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(mailgunRespData)),
			Header:     make(http.Header),
		}, nil
	})
	mailgunClient := mailgun.NewMailgun("domain", "key")
	mailgunClient.SetClient(clientForMailgun)
	emailClient := emailclient.NewEmailClient(mailgunClient)

	subscriberRepository := postgresql.NewSubscriberRepository()
	subscriptionRepository := postgresql.NewSubscriptionRepository()
	tokenRepository := postgresql.NewTokenRepository()

	confirmEmailData := config.EmailData{
		Name:    "confirmation",
		Subject: "Confirm subscription",
		Text:    "To confirm your subscription use %s",
		From:    "sender@test.com",
	}
	confirmSuccessEmailData := config.EmailData{
		Name:    "confirmation-successful",
		Subject: "Confirmation successful",
		Text:    "You have successfully subscribed for weather update. To unsubscribe use %s",
		From:    "sender@test.com",
	}
	unsubEmailData := config.EmailData{
		Name:    "unsubscribe",
		Subject: "End of subscription",
		Text:    "You have successfully unsubscribed",
		From:    "sender@test.com",
	}

	subscriptionService := service.NewSubscriptionService(
		env.DB, weatherApiClient,
		subscriberRepository,
		subscriptionRepository,
		tokenRepository,
		emailClient,
		confirmEmailData,
		confirmSuccessEmailData,
		unsubEmailData,
		env.Log)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, validator.New(), env.Log)

	subscriptionHandler.Subscribe(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	require.NotEmpty(t, sentText)

	// Use regex to extract the token (UUID)
	re := regexp.MustCompile(`[0-9a-fA-F\-]{36}`) // UUID regex
	matches := re.FindStringSubmatch(sentText)
	require.NotEmpty(t, matches)
	token := matches[0]

	// Now use 'token' for your next test steps (e.g., confirm)
	fmt.Println("Extracted confirmation token:", token)

	//CONFIRM
	//req = httptest.NewRequest(http.MethodGet, "/confirm", nil)
	//req.SetPathValue("token", token)
	//subscriptionHandler.Confirm(rr, req)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /confirm/{token}", subscriptionHandler.Confirm)
	mux.HandleFunc("GET /unsubscribe/{token}", subscriptionHandler.Unsubscribe)

	req = httptest.NewRequest(http.MethodGet, "/confirm/"+token, nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	require.NotEmpty(t, sentText)

	matches = re.FindStringSubmatch(sentText)
	require.NotEmpty(t, matches)
	token = matches[0]

	fmt.Println("Extracted unsub token:", token)

	//unsub
	req = httptest.NewRequest(http.MethodGet, "/unsubscribe/"+token, nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.NotEmpty(t, sentText)
	require.Equal(t, unsubEmailData.Text, sentText)
}
