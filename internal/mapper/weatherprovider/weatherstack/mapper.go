package weatherstack

import (
	"errors"
	"fmt"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/weatherstack"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
)

var errTimestampMapping = errors.New("timestamp mapping error")

func CurrentWeatherToWeatherWithLocationDTO(currentWeather weatherstack.CurrentWeather) (dto.WeatherWithLocationDTO, error) {
	description := extractDescription(currentWeather.Current.WeatherDescriptions)

	if currentWeather.Location.Localtime == "" ||
		currentWeather.Current.ObservationTime == "" ||
		currentWeather.Location.Timezone == "" {
		return dto.WeatherWithLocationDTO{}, fmt.Errorf("insufficient data for timestamp mapping: %w", errTimestampMapping)
	}

	lastUpdated, err := parseObservationUnix(
		currentWeather.Location.Localtime,
		currentWeather.Current.ObservationTime,
		currentWeather.Location.Timezone,
	)
	if err != nil {
		return dto.WeatherWithLocationDTO{}, err
	}

	return dto.WeatherWithLocationDTO{
		Weather: dto.WeatherDTO{
			Temperature: float32(currentWeather.Current.Temperature),
			Humidity:    float32(currentWeather.Current.Humidity),
			Description: description,
		},
		Location: dto.Location{
			Name: currentWeather.Location.Name,
		},
		LastUpdated: lastUpdated,
	}, nil
}

func extractDescription(descriptions []string) string {
	if len(descriptions) > 0 {
		return descriptions[0]
	}
	return ""
}

func parseObservationUnix(localtime, obsTime, timezone string) (int64, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, fmt.Errorf("failed to load location: %w", err)
	}
	localTime, err := time.ParseInLocation("2006-01-02 15:04", localtime, loc)
	if err != nil {
		return 0, fmt.Errorf("failed to parse localtime: %w", err)
	}

	utcDateStr := localTime.UTC().Format("2006-01-02")
	obsDateTimeStr := fmt.Sprintf("%s %s", utcDateStr, obsTime)
	obsTimeLayout := "2006-01-02 03:04 PM"

	obsTimeUTC, err := time.ParseInLocation(obsTimeLayout, obsDateTimeStr, time.UTC)
	if err != nil {
		return 0, fmt.Errorf("failed to parse observation time: %w", err)
	}

	if obsTimeUTC.After(localTime.UTC()) {
		obsTimeUTC = obsTimeUTC.AddDate(0, 0, -1)
	}
	return obsTimeUTC.Unix(), nil
}
