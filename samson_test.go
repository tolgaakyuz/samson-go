package samson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	server *httptest.Server
	client *Samson
	token  = "5e4ee7a30136a56db251c54ea9330566324ea40dd8f6de5c3a112a0654a267de"
)

func readTestData(fileName string) string {
	path := "testdata/" + fileName
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(content)
}

func checkHeaders(req *http.Request, assert *assert.Assertions) {
	assert.Equal("Bearer "+token, req.Header.Get("Authorization"))
	assert.Equal("application/json", req.Header.Get("Content-Type"))
}

func setup() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixture := strings.Replace(r.URL.Path, "/", "-", -1)
		fixture = strings.TrimLeft(fixture, "-")
		var path string

		if e := r.URL.Query().Get("error"); e != "" {
			path = "testdata/error-" + e + ".json"
		} else {
			if r.Method == "GET" {
				path = "testdata/" + fixture + ".json"
			}

			if r.Method == "POST" {
				path = "testdata/" + fixture + "-new.json"
			}

			if r.Method == "PUT" {
				path = "testdata/" + fixture + "-updated.json"
			}
		}

		file, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		fmt.Fprintln(w, string(file))
		return
	})

	server = httptest.NewServer(handler)

	client = New(token)
	client.BaseURL = server.URL
}

func teardown() {
	server.Close()
	client = nil
}

func TestSamsonNew(t *testing.T) {
	assert := assert.New(t)

	client = New(token)
	assert.IsType(Samson{}, *client)
	assert.IsType(ProjectService{}, *client.Projects)
	assert.IsType(StageService{}, *client.Stages)
	assert.Equal("http://localhost:9080", client.BaseURL)
	assert.Equal(token, client.token)
	assert.Equal(fmt.Sprintf("Bearer %s", token), client.Headers["Authorization"])
	assert.Equal("application/json", client.Headers["Content-Type"])
	assert.Equal(fmt.Sprintf("sdk samson-go/%s", Version), client.Headers["User-Agent"])
}

func TestNewCall(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	method := "GET"
	path := "/some/path"
	queryParams := map[string]string{
		"foo": "bar",
		"faz": "zoo",
	}

	query := url.Values{}
	for key, value := range queryParams {
		query.Add(key, value)
	}
	query.Add("foo2", "bar2")

	expectedURL, _ := url.Parse(client.BaseURL)
	expectedURL.Path = path
	expectedURL.RawQuery = query.Encode()

	// set the default query params
	client.QueryParams = map[string]string{
		"foo2": "bar2",
	}

	call, err := client.NewCall(method, path, queryParams, nil, nil)
	assert.Nil(err)
	assert.Equal(call.req.Header.Get("Authorization"), "Bearer "+token)
	assert.Equal(call.req.Header.Get("Content-Type"), "application/json")
	assert.Equal(call.req.Method, method)
	assert.Equal(call.req.URL.String(), expectedURL.String())

	method = "POST"
	type RequestBody struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	bodyData := RequestBody{
		Name: "test",
		Age:  10,
	}
	body, _ := json.Marshal(bodyData)
	call, err = client.NewCall(method, path, queryParams, nil, bytes.NewReader(body))
	assert.Nil(err)
	assert.Equal(call.req.Header.Get("Authorization"), "Bearer "+token)
	assert.Equal(call.req.Header.Get("Content-Type"), "application/json")
	assert.Equal(call.req.Method, method)
	assert.Equal(call.req.URL.String(), expectedURL.String())
	defer call.req.Body.Close()
	var requestBody RequestBody
	err = json.NewDecoder(call.req.Body).Decode(&requestBody)
	assert.Nil(err)
	assert.Equal(requestBody, bodyData)
}

func TestNewCall_error_malformedurl(t *testing.T) {
	assert := assert.New(t)

	client = New("token")
	client.BaseURL = "^http://localhost"

	_, err := client.NewCall("GET", "some/path", nil, nil, nil)
	assert.NotNil(err)
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	s := "test string"
	assert.Equal(s, *String(s))
}

func TestInt(t *testing.T) {
	assert := assert.New(t)

	i := 1
	assert.Equal(i, *Int(i))
}
