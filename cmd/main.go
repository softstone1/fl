package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/softstone1/fl/config"
	"github.com/softstone1/fl/internal/client"
	"github.com/softstone1/fl/internal/server"
	"github.com/softstone1/fl/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	 if err != nil {
        slog.Error("failed to load configuration", "err", err)
		os.Exit(1)
    }
	clientCfg := client.Configuration{
        MaxRetries: 	 cfg.ClientMaxRetries,
		RetryWaitMin: 	 cfg.ClientRetryWaitMin,
		RetryWaitMax: 	 cfg.ClientRetryWaitMax,
		RetryMaxWaitTime: cfg.ClientRetryMaxWaitTime,
		Timeout: 		 cfg.ClientTimeout,
	}
	
	clientCfg.BaseURL = cfg.LocationBaseURL
    locationClient := client.NewLocation(client.InitializeClient(clientCfg))
	
	clientCfg.BaseURL = cfg.ForecastBaseURL
    forecastClient := client.NewForecast(client.InitializeClient(clientCfg))

	
	forecastService, err := service.NewForecast(locationClient, forecastClient, 15 * time.Second, cfg.CacheSize)
	if err != nil {
		slog.Error("failed to create forecast service", "err", err)
		os.Exit(1)
	}

	router := server.NewRouter(forecastService)
	slog.Info("server starting", "port", cfg.ServerPort)
	if err := server.NewServer(cfg, router).Run(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed to start", "err", err)
	}
}