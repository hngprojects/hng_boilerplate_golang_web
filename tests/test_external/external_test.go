package external_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hngprojects/hng_boilerplate_golang_web/external"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"

	"github.com/stretchr/testify/assert"
)

func getTestLogger() *utility.Logger {
	return utility.NewLogger()
}

func TestSendRequest_ValidRequestObject(t *testing.T) {
	testLogger := getTestLogger()

	// Mock HTTP server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mock/path/prefix", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	data := map[string]string{"key": "value"}
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(data)

	r := &external.SendRequestObject{
		Method:       http.MethodPost,
		DecodeMethod: "json",
		Path:         ts.URL + "/mock/path",
		UrlPrefix:    "/prefix",
		Data:         data,
		Headers:      map[string]string{"Content-Type": "application/json"},
		Logger:       testLogger,
	}

	expectedResp := map[string]any{
		"message": "success",
	}

	res, err := r.SendRequest()
	fmt.Println(err)
	assert.NoError(t, err)
	fmt.Println(res)
	assert.NotNil(t, res)
	actualResp := res.(map[string]any)
	assert.Equal(t, expectedResp, actualResp)
}

func TestSendRequest_InvalidURL(t *testing.T) {
	testLogger := getTestLogger()

	r := &external.SendRequestObject{
		Method:       http.MethodGet,
		DecodeMethod: "json",
		Path:         "http://invalid-url", // Invalid URL
		Headers:      map[string]string{"Content-Type": "application/json"},
		Logger:       testLogger,
	}

	res, err := r.SendRequest()
	assert.Error(t, err)
	assert.EqualError(t, err, external.ErrNoSuchHost.Error())
	assert.Nil(t, res)
}
