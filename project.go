package samson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// ProjectService service
type ProjectService service

// Project model
type Project struct {
	ID                                     *int                   `json:"id,omitempty"`
	Name                                   *string                `json:"name,omitempty"`
	Description                            *string                `json:"description,omitempty"`
	RepositoryURL                          *string                `json:"repository_url,omitempty"`
	ReleaseBranch                          *string                `json:"release_branch,omitempty"`
	Permalink                              *string                `json:"permalink,omitempty"`
	EnvironmentVariableAttributes          []*EnvironmentVariable `json:"environment_variables_attributes,omitempty"`
	IncludeNewDeployGroups                 *bool                  `json:"include_new_deploy_groups,omitempty"`
	DockerReleaseBranch                    *string                `json:"docker_release_branch,omitempty"`
	DockerImageBuildingDisabled            *bool                  `json:"docker_image_building_disabled,omitempty"`
	Dockerfiles                            *string                `json:"dockerfiles,omitempty"`
	BuildWithGCB                           *bool                  `json:"build_with_gcb,omitempty"`
	ShowGCBVulnerabilities                 *bool                  `json:"show_gcr_vulnerabilities,omitempty"`
	KubernetesAllowWritingToRootFilesystem *bool                  `json:"kubernetes_allow_writing_to_root_filesystem,omitempty"`
	ReleaseSource                          *string                `json:"release_source,omitempty"`
	BuildCommandID                         *int                   `json:"build_command_id,omitempty"`
	Dashboard                              *string                `json:"dashboard,omitempty"`
	RepositoryPath                         *string                `json:"repository_path,omitempty"`
	Owner                                  *string                `json:"owner,omitempty"`
	CreatedAt                              *time.Time             `json:"created_at,omitempty"`
	UpdatedAt                              *time.Time             `json:"updated_at,omitempty"`
}

// EnvironmentVariable model for projects
type EnvironmentVariable struct {
	Name           *string `json:"name,omitempty"`
	Value          *string `json:"value,omitempty"`
	ScopeTypeAndID *string `json:"scope_type_and_id,omitempty"`
}

// List returns all projects
func (service *ProjectService) List() ([]*Project, *Call, error) {
	path := "/projects.json"
	method := "GET"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return nil, call, err
	}

	type ProjectsResponse struct {
		Projects []*Project `json:"projects,omitempty"`
	}
	var projectsResponse ProjectsResponse
	err = call.Do(&projectsResponse)
	if err != nil {
		return nil, call, err
	}

	return projectsResponse.Projects, call, nil
}

// Get returns a single project resource
func (service *ProjectService) Get(id int) (*Project, *Call, error) {
	path := fmt.Sprintf("/projects/%d.json", id)
	method := "GET"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return nil, call, err
	}

	var project Project
	err = call.Do(&project)
	if err != nil {
		return nil, call, err
	}

	return &project, call, nil
}

// Upsert updates or creates a new project resource
func (service *ProjectService) Upsert(project *Project) (*Project, *Call, error) {
	bytesArray, _ := json.Marshal(project)

	var path string
	var method string

	if project.CreatedAt != nil {
		path = fmt.Sprintf("/projects/%d.json", *project.ID)
		method = "PUT"
	} else {
		path = "/projects.json"
		method = "POST"
	}

	call, err := service.s.NewCall(method, path, nil, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return nil, call, err
	}

	err = call.Do(&project)
	if err != nil {
		return nil, call, err
	}

	return project, call, nil
}

// Delete deletes a sinlge project resource
func (service *ProjectService) Delete(id int) (*Call, error) {
	path := fmt.Sprintf("/projects/%d.json", id)
	method := "DELETE"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return call, err
	}

	return call, call.Do(nil)
}
