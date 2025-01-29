
package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/softstone1/fl/internal/client"
	"github.com/softstone1/fl/internal/model"
)

func TestGetRandomForecast_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLocation := client.NewMockLocation(ctrl)
	mockForecast := client.NewMockForecast(ctrl)

	svc, _ := NewForecast(mockLocation, mockForecast, time.Second, 10)

	// Setup
	mockLocation.EXPECT().GetRandomLocation(gomock.Any()).Return(
		&model.Location{
			Name:      "Test Location",
			Latitude:  1.0,
			Longitude: 1.0,
		}, nil,
	)
	mockForecast.EXPECT().GetForecastURL(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		"http://test.url", nil,
	)
	mockForecast.EXPECT().GetForecastPeriods(gomock.Any(), gomock.Any()).Return(
		[]model.ForcastPeriod{
			{
				StartTime:       time.Now().Add(-time.Hour),
				EndTime:         time.Now().Add(time.Hour),
				DetailedForecast: "Sunny",
			},
		}, nil,
	)

	// Execute
	resp, err := svc.GetRandomForecast(context.Background())

	// Verify
	assert.NoError(t, err)
	assert.Contains(t, resp, "The weather in Test Location is: Sunny")
}

func TestGetRandomForecast_LocationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLocation := client.NewMockLocation(ctrl)
	mockForecast := client.NewMockForecast(ctrl)

	svc, _ := NewForecast(mockLocation, mockForecast, time.Second, 10)

	// Setup
	mockLocation.EXPECT().GetRandomLocation(gomock.Any()).Return(
		nil,
		errors.New("location error"),
	)

	// Execute
	resp, err := svc.GetRandomForecast(context.Background())

	// Verify
	assert.Error(t, err)
	assert.Empty(t, resp)
}

func TestGetRandomForecast_NoCurrentDetailedForecast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLocation := client.NewMockLocation(ctrl)
	mockForecast := client.NewMockForecast(ctrl)

	svc, _ := NewForecast(mockLocation, mockForecast, time.Second, 10)

	mockLocation.EXPECT().GetRandomLocation(gomock.Any()).Return(
		&model.Location{
			Name:      "Test Location",
			Latitude:  1.0,
			Longitude: 1.0,
		}, nil,
	)
	mockForecast.EXPECT().GetForecastURL(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		"http://test.url", nil,
	)
	mockForecast.EXPECT().GetForecastPeriods(gomock.Any(), gomock.Any()).Return(
		[]model.ForcastPeriod{}, nil,
	)

	resp, err := svc.GetRandomForecast(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no current detailed forecast found")
	assert.Empty(t, resp)
}