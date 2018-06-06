package samson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// EnvironmentService service
type EnvironmentService service

// Environment model
type Environment struct {
	ID         *int    `json:"id,omitempty"`
	Name       *string `json:"name,omitempty"`
	Production *string `json:"production,omitempty"`
}

// List returns all environments
func (service *EnvironmentService) List() ([]*Environment, *Call, error) {
	path := "/environments.json"
	method := "GET"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return nil, call, err
	}

	type response struct {
		Environments []*Environment `json:"environments,omitempty"`
	}
	var res response
	err = call.Do(&res)
	if err != nil {
		return nil, call, err
	}

	return res.Environments, call, nil
}

// Get returns a single environment resource
func (service *EnvironmentService) Get(id int) (*Environment, *Call, error) {
	path := fmt.Sprintf("/environments/%d.json", id)
	method := "GET"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return nil, call, err
	}

	var environment Environment
	err = call.Do(&environment)
	if err != nil {
		return nil, call, err
	}

	return &environment, call, nil
}

// Upsert updates or creates a new environment resource
func (service *EnvironmentService) Upsert(environment *Environment) (*Environment, *Call, error) {
	bytesArray, _ := json.Marshal(environment)

	var path string
	var method string

	if environment.ID != nil {
		path = fmt.Sprintf("/environments/%d.json", *environment.ID)
		method = "PUT"
	} else {
		path = "/environments.json"
		method = "POST"
	}

	call, err := service.s.NewCall(method, path, nil, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return nil, call, err
	}

	err = call.Do(&environment)
	if err != nil {
		return nil, call, err
	}

	return environment, call, nil
}

// Delete deletes a sinlge environment resource
func (service *EnvironmentService) Delete(id int) (*Call, error) {
	path := fmt.Sprintf("/environments/%d.json", id)
	method := "DELETE"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return call, err
	}

	return call, call.Do(nil)
}
