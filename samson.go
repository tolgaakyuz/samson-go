package samson

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Samson model
type Samson struct {
	token       string
	QueryParams map[string]string
	Headers     map[string]string
	BaseURL     string

	Projects *ProjectService
	Stages   *StageService
}

type service struct {
	s *Samson
}

// New returns a Samson client
func New(token string) *Samson {
	s := &Samson{
		token:       token,
		QueryParams: map[string]string{},
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
			"Content-Type":  "application/json",
			"User-Agent":    fmt.Sprintf("sdk samson-go/%s", Version),
		},
		BaseURL: "https://localhost:9080",
	}

	s.Projects = &ProjectService{s: s}
	s.Stages = &StageService{s: s}

	return s
}

// NewCall creates a new api call object
func (s *Samson) NewCall(method, path string, queryParams, headers map[string]string, body io.Reader) (*Call, error) {
	// prepare the request url
	u, err := url.Parse(s.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path

	// merge queryParams
	query := url.Values{}
	for key, value := range s.QueryParams {
		query.Set(key, value)
	}
	if queryParams != nil {
		for key, value := range queryParams {
			query.Set(key, value)
		}
	}

	// merge headers
	if headers == nil {
		headers = map[string]string{}
	}
	for key, value := range s.Headers {
		headers[key] = value
	}

	call := &Call{
		client:      http.DefaultClient,
		method:      method,
		url:         u,
		queryParams: query,
		headers:     headers,
		body:        body,
	}

	err = call.prepareRequest()
	if err != nil {
		return nil, err
	}

	return call, nil
}

func String(s string) *string {
	return &s
}

func Int(i int) *int {
	return &i
}
