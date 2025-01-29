package handler

import (
	"testing"
	"time"

	"github.com/softstone1/fl/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestForcastPeriod(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(2 * time.Hour)

	forcast := model.ForcastPeriod{
		Number:           1,
		Name:             "Today",
		StartTime:        startTime,
		EndTime:          endTime,
		DetailedForecast: "Sunny with a chance of rain",
	}

	assert.Equal(t, 1, forcast.Number)
	assert.Equal(t, "Today", forcast.Name)
	assert.Equal(t, startTime, forcast.StartTime)
	assert.Equal(t, endTime, forcast.EndTime)
	assert.Equal(t, "Sunny with a chance of rain", forcast.DetailedForecast)
}
