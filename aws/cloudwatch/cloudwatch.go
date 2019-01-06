package cloudwatch

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

// ClientInterface is an interface of a wrapper of CloudWatch API client
type ClientInterface interface {
	PutTemperature(deviceID string, temperature float64) error
}

// Client is a wrapper of CloudWatch API client
type Client struct {
	api cloudwatchiface.CloudWatchAPI
}

// NewClient creates new Client object with the given API client
func NewClient(api cloudwatchiface.CloudWatchAPI) *Client {
	return &Client{
		api: api,
	}
}
