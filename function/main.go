package function

import (
	"context"
	"time"

	"github.com/dtan4/nature-remo-to-cloud-watch-function/aws/cloudwatch"
	"github.com/dtan4/nature-remo-to-cloud-watch-function/aws/ssm"
	"github.com/dtan4/nature-remo-to-cloud-watch-function/natureremo"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	cloudwatchapi "github.com/aws/aws-sdk-go/service/cloudwatch"
	ssmapi "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-xray-sdk-go/xray"
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
	xray.Configure(xray.Config{LogLevel: "trace"})

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

func main() {
	sess := session.New()

	cwapi := cloudwatchapi.New(sess)
	xray.AWS(cwapi.Client)
	CloudWatchClient = cloudwatch.NewClient(cwapi)

	ssmapi := ssmapi.New(sess)
	xray.AWS(ssmapi.Client)
	SSMClient = ssm.NewClient(ssmapi)

	accessToken, err := SSMClient.LoadSecret(natureRemoAccessTokenKey)
	if err != nil {
		panic(err)
	}

	NatureRemoClient = natureremo.NewClient(accessToken)

	lambda.Start(Handler)
}
