package external_test

import (
	"testing"

	"github.com/hngprojects/hng_boilerplate_golang_web/external"
	"github.com/hngprojects/hng_boilerplate_golang_web/external/ipstack"
	"github.com/stretchr/testify/assert"
)

// Test SetupExternalRequest
func TestSetupExternalRequest_IPStack(t *testing.T) {
	logger := getTestLogger()

	ipReq := ipstack.IpstackRequest{
		Logger: logger,
	}

	err := ipReq.SetupExternalRequest(external.IpstackResolveIp, "192.168.1.1")

	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", ipReq.RequestData)
	assert.Equal(t, "GET", ipReq.Method)
	assert.Equal(t, 200, ipReq.SuccessCode)
}

// Test SendRequest (successful response)
func TestSendRequest_IPstack(t *testing.T) {
	logger := getTestLogger()

	ipReq := ipstack.IpstackRequest{
		Logger: logger,
	}
	ipReq.SetupExternalRequest(external.IpstackResolveIp, "192.168.1.1")
	response, err := ipReq.SendRequest()

	assert.Error(t, err)
	assert.ErrorIs(t, err, ipstack.ErrInvalidIPStackResolveIPResponse) //Because IPstack.Key is not set in env
	assert.Equal(t, "", response.Ip)
}
