package samson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// DeployGroupService service
type DeployGroupService service

// DeployGroup model
type DeployGroup struct {
	ID            *int       `json:"id,omitempty"`
	Name          *string    `json:"name,omitempty"`
	EnvironmentID *int       `json:"environment_id,omitempty"`
	EnvValue      *string    `json:"env_value,omitempty"`
	Permalink     *string    `json:"permalink,omitempty"`
	VaultServerID *int       `json:"vault_server_id,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	DeleteAt      *time.Time `json:"deleted_at,omitempty"`
}

// List returns all deploy groups
func (service *DeployGroupService) List() ([]*DeployGroup, *Call, error) {
	path := "/deploy_groups.json"
	method := "GET"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return nil, call, err
	}

	type response struct {
		DeployGroups []*DeployGroup `json:"deploy_groups,omitempty"`
	}
	var res response
	err = call.Do(&res)
	if err != nil {
		return nil, call, err
	}

	return res.DeployGroups, call, nil
}

// Get returns a single deploy group resource
func (service *DeployGroupService) Get(id int) (*DeployGroup, *Call, error) {
	path := fmt.Sprintf("/deploy_groups/%d.json", id)
	method := "GET"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return nil, call, err
	}

	var dg *DeployGroup
	err = call.Do(&dg)
	if err != nil {
		return nil, call, err
	}

	return dg, call, nil
}

// Upsert updates or creates a new deploy group resource
func (service *DeployGroupService) Upsert(dg *DeployGroup) (*DeployGroup, *Call, error) {
	bytesArray, _ := json.Marshal(dg)

	var path string
	var method string

	if dg.CreatedAt != nil {
		path = fmt.Sprintf("/deploy_groups/%d.json", *dg.ID)
		method = "PUT"
	} else {
		path = "/deploy_groups.json"
		method = "POST"
	}

	call, err := service.s.NewCall(method, path, nil, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return nil, call, err
	}

	err = call.Do(&dg)
	if err != nil {
		return nil, call, err
	}

	return dg, call, nil
}

// Delete deletes a single deploy group resource
func (service *DeployGroupService) Delete(id int) (*Call, error) {
	path := fmt.Sprintf("/deploy_groups/%d.json", id)
	method := "DELETE"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return call, err
	}

	return call, call.Do(nil)
}
