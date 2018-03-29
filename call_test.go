package samson

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo_success(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, string("{\"bar\":\"foo\"}"))
	})
	server = httptest.NewServer(handler)

	client = New("token")
	client.BaseURL = server.URL

	call, err := client.NewCall("GET", "some/path", nil, nil, nil)
	assert.Nil(err)

	err = call.Do(nil)
	assert.Nil(err)
}

func TestDo_success_with_jsonbody(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, string("{\"foo\":\"bar\"}"))
	})
	server = httptest.NewServer(handler)

	client = New("token")
	client.BaseURL = server.URL

	call, err := client.NewCall("GET", "some/path", nil, nil, nil)
	assert.Nil(err)

	type Testbody struct{ Foo string }
	var testBody Testbody
	err = call.Do(&testBody)
	assert.Nil(err)
}

func TestDo_fail(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		fmt.Fprintln(w, string("{\"message\":\"error\"}"))
	})
	server = httptest.NewServer(handler)

	client = New("token")
	client.BaseURL = server.URL

	call, err := client.NewCall("GET", "some/path", nil, nil, nil)
	assert.Nil(err)

	err = call.Do(nil)
	assert.NotNil(err)
	assert.Equal(call.err.Error(), "error")
}

func TestDo_fail_malformedrequest_1(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, string("doesnt really matter"))
	})
	server = httptest.NewServer(handler)

	client = New("token")
	client.BaseURL = server.URL

	call, err := client.NewCall("GET", "some/path", nil, nil, nil)
	assert.Nil(err)

	// modify the call's request
	call.req.Method = "unknown method"

	err = call.Do(nil)
	assert.NotNil(err)
}

func TestDo_fail_malformedrequest_2(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, string("doesnt really matter"))
	})
	server = httptest.NewServer(handler)

	client = New("token")
	client.BaseURL = server.URL

	_, err := client.NewCall("<", "some/path", nil, nil, nil)
	assert.NotNil(err)
}

func TestDo_fail_malformedresponse_1(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		fmt.Fprintln(w, string("this is not a valid json"))
	})
	server = httptest.NewServer(handler)

	client = New("token")
	client.BaseURL = server.URL

	call, err := client.NewCall("GET", "some/path", nil, nil, nil)
	assert.Nil(err)

	err = call.Do(nil)
	assert.NotNil(err)
	assert.IsType(&json.SyntaxError{}, err)
}

func TestDo_fail_malformedresponse_2(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, string("this is not a valid json"))
	})
	server = httptest.NewServer(handler)

	client = New("token")
	client.BaseURL = server.URL

	call, err := client.NewCall("GET", "some/path", nil, nil, nil)
	assert.Nil(err)

	type Testbody struct{ Foo string }
	var testBody Testbody
	err = call.Do(testBody)
	assert.NotNil(err)
	assert.IsType(&json.SyntaxError{}, err)
}
