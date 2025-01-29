package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/go-resty/resty/v2"
)

// Configuration holds the settings for the Resty client
type Configuration struct {
    BaseURL          string
    MaxRetries       int
    RetryWaitMin     time.Duration
    RetryWaitMax     time.Duration
    RetryMaxWaitTime time.Duration
    Timeout          time.Duration
}

// APIError represents a generic error response from the API
type APIError struct {
	StatusCode int
	Body       map[string]interface{}
}

func (e *APIError) Error() string {
	bodyBytes, _ := json.Marshal(e.Body)
	return fmt.Sprintf("API Error %d: %s", e.StatusCode, string(bodyBytes))
}

// InitializeClient sets up the Resty client with retry capabilities
func InitializeClient(config Configuration) *resty.Client {
	client := resty.New()

	client.
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout)

	localRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Configure Retry Mechanism
	client.
		SetRetryCount(config.MaxRetries).
		SetRetryWaitTime(time.Duration(config.RetryWaitMin) * time.Millisecond).
		SetRetryMaxWaitTime(time.Duration(config.RetryWaitMax) * time.Millisecond).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				// Retry on network errors
				if err != nil {
					log.Printf("Retry condition met due to error: %v", err)
					return true
				}
				// Retry on server errors (5xx) and too many requests (429)
				if r.StatusCode() >= 500 || r.StatusCode() == 429 {
					log.Printf("Retry condition met due to status code: %d", r.StatusCode())
					return true
				}
				return false
			},
		).
		SetRetryAfter(
			func(c *resty.Client, r *resty.Response) (time.Duration, error) {
				countIface, ok := r.Request.Context().Value("retry_count").(int)
				var count int
				if !ok {
					count = 1
				} else {
					count = countIface + 1
				}

				ctx := context.WithValue(r.Request.Context(), "retry_count", count)
				r.Request.SetContext(ctx)

				// Calculate exponential backoff with jitter
				min := float64(config.RetryWaitMin)
				max := float64(config.RetryWaitMax)
				backoff := min * math.Pow(2, float64(count))
				if backoff > max {
					backoff = max
				}
				// Add jitter: random value between 0 and 100ms
				jitter := localRand.Float64() * 100
				return time.Duration(backoff+jitter) * time.Millisecond, nil
			},
		)

	// Minimal Logging: Log requests and responses
	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		log.Printf("Sending %s request to %s", r.Method, r.URL)
		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		log.Printf("Received response with status %d from %s", r.StatusCode(), r.Request.URL)
		return nil
	})

	// Centralized Error Handling
	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		if r.IsError() {
			var apiErr map[string]interface{}
			if err := json.Unmarshal(r.Body(), &apiErr); err != nil {
				// If unmarshalling fails, return a generic error
				return fmt.Errorf("API returned status %d: %s", r.StatusCode(), r.Status())
			}
			return &APIError{
				StatusCode: r.StatusCode(),
				Body:       apiErr,
			}
		}
		return nil
	})

	return client
}
