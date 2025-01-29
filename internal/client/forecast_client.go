package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/softstone1/fl/internal/model"
)

type Forecast interface {
	GetForecastURL(ctx context.Context, lat, lng float64) (string, error)
	GetForecastPeriods(ctx context.Context, forecastURL string) ([]model.ForcastPeriod, error)
}

type forecast struct {
	client *resty.Client
}

// make sure forcast implements the Forcast interface
var _ Forecast = (*forecast)(nil)

// ForecastURLResponse represents the top-level JSON structure
type ForecastURLResponse struct {
	Properties ForcastURLProperties `json:"properties"`
}

// ForcastURLProperties holds the forecast URL
type ForcastURLProperties struct {
	Forecast string `json:"forecast"`
}

// ForcastPeriodResponse represents the top-level JSON structure
type ForcastPeriodResponse struct {
	Properties ForcastPeriodsProperties `json:"properties"`
}

// ForcastPeriodsProperties holds the forecast periods
type ForcastPeriodsProperties struct {
	Periods []model.ForcastPeriod `json:"periods"`
}

// NewForecast initializes a new Forecast Client with a shared http.Client
func NewForecast(client *resty.Client) *forecast {
	return &forecast{
		client: client,
	}
}

// GetForecastURL returns the forcast URL for a given location
func (f *forecast) GetForecastURL(ctx context.Context, lat, lng float64) (string, error) {
	url := fmt.Sprintf("%s/points/%f,%f", f.client.BaseURL, lat, lng)
	forcastURL := &ForecastURLResponse{}
	resp, err := f.client.R().
		SetResult(forcastURL).
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Get(url)

	if err != nil {
		return "", fmt.Errorf("failed to fetch forecast: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to fetch forecast: %s", resp.Status())
	}
	return forcastURL.Properties.Forecast, nil
}

// GetForecastPeriods returns the forecast periods for a given forecast URL
func (f *forecast) GetForecastPeriods(ctx context.Context, forecastURL string) ([]model.ForcastPeriod, error) {
	forecastResponse := &ForcastPeriodResponse{}
	resp, err := f.client.R().
		SetResult(forecastResponse).
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Get(forecastURL)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast periods: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	return forecastResponse.Properties.Periods, nil
}
