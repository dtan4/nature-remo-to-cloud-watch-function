package cloudwatch

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

type mockCloudWatchAPI struct {
	cloudwatchiface.CloudWatchAPI
	putMetricDataWithContextFunc func(ctx aws.Context, input *cloudwatch.PutMetricDataInput, opts ...request.Option) (*cloudwatch.PutMetricDataOutput, error)
}

func (m *mockCloudWatchAPI) PutMetricDataWithContext(ctx aws.Context, input *cloudwatch.PutMetricDataInput, opts ...request.Option) (*cloudwatch.PutMetricDataOutput, error) {
	return m.putMetricDataWithContextFunc(ctx, input, opts...)
}

func TestNewClient(t *testing.T) {
	api := &mockCloudWatchAPI{}
	client := NewClient(api)

	if client.api != api {
		t.Error("api does not match")
	}
}

func TestPutTemperature(t *testing.T) {
	testcases := []struct {
		timestamp   time.Time
		deviceID    string
		temperature float64
	}{
		{
			timestamp:   time.Date(2019, 1, 7, 2, 39, 24, 0, time.UTC),
			deviceID:    "91246eb0-4e06-4f1a-a400-42874839aee1",
			temperature: 18.17,
		},
	}

	for _, tc := range testcases {
		ctx := context.Background()

		client := &Client{
			api: &mockCloudWatchAPI{
				putMetricDataWithContextFunc: func(ctx aws.Context, input *cloudwatch.PutMetricDataInput, opts ...request.Option) (*cloudwatch.PutMetricDataOutput, error) {
					return &cloudwatch.PutMetricDataOutput{}, nil
				},
			},
		}

		err := client.PutTemperature(ctx, tc.timestamp, tc.deviceID, tc.temperature)
		if err != nil {
			t.Errorf("want no error, got: %s", err)
		}
	}
}

func TestPutTemperature_error(t *testing.T) {
	testcases := []struct {
		subtitle    string
		timestamp   time.Time
		deviceID    string
		temperature float64
		err         string
		want        string
	}{
		{
			subtitle:    "api error",
			timestamp:   time.Date(2019, 1, 7, 2, 39, 24, 0, time.UTC),
			deviceID:    "91246eb0-4e06-4f1a-a400-42874839aee1",
			temperature: 18.17,
			err:         "unexpected error",
			want:        "cannot put metric: unexpected error",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.subtitle, func(t *testing.T) {
			ctx := context.Background()

			client := &Client{
				api: &mockCloudWatchAPI{
					putMetricDataWithContextFunc: func(ctx aws.Context, input *cloudwatch.PutMetricDataInput, opts ...request.Option) (*cloudwatch.PutMetricDataOutput, error) {
						return nil, fmt.Errorf(tc.err)
					},
				},
			}

			err := client.PutTemperature(ctx, tc.timestamp, tc.deviceID, tc.temperature)
			if err == nil {
				t.Fatalf("want error, got nil")
			}

			if err.Error() != tc.want {
				t.Errorf("want: %q, got: %q", tc.want, err.Error())
			}
		})
	}
}
