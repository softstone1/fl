package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/softstone1/fl/internal/client"
	"github.com/softstone1/fl/internal/model"
)

type Forecast interface {
	GetRandomForecast(ctx context.Context) (string, error)
}

type forecast struct {
	LocationClient       client.Location
	ForcastClient        client.Forecast
	Timeout              time.Duration
	forecastURLCache     *lru.Cache[string, string]
	forecastPeriodsCache *lru.Cache[string, []model.ForcastPeriod]
}

// make sure forcast implements the Forcast interface
var _ Forecast = (*forecast)(nil)

func NewForecast(locationClient client.Location, forcastClient client.Forecast, timeout time.Duration, cacheSize int) (*forecast, error) {

	forecastURLCache, err := lru.New[string, string](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create forecastURLCache: %w", err)
	}
	forecastPeriodsCache, err := lru.New[string, []model.ForcastPeriod](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create forecastPeriodsCache: %w", err)
	}

	return &forecast{
		LocationClient:       locationClient,
		ForcastClient:        forcastClient,
		Timeout:              timeout,
		forecastURLCache:     forecastURLCache,
		forecastPeriodsCache: forecastPeriodsCache,
	}, nil
}

// GetRandomForecast orchestrates fetching random location, forecast URL, and current detailed forcast with timeout and caching
func (s *forecast) GetRandomForecast(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()

	location, err := s.LocationClient.GetRandomLocation(ctx)
	if err != nil {
		return "", fmt.Errorf("Stage 1 - FetchLocation error: %w", err)
	}

	forecastURL, err := s.getForecastURL(ctx, location.Latitude, location.Longitude)
	if err != nil {
		return "", fmt.Errorf("Stage 2 - GetForecastURL error: %w", err)
	}

	detailedForcast, err := s.getCurrentDetailedForcast(ctx, forecastURL)
	if err != nil {
		return "", fmt.Errorf("Stage 3 - GetForecastResponse error: %w", err)
	}

	return fmt.Sprintf("The weather in %s is: %s", location.Name, detailedForcast), nil
}

// getForecastURL retrieves the forecast URL, utilizing the cache if available
func (s *forecast) getForecastURL(ctx context.Context, lat, lng float64) (string, error) {
	// Define cache key
	cacheKey := fmt.Sprintf("%f,%f", lat, lng)

	// Check if forecast URL is cached
	if cachedURL, found := s.forecastURLCache.Get(cacheKey); found {
		return cachedURL, nil
	}

	// Fetch forecast URL from external API
	forecastURL, err := s.ForcastClient.GetForecastURL(ctx, lat, lng)
	if err != nil {
		return "", err
	}

	// Store the fetched forecast URL in cache
	s.forecastURLCache.Add(cacheKey, forecastURL)

	return forecastURL, nil
}

// getCurrentDetailedForcast retrieves the current detailed forecast, utilizing the cache if available
func (s *forecast) getCurrentDetailedForcast(ctx context.Context, forecastURL string) (string, error) {

	// Define cache key
	cacheKey := forecastURL

	// Check if forecast response is cached
	if cachedResp, found := s.forecastPeriodsCache.Get(cacheKey); found {
		current := time.Now()
		for _, period := range cachedResp {
			if period.StartTime.Before(current) && period.EndTime.After(current) {
				return period.DetailedForecast, nil
			}
		}
	}

	// Fetch forecast response from external API
	forecastResponse, err := s.ForcastClient.GetForecastPeriods(ctx, forecastURL)
	if err != nil {
		return "", err
	}

	// Find the current detailed forecast
	current := time.Now()
	for _, period := range forecastResponse {
		if period.StartTime.Before(current) && period.EndTime.After(current) {
			// Store the fetched forecast response in cache
			s.forecastPeriodsCache.Add(cacheKey, forecastResponse)
			return period.DetailedForecast, nil
		}
	}

	return "", errors.New("no current detailed forecast found")

}
