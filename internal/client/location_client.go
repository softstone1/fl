package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"github.com/go-resty/resty/v2"
	"github.com/softstone1/fl/internal/model"
)

type Location interface {
	GetRandomLocation(ctx context.Context) (*model.Location, error)
}

type location struct {
	client *resty.Client
}
// make sure location implements the Location interface
var _ Location = (*location)(nil)

// LocationResponse represents the entire JSON response
type LocationResponse struct {
    Locations []model.Location `json:"locations"`
}

// NewLocation initializes a new Location Client with a shared http.Client
func NewLocation(client *resty.Client) *location {
	return &location{
		client: client,
	}
}

// GetRandomLocation fetches a random location from the API
func (l *location) GetRandomLocation(ctx context.Context) (*model.Location, error) {
	url := fmt.Sprintf("%s/api/random", l.client.BaseURL)
	locationResponse := &LocationResponse{}
	resp, err := l.client.R().
	SetResult(locationResponse).
	SetContext(ctx).
	SetHeader("Accept", "application/json").
	Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get random location: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	
	if len(locationResponse.Locations) == 0 {
		return nil, errors.New("no locations found")
	}
	return &locationResponse.Locations[0], nil
}
