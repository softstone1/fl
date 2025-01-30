package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/softstone1/fl/internal/service"
	"go.uber.org/mock/gomock"
)

func TestGetRandomForecast(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockForecastSvc := service.NewMockForecast(ctrl)
    h := NewForecast(mockForecastSvc)

    tests := []struct {
        name           string
        mockSetup      func()
        expectedStatus int
        expectedBody   string
    }{
        {
            name: "Success returns 200 with 'sunny'",
            mockSetup: func() {
                mockForecastSvc.
                    EXPECT().
                    GetRandomForecast(gomock.Any()).
                    Return("sunny", nil)
            },
            expectedStatus: http.StatusOK,
            expectedBody:   "sunny",
        },
        {
            name: "Service error returns 500",
            mockSetup: func() {
                mockForecastSvc.
                    EXPECT().
                    GetRandomForecast(gomock.Any()).
                    Return("", errors.New("some error"))
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody: "some error\n",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.mockSetup()

            req := httptest.NewRequest(http.MethodGet, "/forecast", nil)
            rr := httptest.NewRecorder()

            h.GetRandomForecast(rr, req)

            if rr.Code != tc.expectedStatus {
                t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
            }

            if body := rr.Body.String(); body != tc.expectedBody {
                t.Errorf("expected body %q, got %q", tc.expectedBody, body)
            }
        })
    }
}