package cloudwatch

import (
	"testing"

	"github.com/dtan4/nature-remo-to-cloud-watch-function/aws/mock"

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
