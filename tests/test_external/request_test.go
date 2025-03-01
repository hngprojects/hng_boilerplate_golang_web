package external_test

import (
	"testing"

	"github.com/hngprojects/hng_boilerplate_golang_web/external"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/ipstack"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/stretchr/testify/assert"
)

func TestNewExternalRequest(t *testing.T) {
	logger := getTestLogger()

	// Test with IPstack option
	ipstackReq := &ipstack.IpstackRequest{}
	extReq := request.NewExternalRequest(logger, request.WithIPstack(ipstackReq))

	assert.NotNil(t, extReq)
	assert.Equal(t, logger, extReq.Logger)
	assert.Equal(t, ipstackReq, extReq.Ipstack)
}

func TestSetupExternalRequest(t *testing.T) {
	logger := getTestLogger()
	ipstackReq := &ipstack.IpstackRequest{} // Ensure Ipstack is initialized

	extReq := request.NewExternalRequest(logger, request.WithIPstack(ipstackReq))

	// Valid request setup
	err := extReq.SetupExternalRequest(external.IpstackResolveIp, "8.8.8.8")
	assert.NoError(t, err)

	// Invalid request setup
	err = extReq.SetupExternalRequest("invalid_request", nil)
	assert.ErrorIs(t, err, request.ErrRequestNotFound)
}

func TestSendRequest_ExternalRequest(t *testing.T) {
	logger := getTestLogger()
	extReq := request.NewExternalRequest(logger, request.WithIPstack(&ipstack.IpstackRequest{}), request.WithTest(true))
	err := extReq.SetupExternalRequest("ipstack_resolve_ip", "8.8.8.8")
	assert.NoError(t, err)
	// Mock the external request
	resp, err := extReq.SendRequest()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
