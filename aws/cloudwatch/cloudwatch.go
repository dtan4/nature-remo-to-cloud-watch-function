package cloudwatch

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/pkg/errors"
)

const (
	metricNamespace   = "NatureRemo/RoomMetrics"
	metricName        = "Temperature"
	dimensionDeviceID = "DeviceID"
)

// ClientInterface is an interface of a wrapper of CloudWatch API client
type ClientInterface interface {
	PutTemperature(ctx context.Context, timestamp time.Time, deviceID string, temperature float64) error
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

// PutTemperature puts the given room temperature as CloudWatch metric
func (c *Client) PutTemperature(ctx context.Context, timestamp time.Time, deviceID string, temperature float64) error {
	_, err := c.api.PutMetricDataWithContext(ctx, &cloudwatch.PutMetricDataInput{
		Namespace: aws.String(metricNamespace),
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String(metricName),
				Timestamp:  aws.Time(timestamp),
				Value:      aws.Float64(temperature),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String(dimensionDeviceID),
						Value: aws.String(deviceID),
					},
				},
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "cannot put metric")
	}

	return nil
}
