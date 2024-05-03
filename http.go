package jwtsdk

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/correctinho/correct-mlt-go/qlog"
	utilsdk "github.com/correctinho/correct-util-sdk-go"
	stg "github.com/correctinho/correct-util-sdk-go/stg"
)

// HTTP - cliente http
type HTTP struct {
	sync.RWMutex
	Context interface{}
	Host    string
}

// newHTTP - novo objeto http
func newHTTP(ctx interface{}) HTTP {
	return HTTP{
		Context: ctx,
		Host:    os.Getenv("JWT_URL"),
	}
}

// NewHTTPClient - novo cliente http
func NewHTTPClient() http.Client {
	return http.Client{Timeout: 120 * time.Second}
}

// Do - envio de requisicao
func (h *HTTP) Do(method string, path string, payload interface{}, out interface{}, logging bool) *JwtError {
	logger := qlog.NewProduction(h.Context)
	defer logger.Sync()

	setLogging := func(url, path, msg string, bd []byte, statusCode int) {
		logger.InfoJSON(msg, string(bd), qlog.LoggerExtras{Key: "http_request", Value: map[string]interface{}{
			"host":   url,
			"method": method,
			"path":   path,
			"url":    url + path,
			"status": statusCode,
		}})
	}

	structToMap := func(in interface{}, out map[string]interface{}) error {
		h.Lock()
		defer h.Unlock()
		data, e := json.Marshal(in)
		if e != nil {
			return e
		}
		e = json.Unmarshal(data, &out)
		if e != nil {
			return e
		}
		return nil
	}

	var body []byte

	if payload != nil {
		data := make(map[string]interface{}, 0)

		if e := structToMap(payload, data); e != nil {
			logger.Error(e.Error())
			return &ErrServiceUnavailable
		}

		var e error
		body, e = json.Marshal(data)
		if e != nil {
			logger.Error(e.Error())
			return &ErrServiceUnavailable
		}

		setLogging(h.Host, path, "REQUEST SENT", []byte(""), 0)

		if _, ok := os.LookupEnv("GO_DEBUG"); ok {
			var prettyJSON bytes.Buffer
			json.Indent(&prettyJSON, body, "", "\t")
			println("======= HTTP REQUEST =========")
			fmt.Printf("%v\n\n", prettyJSON.String())
		}
	}

	request, e := http.NewRequest(method, h.Host+path, bytes.NewReader(body))
	if e != nil {
		logger.Error(e.Error())
		return &ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(request.Context(), (60 * time.Second))
	defer cancel()

	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")

	client := NewHTTPClient()
	response, e := client.Do(request)
	if e != nil {
		logger.Error(e.Error())
		return &ErrServiceUnavailable
	}

	defer response.Body.Close()
	respBody, e := io.ReadAll(response.Body)
	if e != nil {
		logger.Error(e.Error())
		setLogging(h.Host, path, "RESPONSE RECEIVED", []byte(e.Error()), response.StatusCode)
		return &ErrServiceUnavailable
	}

	contentEncoding := []byte(response.Header.Get("Content-Encoding"))
	var responseBody []byte
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		reader, e := gzip.NewReader(response.Body)
		if e != nil {
			logger.Error(e.Error())
			return &ErrServiceUnavailable
		}
		defer reader.Close()
		responseBody, e = io.ReadAll(reader)
		if e != nil {
			logger.Error(e.Error())
			return &ErrServiceUnavailable
		}
	} else {
		responseBody = respBody
	}

	if _, ok := os.LookupEnv("GO_DEBUG"); ok {
		var prettyJSON bytes.Buffer
		json.Indent(&prettyJSON, responseBody, "", "\t")
		println("======= HTTP RESPONSE =========")
		fmt.Printf("%v\n\n", prettyJSON.String())
	}

	if response.StatusCode == 500 {
		return &ErrServiceUnavailable
	}
	if response.StatusCode >= 300 {
		logger.Error(string(responseBody))
		return &ErrServiceUnavailable
	}
	if utilsdk.IsNil(out) {
		return nil
	}
	if e := json.Unmarshal(responseBody, &out); e != nil {
		logger.Error(e.Error())
		return &ErrServiceUnavailable
	}
	return nil
}

// tryUnmarshal - Função que converte um []byte em um json
func (h *HTTP) tryUnmarshal(buffer []byte) string {
	out := struct {
		Response map[string]interface{} `json:"response"`
		Status   int                    `json:"status"`
	}{}
	tmp := &[]struct {
		Response string `json:"response"`
		Status   int    `json:"status"`
	}{}
	value := string(buffer)
	if len(value) > 0 {
		if e := json.Unmarshal(buffer, &tmp); e != nil {
			println(e.Error())
			return string(buffer)
		}
		if len(*tmp) > 0 {
			obj := *tmp
			if e := json.Unmarshal([]byte(obj[0].Response), &out.Response); e != nil {
				println(e.Error())
				return string(buffer)
			}
		}
	}
	return stg.ToJSON(out)
}
