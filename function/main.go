package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dtan4/nature-remo-to-cloud-watch-function/aws/cloudwatch"
	"github.com/dtan4/nature-remo-to-cloud-watch-function/aws/ssm"
	"github.com/dtan4/nature-remo-to-cloud-watch-function/natureremo"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	cloudwatchapi "github.com/aws/aws-sdk-go/service/cloudwatch"
	ssmapi "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

const (
	natureRemoAccessTokenKey = "/natureRemoToCloudWatchFunction/natureRemoAccessToken"
	deviceIDKey              = "/natureRemoToCloudWatchFunction/deviceID"
	sentryDsnKey             = "/natureRemoToCloudWatchFunction/sentryDsn"
)

var (
	CloudWatchClient cloudwatch.ClientInterface
	SSMClient        ssm.ClientInterface
	NatureRemoClient natureremo.ClientInterface

	sentryEnabled = false
)

func init() {
	if strings.ToLower(os.Getenv("ENABLE_SENTRY")) == "true" {
		sentryEnabled = true
	}
}

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

	if sentryEnabled {
		sentryDsn, err := SSMClient.LoadSecret(ctx, sentryDsnKey)
		if err != nil {
			return errors.Wrap(err, "cannot load Sentry DSN")
		}

		if err := sentry.Init(sentry.ClientOptions{
			Dsn: sentryDsn,
			Transport: &sentry.HTTPSyncTransport{
				Timeout: 5 * time.Second,
			},

			// Release: version.Version,
			// https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html#configuration-envvars-runtime
			ServerName: os.Getenv("AWS_LAMBDA_FUNCTION_NAME"),
		}); err != nil {
			return errors.Wrap(err, "cannot initialize Sentry client")
		}
	}

	if err := RealHandler(ctx); err != nil {
		if sentryEnabled {
			sentry.CaptureException(err)
		}
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
