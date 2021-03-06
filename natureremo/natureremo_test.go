package natureremo

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	natureremoapi "github.com/tenntenn/natureremo"
)

type fakeDeviceService struct {
	Devices []*natureremoapi.Device
	Err     error
}

func (f fakeDeviceService) GetAll(ctx context.Context) ([]*natureremoapi.Device, error) {
	if f.Err != nil {
		return []*natureremoapi.Device{}, f.Err
	}

	return f.Devices, nil
}

func (f fakeDeviceService) Update(ctx context.Context, device *natureremoapi.Device) (*natureremoapi.Device, error) {
	return nil, nil
}

func (f fakeDeviceService) Delete(ctx context.Context, device *natureremoapi.Device) error {
	return nil
}

func (f fakeDeviceService) UpdateTemperatureOffset(ctx context.Context, device *natureremoapi.Device) (*natureremoapi.Device, error) {
	return nil, nil
}

func (f fakeDeviceService) UpdateHumidityOffset(ctx context.Context, device *natureremoapi.Device) (*natureremoapi.Device, error) {
	return nil, nil
}

func TestNewClient(t *testing.T) {
	testcases := []struct {
		accessToken string
		httpClient  *http.Client
	}{
		{
			accessToken: "",
			httpClient:  nil,
		},
		{
			accessToken: "",
			httpClient:  http.DefaultClient,
		},
		{
			accessToken: "dummyaccesstoken",
			httpClient:  nil,
		},
	}

	for _, tc := range testcases {
		got := NewClient(tc.accessToken, tc.httpClient)
		if got == nil {
			t.Error("want client, got nil")
		}
	}
}

func TestFetchTemperature(t *testing.T) {
	testcases := []struct {
		deviceID string
		devices  []*natureremoapi.Device
		want     float64
	}{
		{
			deviceID: "91246eb0-4e06-4f1a-a400-42874839aee1",
			devices: []*natureremoapi.Device{
				{
					DeviceCore: natureremoapi.DeviceCore{
						ID:   "1ffdbcec-12ed-4694-aadc-3f773d8363d4",
						Name: "Main Room",
					},
					NewestEvents: map[natureremoapi.SensorType]natureremoapi.SensorValue{
						natureremoapi.SensorTypeTemperature: {
							Value: 18.17,
						},
					},
				},
				{
					DeviceCore: natureremoapi.DeviceCore{
						ID:   "91246eb0-4e06-4f1a-a400-42874839aee1",
						Name: "Bed Room",
					},
					NewestEvents: map[natureremoapi.SensorType]natureremoapi.SensorValue{
						natureremoapi.SensorTypeTemperature: {
							Value: 21.39,
						},
					},
				},
			},
			want: 21.39,
		},
	}

	for _, tc := range testcases {
		ctx := context.Background()
		client := &Client{
			api: &natureremoapi.Client{
				DeviceService: &fakeDeviceService{
					Devices: tc.devices,
				},
			},
		}

		got, err := client.FetchTemperature(ctx, tc.deviceID)
		if err != nil {
			t.Fatalf("want no error, got: %s", err)
		}

		if got != tc.want {
			t.Errorf("want: %f, got: %f", tc.want, got)
		}
	}
}

func TestFetchTemperature_error(t *testing.T) {
	testcases := []struct {
		subtitle string
		deviceID string
		devices  []*natureremoapi.Device
		err      error
		want     string
	}{
		{
			subtitle: "api error",
			deviceID: "91246eb0-4e06-4f1a-a400-42874839aee1",
			err:      fmt.Errorf("unexpected error"),
			want:     `cannot get all devices: unexpected error`,
		},
		{
			subtitle: "device not found",
			deviceID: "91246eb0-4e06-4f1a-a400-42874839aee1",
			devices: []*natureremoapi.Device{
				{
					DeviceCore: natureremoapi.DeviceCore{
						ID:   "1ffdbcec-12ed-4694-aadc-3f773d8363d4",
						Name: "Main Room",
					},
					NewestEvents: map[natureremoapi.SensorType]natureremoapi.SensorValue{
						natureremoapi.SensorTypeTemperature: {
							Value: 18.17,
						},
					},
				},
			},
			want: `device "91246eb0-4e06-4f1a-a400-42874839aee1" not found`,
		},
		{
			subtitle: "temperature not found",
			deviceID: "91246eb0-4e06-4f1a-a400-42874839aee1",
			devices: []*natureremoapi.Device{
				{
					DeviceCore: natureremoapi.DeviceCore{
						ID:   "91246eb0-4e06-4f1a-a400-42874839aee1",
						Name: "Main Room",
					},
					NewestEvents: map[natureremoapi.SensorType]natureremoapi.SensorValue{},
				},
			},
			want: `no temperature in device "91246eb0-4e06-4f1a-a400-42874839aee1"`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.subtitle, func(t *testing.T) {
			ctx := context.Background()
			client := &Client{
				api: &natureremoapi.Client{
					DeviceService: &fakeDeviceService{
						Devices: tc.devices,
						Err:     tc.err,
					},
				},
			}

			_, err := client.FetchTemperature(ctx, tc.deviceID)
			if err == nil {
				t.Fatal("want error, got nil")
			}

			if err.Error() != tc.want {
				t.Errorf("want: %q, got: %q", tc.want, err.Error())
			}
		})
	}
}
