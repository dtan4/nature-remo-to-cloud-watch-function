package main

import (
	"context"
	"net/http"
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

func RealHandler(ctx context.Context) error {
	deviceID, err := SSMClient.LoadSecret(ctx, deviceIDKey)
	if err != nil {
		return errors.Wrap(err, "cannot load device ID")
	}

	temperature, err := NatureRemoClient.FetchTemperature(ctx, deviceID)
	if err != nil {
		return errors.Wrap(err, "cannot fetch room temperature")
	}

	if err := CloudWatchClient.PutTemperature(ctx, time.Now(), deviceID, temperature); err != nil {
		return errors.Wrap(err, "cannot put room temperature")
	}

	return nil
}

func Handler(ctx context.Context) error {
	sess := session.Must(session.NewSession())

	cwapi := cloudwatchapi.New(sess)
	xray.AWS(cwapi.Client)
	CloudWatchClient = cloudwatch.NewClient(cwapi)

	ssmapi := ssmapi.New(sess)
	xray.AWS(ssmapi.Client)
	SSMClient = ssm.NewClient(ssmapi)

	accessToken, err := SSMClient.LoadSecret(ctx, natureRemoAccessTokenKey)
	if err != nil {
		return errors.Wrap(err, "cannot load Nature Remo access token")
	}

	NatureRemoClient = natureremo.NewClient(accessToken, xray.Client(http.DefaultClient))

	return RealHandler(ctx)
}

func main() {
	lambda.Start(Handler)
}
