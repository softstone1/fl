package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort             int
	ServerReadTimeout      time.Duration
	ServerWriteTimeout     time.Duration
	ServerIdleTimeout      time.Duration
	LocationBaseURL        string
	ForecastBaseURL        string
	ClientMaxRetries       int
	ClientRetryWaitMin     time.Duration
	ClientRetryWaitMax     time.Duration
	ClientRetryMaxWaitTime time.Duration
	ClientTimeout          time.Duration
	CacheSize              int
}

// LoadConfig parses configuration from environment variables and command-line flags
func LoadConfig() (*Config, error) {
	config := &Config{}

	flag.IntVar(&config.ServerPort, "port", 5000, "server port")
	flag.DurationVar(&config.ServerReadTimeout, "read-timeout", 15*time.Second, "server read timeout")
	flag.DurationVar(&config.ServerWriteTimeout, "write-timeout", 15*time.Second, "server write timeout")
	flag.DurationVar(&config.ServerIdleTimeout, "idle-timeout", 60*time.Second, "server idle timeout")
	flag.StringVar(&config.LocationBaseURL, "location-base-url", "https://locations.patch3s.dev", "location base URL")
	flag.StringVar(&config.ForecastBaseURL, "forecast-base-url", "https://api.weather.gov", "forecast base URL")
	flag.IntVar(&config.ClientMaxRetries, "client-max-retries", 5, "client max retries")
	flag.DurationVar(&config.ClientRetryWaitMin, "client-retry-wait-min", 500*time.Millisecond, "client retry wait min")
	flag.DurationVar(&config.ClientRetryWaitMax, "client-retry-wait-max", 5000*time.Millisecond, "client retry wait max")
	flag.DurationVar(&config.ClientRetryMaxWaitTime, "client-retry-max-wait-time", 30000*time.Millisecond, "client retry max wait time")
	flag.DurationVar(&config.ClientTimeout, "client-timeout", 15*time.Second, "client timeout")
	flag.IntVar(&config.CacheSize, "cache-size", 1000, "cache size")

	flag.Parse()

	// override with environment variables

	if portEnv, ok := os.LookupEnv("PORT"); ok {
		port, err := strconv.Atoi(portEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PORT: %w", err)
		}
		config.ServerPort = port
	}

	if readTimeoutEnv, ok := os.LookupEnv("SERVER_READ_TIMEOUT"); ok {
		readTimeout, err := time.ParseDuration(readTimeoutEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse SERVER_READ_TIMEOUT: %w", err)
		}
		config.ServerReadTimeout = readTimeout
	}

	if writeTimeoutEnv, ok := os.LookupEnv("SERVER_WRITE_TIMEOUT"); ok {
		writeTimeout, err := time.ParseDuration(writeTimeoutEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse SERVER_WRITE_TIMEOUT: %w", err)
		}
		config.ServerWriteTimeout = writeTimeout
	}

	if idleTimeoutEnv, ok := os.LookupEnv("SERVER_IDLE_TIMEOUT"); ok {
		idleTimeout, err := time.ParseDuration(idleTimeoutEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse SERVER_IDLE_TIMEOUT: %w", err)
		}
		config.ServerIdleTimeout = idleTimeout
	}

	if locationBaseURLEnv, ok := os.LookupEnv("LOCATION_BASE_URL"); ok {
		config.LocationBaseURL = locationBaseURLEnv
	}

	if forecastBaseURLEnv, ok := os.LookupEnv("FORECAST_BASE_URL"); ok {
		config.ForecastBaseURL = forecastBaseURLEnv
	}

	if clientMaxRetriesEnv, ok := os.LookupEnv("CLIENT_MAX_RETRIES"); ok {
		clientMaxRetries, err := strconv.Atoi(clientMaxRetriesEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CLIENT_MAX_RETRIES: %w", err)
		}
		config.ClientMaxRetries = clientMaxRetries
	}

	if clientRetryWaitMinEnv, ok := os.LookupEnv("CLIENT_RETRY_WAIT_MIN"); ok {
		clientRetryWaitMin, err := time.ParseDuration(clientRetryWaitMinEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CLIENT_RETRY_WAIT_MIN: %w", err)
		}
		config.ClientRetryWaitMin = clientRetryWaitMin
	}

	if clientRetryWaitMaxEnv, ok := os.LookupEnv("CLIENT_RETRY_WAIT_MAX"); ok {
		clientRetryWaitMax, err := time.ParseDuration(clientRetryWaitMaxEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CLIENT_RETRY_WAIT_MAX: %w", err)
		}
		config.ClientRetryWaitMax = clientRetryWaitMax
	}

	if clientRetryMaxWaitTimeEnv, ok := os.LookupEnv("CLIENT_RETRY_MAX_WAIT_TIME"); ok {
		clientRetryMaxWaitTime, err := time.ParseDuration(clientRetryMaxWaitTimeEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CLIENT_RETRY_MAX_WAIT_TIME: %w", err)
		}
		config.ClientRetryMaxWaitTime = clientRetryMaxWaitTime
	}

	if clientTimeoutEnv, ok := os.LookupEnv("CLIENT_TIMEOUT"); ok {
		clientTimeout, err := time.ParseDuration(clientTimeoutEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CLIENT_TIMEOUT: %w", err)
		}
		config.ClientTimeout = clientTimeout
	}

	if cacheSizeEnv, ok := os.LookupEnv("CACHE_SIZE"); ok {
		cacheSize, err := strconv.Atoi(cacheSizeEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CACHE_SIZE: %w", err)
		}
		config.CacheSize = cacheSize
	}

	// validate configuration
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return fmt.Errorf("invalid port number: %d", c.ServerPort)
	}
	if c.ServerReadTimeout <= 0 {
		return fmt.Errorf("read_timeout must be positive")
	}
	if c.ServerWriteTimeout <= 0 {
		return fmt.Errorf("write_timeout must be positive")
	}
	if c.ServerIdleTimeout <= 0 {
		return fmt.Errorf("idle_timeout must be positive")
	}
	if c.LocationBaseURL == "" {
		return fmt.Errorf("location_base_url cannot be empty")
	}
	if c.ForecastBaseURL == "" {
		return fmt.Errorf("forecast_base_url cannot be empty")
	}
	if c.ClientMaxRetries < 0 {
		return fmt.Errorf("client_max_retries must be non-negative")
	}
	if c.ClientRetryWaitMin <= 0 {
		return fmt.Errorf("client_retry_wait_min must be positive")
	}
	if c.ClientRetryWaitMax <= 0 {
		return fmt.Errorf("client_retry_wait_max must be positive")
	}
	if c.ClientRetryMaxWaitTime <= 0 {
		return fmt.Errorf("client_retry_max_wait_time must be positive")
	}
	if c.ClientTimeout <= 0 {
		return fmt.Errorf("client_timeout must be positive")
	}

	if c.CacheSize <= 0 {
		return fmt.Errorf("cache_size must be positive")
	}

	return nil
}
