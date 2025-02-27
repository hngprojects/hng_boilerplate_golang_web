package request

import (
	"errors"

	"github.com/hngprojects/hng_boilerplate_golang_web/external"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/ipstack"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/mocks"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

// ExternalRequest holds fields for communicating with external services
//
// Create a new ExternalRequest using NewExternalRequest,
// pass in logger and request config options(WithTest, WithIpstack) to configure ExternalRequest
//
// Call SetupExternalRequest before SendRequest to send requests to external-services
type ExternalRequest struct {
	Logger      *utility.Logger
	Test        bool
	Ipstack     *ipstack.IpstackRequest
	requestName string
}

type sendExtReqOpts func(*ExternalRequest) *ExternalRequest

var (
	ErrRequestNotFound = errors.New("request not found")
)

func WithTest(test bool) sendExtReqOpts {
	return func(sr *ExternalRequest) *ExternalRequest {
		sr.Test = test
		return sr
	}
}

func WithIPstack(ir *ipstack.IpstackRequest) sendExtReqOpts {
	return func(sr *ExternalRequest) *ExternalRequest {
		sr.Ipstack = ir
		sr.Ipstack.Logger = sr.Logger
		return sr
	}
}

func NewExternalRequest(logger *utility.Logger, opts ...sendExtReqOpts) *ExternalRequest {
	sendExtRequest := ExternalRequest{
		Logger: logger,
	}

	for _, op := range opts {
		op(&sendExtRequest)
	}
	return &sendExtRequest
}

func (sr *ExternalRequest) SetupExternalRequest(name string, data any) error {
	switch name {
	case external.IpstackResolveIp:
		err := sr.Ipstack.SetupExternalRequest(name, data)
		if err != nil {
			return err
		}
	default:
		return ErrRequestNotFound
	}
	return nil
}
func (sr *ExternalRequest) SendRequest() (any, error) {
	if !sr.Test {
		switch sr.requestName {
		case external.IpstackResolveIp:
			return sr.Ipstack.SendRequest()
		default:
			return nil, ErrRequestNotFound
		}
	} else {
		switch sr.requestName {
		case external.IpstackResolveIp:
			mer := mocks.ExternalRequest{Logger: sr.Logger, Test: true}
			resp, err := mer.SendExternalRequest(sr.requestName, sr.Ipstack.RequestData)
			return resp, err
		}

	}
	return nil, nil
}
