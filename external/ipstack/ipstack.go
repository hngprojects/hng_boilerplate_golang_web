package ipstack

import (
	"errors"
	"fmt"

	"github.com/hngprojects/hng_boilerplate_golang_web/external"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type IpstackRequest struct {
	RequestData  string
	Headers      map[string]string
	Logger       *utility.Logger
	Name         string
	Path         string
	Method       string
	DecodeMethod string
	UrlPrefix    string
	SuccessCode  int
}

var (
	ErrFailedToResolveIP               error = errors.New("ipstack failed to resolve Ip stack")
	ErrInvalidIPStackResolveIPResponse error = errors.New("invalid response type, for ipstack")
)

// NewIpstackRequest creates an external.ExternalRequest interface for making IpStack Requests
func NewIpstackRequest(name, path, method, decodeMethod string, successCode int, ip string) IpstackRequest {
	return IpstackRequest{
		Name:         name,
		Path:         path,
		Method:       method,
		DecodeMethod: decodeMethod,
		SuccessCode:  successCode,
		RequestData:  ip,
	}
}

func (ir *IpstackRequest) SetupExternalRequest(name string, data any) error {
	config := config.Setup(ir.Logger, "./../../app")

	key := config.IPStack.Key
	logger := ir.Logger
	idata := data

	ip, ok := idata.(string)
	if !ok {
		logger.Error("ipstack resolve ip", idata, "request data format error")
		return fmt.Errorf("request data format error")
	}

	ir.Name = name
	ir.UrlPrefix = "/" + ip + "?access_key=" + key
	ir.Path = fmt.Sprintf("%v", config.IPStack.BaseUrl)
	ir.Method = "GET"
	ir.SuccessCode = 200
	ir.DecodeMethod = external.JsonDecodeMethod
	ir.RequestData = ip

	return nil
}

func (ir *IpstackRequest) SendRequest() (IPStackResolveIPResponse, error) {

	reqObj := external.SendRequestObject{
		Name:         ir.Name,
		Logger:       ir.Logger,
		Path:         ir.Path,
		Method:       ir.Method,
		Headers:      ir.Headers,
		SuccessCode:  ir.SuccessCode,
		Data:         ir.RequestData,
		DecodeMethod: ir.DecodeMethod,
		UrlPrefix:    ir.UrlPrefix,
	}
	response, err := reqObj.SendRequest()
	if err != nil {
		return IPStackResolveIPResponse{}, ErrFailedToResolveIP
	}
	resolveIPResponse, ok := response.(IPStackResolveIPResponse)
	if !ok {
		return IPStackResolveIPResponse{}, ErrInvalidIPStackResolveIPResponse
	}
	return resolveIPResponse, nil
}
