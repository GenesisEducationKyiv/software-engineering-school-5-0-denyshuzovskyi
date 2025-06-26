package weatherstack

import (
	"errors"
	"fmt"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
)

var errTimestampMapping = errors.New("timestamp mapping error")

func CurrentWeatherToWeatherWithLocationDTO(currentWeather CurrentWeather) (dto.WeatherWithLocationDTO, error) {
	var description string
	if len(currentWeather.Current.WeatherDescriptions) > 0 {
		description = currentWeather.Current.WeatherDescriptions[0]
	}

	if currentWeather.Location.Localtime != "" &&
		currentWeather.Current.ObservationTime != "" &&
		currentWeather.Location.Timezone != "" {

		loc, locErr := time.LoadLocation(currentWeather.Location.Timezone)
		if locErr != nil {
			return dto.WeatherWithLocationDTO{}, fmt.Errorf("failed to load location: %w", locErr)
		}

		localTimeParsed, err := time.ParseInLocation("2006-01-02 15:04", currentWeather.Location.Localtime, loc)
		if err != nil {
			return dto.WeatherWithLocationDTO{}, fmt.Errorf("failed to parse localtime: %w", err)
		}

		utcDateStr := localTimeParsed.UTC().Format("2006-01-02")

		datetimeStr := utcDateStr + " " + currentWeather.Current.ObservationTime
		layout := "2006-01-02 03:04 PM"
		obsTimeUTC, err := time.ParseInLocation(layout, datetimeStr, time.UTC)
		if err != nil {
			return dto.WeatherWithLocationDTO{}, fmt.Errorf("failed to parse observation time: %w", err)
		}

		if obsTimeUTC.After(localTimeParsed.UTC()) {
			obsTimeUTC = obsTimeUTC.AddDate(0, 0, -1)
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
			LastUpdated: obsTimeUTC.Unix(),
		}, nil
	}

	return dto.WeatherWithLocationDTO{}, fmt.Errorf("insufficient data for timestamp mapping: %w", errTimestampMapping)
}
