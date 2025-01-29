package handler

import (
	"net/http"

	"github.com/softstone1/fl/internal/service"
)


type Forecast struct {
	ForcastService service.Forecast
}

func NewForecast(forcastService service.Forecast) *Forecast {
	return &Forecast{
		ForcastService: forcastService,
	}
}

// GetRandomForecast get current detailed forcast for a random location
func (f *Forecast) GetRandomForecast(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp, err := f.ForcastService.GetRandomForecast(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte(resp))
}
