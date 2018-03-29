package samson

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Call represents an api call
type Call struct {
	client      *http.Client
	url         *url.URL
	method      string
	queryParams url.Values
	headers     map[string]string
	body        io.Reader
	req         *http.Request
	res         *http.Response
	err         *ErrorResponse
}

func (call *Call) prepareRequest() error {
	call.url.RawQuery = call.queryParams.Encode()

	req, err := http.NewRequest(call.method, call.url.String(), call.body)
	if err != nil {
		return err
	}

	for key, value := range call.headers {
		req.Header.Set(key, value)
	}

	call.req = req

	return nil
}

// Do makes the call
func (call *Call) Do(v interface{}) error {
	res, err := call.client.Do(call.req)
	if err != nil {
		return err
	}

	call.res = res

	if call.res.StatusCode >= 200 && call.res.StatusCode < 400 {
		if v != nil {
			defer call.res.Body.Close()
			err = json.NewDecoder(call.res.Body).Decode(v)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return call.handleError()
}

func (call *Call) handleError() error {
	var e ErrorResponse
	defer call.res.Body.Close()
	err := json.NewDecoder(call.res.Body).Decode(&e)
	if err != nil {
		return err
	}

	call.err = &e

	return e
}
