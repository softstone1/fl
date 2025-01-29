package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/softstone1/fl/internal/handler"
	"github.com/softstone1/fl/internal/service"
)

// NewRouter constructs the main router for your app.
func NewRouter(forcastService service.Forecast) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		// forcast handler
		forcastHandler := handler.NewForecast(forcastService)
		r.Get("/", forcastHandler.GetRandomForecast)
	})

	return r
}

