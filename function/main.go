package function

import (
	"context"
	"time"

	"github.com/dtan4/nature-remo-to-cloud-watch-function/aws/cloudwatch"
	"github.com/dtan4/nature-remo-to-cloud-watch-function/aws/ssm"
	"github.com/dtan4/nature-remo-to-cloud-watch-function/natureremo"

	"github.com/pkg/errors"
)

const (
	natureRemoAccessTokenKey = "/natureRemoToCloudWatchFunction/natureRemoAccessToken"
	deviceIDKey              = "/natureRemoToCloudWatchFunction/deviceID"
)

var (
	CloudWatchClient cloudwatch.ClientInterface
	SSMClient        ssm.ClientInterface
	NatureRemoClient natureremo.ClientInterface
)

func Handler(ctx context.Context) error {
	deviceID, err := SSMClient.LoadSecret(deviceIDKey)
	if err != nil {
		return errors.Wrap(err, "cannot load device ID")
	}

	temperature, err := NatureRemoClient.FetchTemperature(ctx, deviceID)
	if err != nil {
		return errors.Wrap(err, "cannot fetch room temperature")
	}

	if err := CloudWatchClient.PutTemperature(time.Now(), deviceID, temperature); err != nil {
		return errors.Wrap(err, "cannot put room temperature")
	}

	return nil
}
