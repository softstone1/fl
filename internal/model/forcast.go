package model

import "time"

type ForcastPeriod struct {
	Number           int       `json:"number"`
	Name             string    `json:"name"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	DetailedForecast string    `json:"detailedForecast"`
}
