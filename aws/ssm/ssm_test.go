package ssm

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type mockSSMAPI struct {
	ssmiface.SSMAPI
	getParameterWithContextFunc func(ctx aws.Context, input *ssm.GetParameterInput, opts ...request.Option) (*ssm.GetParameterOutput, error)
}

func (m *mockSSMAPI) GetParameterWithContext(ctx aws.Context, input *ssm.GetParameterInput, opts ...request.Option) (*ssm.GetParameterOutput, error) {
	return m.getParameterWithContextFunc(ctx, input, opts...)
}

func TestNewClient(t *testing.T) {
	api := &mockSSMAPI{}
	client := NewClient(api)

	if client.api != api {
		t.Error("api does not match")
	}
}

func TestLoadSecret(t *testing.T) {
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
			ctx := context.Background()

			client := &Client{
				api: &mockSSMAPI{
					getParameterWithContextFunc: func(ctx aws.Context, input *ssm.GetParameterInput, opts ...request.Option) (*ssm.GetParameterOutput, error) {
						return &ssm.GetParameterOutput{
							Parameter: &ssm.Parameter{
								Name:  aws.String("/foobarbaz/foo"),
								Value: aws.String("abcdef"),
							},
						}, nil
					},
				},
			}

			got, err := client.LoadSecret(ctx, tc.name)
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
	testcases := []struct {
		name string
		err  string
		want string
	}{
		{
			name: "/foobarbaz/foo",
			err:  "unexpected error",
			want: "cannot retrieve secret from Parameter Store: unexpected error",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			client := &Client{
				api: &mockSSMAPI{
					getParameterWithContextFunc: func(ctx aws.Context, input *ssm.GetParameterInput, opts ...request.Option) (*ssm.GetParameterOutput, error) {
						return nil, fmt.Errorf(tc.err)
					},
				},
			}

			_, err := client.LoadSecret(ctx, tc.name)
			if err == nil {
				t.Fatal("want error, got: nil")
			}

			if err.Error() != tc.want {
				t.Errorf("want: %q, got: %q", tc.want, err.Error())
			}
		})
	}
}
