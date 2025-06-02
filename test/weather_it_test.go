package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/weatherapi"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/handler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/httputil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/logger/noophandler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/repository/posgresql"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/migrations"
	"github.com/golang-migrate/migrate/v4"
	mpostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime"
	"testing"
)

type TestEnv struct {
	Log     *slog.Logger
	DB      *sql.DB
	Cleanup func()
}

func SetupTestEnv(t *testing.T) *TestEnv {
	if testing.Short() {
		t.Skip()
	}

	if runtime.GOOS != "linux" {
		t.Skip("Works only on Linux (Testcontainers)")
	}

	ctx := context.Background()

	ctr, err := postgres.Run(ctx, "postgres:17-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	require.NoError(t, err)

	conn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	t.Log(conn)

	db, err := sql.Open("pgx", conn)
	require.NoError(t, err)

	driver, err := mpostgres.WithInstance(db, &mpostgres.Config{})
	require.NoError(t, err)

	d, err := iofs.New(migrations.Files, ".")
	require.NoError(t, err)

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	require.NoError(t, err)

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("unable to apply migrations: %v", err)
	} else {
		t.Log("migration completed successfully")
	}

	log := slog.New(noophandler.NewNoOpHandler())

	return &TestEnv{
		Log: log,
		DB:  db,
		Cleanup: func() {
			require.NoError(t, db.Close())
			require.NoError(t, ctr.Terminate(ctx))
		},
	}
}

func TestWeatherHandlerIT(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	currentWeatherData, err := os.ReadFile("./test_data/current_weather_resp.json")
	require.NoError(t, err)

	testClient := httputil.NewTestHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(currentWeatherData)),
			Header:     make(http.Header),
		}, nil
	})

	weatherApiClient := weatherapi.NewClient("https://api.weatherapi.com/v1", "key", testClient, env.Log)
	locationRepository := posgresql.NewLocationRepository()
	weatherRepository := posgresql.NewWeatherRepository()
	weatherService := service.NewWeatherService(env.DB, weatherApiClient, locationRepository, weatherRepository, env.Log)
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

	require.Equal(t, expectedTemp, actualWeatherDto.Temperature)
	require.Equal(t, expectedHum, actualWeatherDto.Humidity)
	require.Equal(t, expectedDesc, actualWeatherDto.Description)

	actualWeatherFroDB, err := weatherRepository.FindLastUpdatedByLocation(context.Background(), env.DB, city)
	require.NoError(t, err)
	require.NotNil(t, actualWeatherFroDB)
	require.Equal(t, expectedTemp, actualWeatherDto.Temperature)
	require.Equal(t, expectedHum, actualWeatherDto.Humidity)
	require.Equal(t, expectedDesc, actualWeatherDto.Description)
}
