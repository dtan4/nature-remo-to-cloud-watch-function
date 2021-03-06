package natureremo

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	natureremoapi "github.com/tenntenn/natureremo"
)

// ClientInterface is an interface of a wrapper of Nature Remo API client
type ClientInterface interface {
	FetchTemperature(ctx context.Context, deviceID string) (float64, error)
}

// Client represents a wrapper of Nature Remo API client
type Client struct {
	api *natureremoapi.Client
}

// NewClient creates new Client object with a configured API client
func NewClient(accessToken string, httpClient *http.Client) *Client {
	api := natureremoapi.NewClient(accessToken)

	if httpClient != nil {
		api.HTTPClient = httpClient
	}

	return &Client{
		api: api,
	}
}

// FetchTemperature fetches the current room temperature from the specified device
func (c *Client) FetchTemperature(ctx context.Context, deviceID string) (float64, error) {
	devices, err := c.api.DeviceService.GetAll(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "cannot get all devices")
	}

	for _, d := range devices {
		if d.ID == deviceID {
			if t, ok := d.NewestEvents[natureremoapi.SensorTypeTemperature]; ok {
				return t.Value, nil
			}

			return 0, errors.Errorf("no temperature in device %q", deviceID)
		}
	}

	return 0, errors.Errorf("device %q not found", deviceID)
}
