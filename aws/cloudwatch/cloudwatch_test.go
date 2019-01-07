package cloudwatch

import (
	"fmt"
	"testing"
	"time"

	"github.com/dtan4/nature-remo-to-cloud-watch-function/aws/mock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/golang/mock/gomock"
)

//go:generate mockgen -source $GOPATH/pkg/mod/github.com/aws/aws-sdk-go@v1.16.14/service/cloudwatch/cloudwatchiface/interface.go -destination ../mock/cloudwatch.go -package mock

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cloudwatchMock := mock.NewMockCloudWatchAPI(ctrl)

	got := NewClient(cloudwatchMock)
	if got == nil {
		t.Error("want objcet, got nil")
	}
}

func TestPutTemperature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
		cloudwatchMock := mock.NewMockCloudWatchAPI(ctrl)
		cloudwatchMock.EXPECT().PutMetricData(&cloudwatch.PutMetricDataInput{
			Namespace: aws.String("NatureRemo/RoomMetrics"),
			MetricData: []*cloudwatch.MetricDatum{
				{
					MetricName: aws.String("Temperature"),
					Timestamp:  aws.Time(tc.timestamp),
					Value:      aws.Float64(tc.temperature),
					Dimensions: []*cloudwatch.Dimension{
						{
							Name:  aws.String("DeviceID"),
							Value: aws.String(tc.deviceID),
						},
					},
				},
			},
		}).Return(&cloudwatch.PutMetricDataOutput{}, nil)

		client := &Client{
			api: cloudwatchMock,
		}

		err := client.PutTemperature(tc.timestamp, tc.deviceID, tc.temperature)
		if err != nil {
			t.Errorf("want no error, got: %s", err)
		}
	}
}

func TestPutTemperature_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		subtitle    string
		timestamp   time.Time
		deviceID    string
		temperature float64
		want        string
	}{
		{
			subtitle:    "api error",
			timestamp:   time.Date(2019, 1, 7, 2, 39, 24, 0, time.UTC),
			deviceID:    "91246eb0-4e06-4f1a-a400-42874839aee1",
			temperature: 18.17,
			want:        "cannot put metric: unexpected error",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.subtitle, func(t *testing.T) {
			cloudwatchMock := mock.NewMockCloudWatchAPI(ctrl)
			cloudwatchMock.EXPECT().PutMetricData(&cloudwatch.PutMetricDataInput{
				Namespace: aws.String("NatureRemo/RoomMetrics"),
				MetricData: []*cloudwatch.MetricDatum{
					{
						MetricName: aws.String("Temperature"),
						Timestamp:  aws.Time(tc.timestamp),
						Value:      aws.Float64(tc.temperature),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("DeviceID"),
								Value: aws.String(tc.deviceID),
							},
						},
					},
				},
			}).Return(nil, fmt.Errorf("unexpected error"))

			client := &Client{
				api: cloudwatchMock,
			}

			err := client.PutTemperature(tc.timestamp, tc.deviceID, tc.temperature)
			if err == nil {
				t.Fatalf("want error, got nil")
			}

			if err.Error() != tc.want {
				t.Errorf("want: %q, got: %q", tc.want, err.Error())
			}
		})
	}
}
