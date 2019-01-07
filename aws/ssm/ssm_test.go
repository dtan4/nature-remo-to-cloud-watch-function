package ssm

import (
	"fmt"
	"testing"

	"github.com/dtan4/nature-remo-to-cloud-watch-function/aws/mock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
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

func TestLoadSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name string
		want string
	}{
		{
			name: "/foobarbaz/foo",
			want: "abcdef",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ssmMock := mock.NewMockSSMAPI(ctrl)
			ssmMock.EXPECT().GetParameter(&ssm.GetParameterInput{
				Name:           aws.String(tc.name),
				WithDecryption: aws.Bool(true),
			}).Return(&ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{
					Name:  aws.String("/foobarbaz/foo"),
					Value: aws.String("abcdef"),
				},
			}, nil)

			client := &Client{
				api: ssmMock,
			}

			got, err := client.LoadSecret(tc.name)
			if err != nil {
				t.Fatalf("want no error, got: %s", err)
			}

			if got != tc.want {
				t.Errorf("want: %q, got: %q", tc.want, got)
			}
		})
	}
}

func TestLoadSecret_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name string
		want string
	}{
		{
			name: "/foobarbaz/foo",
			want: "cannot retrieve secret from Parameter Store: unexpected error",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ssmMock := mock.NewMockSSMAPI(ctrl)
			ssmMock.EXPECT().GetParameter(&ssm.GetParameterInput{
				Name:           aws.String(tc.name),
				WithDecryption: aws.Bool(true),
			}).Return(nil, fmt.Errorf("unexpected error"))

			client := &Client{
				api: ssmMock,
			}

			_, err := client.LoadSecret(tc.name)
			if err == nil {
				t.Fatal("want error, got: nil")
			}

			if err.Error() != tc.want {
				t.Errorf("want: %q, got: %q", tc.want, err.Error())
			}
		})
	}
}
