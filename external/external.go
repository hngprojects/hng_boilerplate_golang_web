package external

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"io"
	"net/http"
	"strconv"

	"github.com/elliotchance/phpserialize"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

var (
	ResponseCode int
	ResponseBody string

	// serializers
	JsonDecodeMethod    = "json"
	PhpSerializerMethod = "phpserializer"

	// requests
	IpstackResolveIp                = "ipstack_resolve_ip"
	ErrNoSuchHost                   = errors.New("no such host")
	ErrFailedToReadReqBody          = errors.New("failed to read request body")
	ErrFaildToUnmarshalJsonResponse = errors.New("failed to unmarshal json")
	ErrFailedToCreateHttpRequest    = errors.New("failed to create http request")
	ErrFailedToEncodeData           = errors.New("faild to encode request data")
	ErrFailedToUnmarshalPhpResponse = errors.New("failed to unmarshal php response")
)

type SendRequestObject struct {
	Name         string
	Logger       *utility.Logger
	Path         string
	Method       string
	Headers      map[string]string
	SuccessCode  int
	Data         interface{}
	DecodeMethod string
	UrlPrefix    string
}

// SendRequest is a general external function for sending requests to
// external services needed by hng_boilerplate
func (r *SendRequestObject) SendRequest() (any, error) {
	var (
		data   = r.Data
		logger = r.Logger
		name   = r.Name
		err    error
	)

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(data)
	if err != nil {
		logger.Error("encoding error", name, err.Error())
		return nil, ErrFailedToEncodeData
	}

	logger.Info("before prefix", name, r.Path, data, buf)
	if r.UrlPrefix != "" {
		r.Path += r.UrlPrefix
	}
	logger.Info("after prefix", name, r.Path, data, buf)

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, r.Path, buf)
	if err != nil {
		logger.Error("request creation error", name, err.Error())
		return nil, ErrFailedToCreateHttpRequest
	}

	for key, value := range r.Headers {
		req.Header.Add(key, value)
	}

	logger.Info("request", name, r.Path, r.Method, r.Headers)

	res, err := client.Do(req)
	if err != nil {
		logger.Error("client do", name, err.Error())
		return nil, ErrNoSuchHost
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("reading body error", name, err.Error())
		return nil, ErrFailedToReadReqBody
	}
	var response any
	if r.DecodeMethod != PhpSerializerMethod {

		err = json.Unmarshal(body, &response)
		if err != nil {
			logger.Error("json decoding error", name, err.Error())
			return nil, ErrFaildToUnmarshalJsonResponse
		}
	}

	logger.Info("response body", name, r.Path, string(body))

	if r.DecodeMethod == PhpSerializerMethod {
		err := phpserialize.Unmarshal(body, &response)
		if err != nil {
			logger.Error("php serializer decoding error", name, err.Error())
			return nil, ErrFailedToUnmarshalPhpResponse
		}
	}

	defer res.Body.Close()
	ResponseCode = res.StatusCode

	if res.StatusCode == r.SuccessCode {
		return response, nil
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("external requests error for request %v, code %v", name, strconv.Itoa(res.StatusCode))
	}

	return response, nil
}
