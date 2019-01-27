package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type fakeCloudWatchClient struct {
	Err error
}

func (c fakeCloudWatchClient) PutTemperature(ctx context.Context, timestamp time.Time, deviceID string, temperature float64) error {
	if c.Err != nil {
		return c.Err
	}

	return nil
}

type fakeSSMClient struct {
	Err   error
	Value string
}

func (c fakeSSMClient) LoadSecret(ctx context.Context, name string) (string, error) {
	if c.Err != nil {
		return "", c.Err
	}

	return c.Value, nil
}

type fakeNatureRemoClient struct {
	Err         error
	Temperature float64
}

func (c fakeNatureRemoClient) FetchTemperature(ctx context.Context, deviceID string) (float64, error) {
	if c.Err != nil {
		return 0, c.Err
	}

	return c.Temperature, nil
}

func TestRealHandler(t *testing.T) {
	testcases := []struct {
		temperature float64
		value       string
	}{
		{
			temperature: 18.17,
			value:       "foobarbaz",
		},
	}

	for _, tc := range testcases {
		ctx := context.Background()
		CloudWatchClient = fakeCloudWatchClient{}
		SSMClient = fakeSSMClient{Value: tc.value}
		NatureRemoClient = fakeNatureRemoClient{Temperature: tc.temperature}

		err := RealHandler(ctx)
		if err != nil {
			t.Errorf("want no error, got: %s", err)
		}
	}
}

func TestRealHandler_error(t *testing.T) {
	testcases := []struct {
		subtitle      string
		cloudWatchErr error
		ssmErr        error
		natureRemoErr error
		want          string
	}{
		{
			subtitle:      "cloudwatch",
			cloudWatchErr: fmt.Errorf("unexpected error"),
			want:          "cannot put room temperature: unexpected error",
		},
		{
			subtitle: "ssm",
			ssmErr:   fmt.Errorf("unexpected error"),
			want:     "cannot load device ID: unexpected error",
		},
		{
			subtitle:      "natureremo",
			natureRemoErr: fmt.Errorf("unexpected error"),
			want:          "cannot fetch room temperature: unexpected error",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.subtitle, func(t *testing.T) {
			ctx := context.Background()
			CloudWatchClient = fakeCloudWatchClient{Err: tc.cloudWatchErr}
			SSMClient = fakeSSMClient{Err: tc.ssmErr}
			NatureRemoClient = fakeNatureRemoClient{Err: tc.natureRemoErr}

			err := RealHandler(ctx)
			if err == nil {
				t.Fatalf("want error, got nil")
			}

			if err.Error() != tc.want {
				t.Errorf("want: %q, got: %q", tc.want, err.Error())
			}
		})
	}
}
