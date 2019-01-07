package ssm

import (
	"testing"

	"github.com/dtan4/nature-remo-to-cloud-watch-function/aws/mock"

	"github.com/golang/mock/gomock"
)

//go:generate mockgen -source $GOPATH/pkg/mod/github.com/aws/aws-sdk-go@v1.16.14/service/ssm/ssmiface/interface.go -destination ../mock/ssm.go -package mock

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ssmMock := mock.NewMockSSMAPI(ctrl)

	got := NewClient(ssmMock)
	if got == nil {
		t.Error("want object, got nil")
	}
}
