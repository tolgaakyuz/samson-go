package samson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// StageService service
type StageService service

// Stage model
type Stage struct {
	ID                                     *int       `json:"id,omitempty"`
	Name                                   *string    `json:"name,omitempty"`
	ProjectID                              *int       `json:"project_id,omitempty"`
	NotifyEmailAddess                      *string    `json:"notify_email_address,omitempty"`
	Order                                  *int       `json:"order,omitempty"`
	Confirm                                *bool      `json:"confirm,omitempty"`
	DatadogTags                            *string    `json:"datadog_tags,omitempty"`
	DatadogMonitorIds                      *string    `json:"datadog_monitor_ids,omitempty"`
	UpdateGithubPullRequests               *bool      `json:"update_github_pull_requests,omitempty"`
	DeployOnRelease                        *bool      `json:"deploy_on_release,omitempty"`
	CommentOnZendeskTickets                *bool      `json:"comment_on_zendesk_tickets,omitempty"`
	Production                             *bool      `json:"production,omitempty"`
	Permalink                              *string    `json:"permalink,omitempty"`
	UseGithubDeploymentAPI                 *bool      `json:"use_github_deployment_api,omitempty"`
	Dashboard                              *string    `json:"dashboard,omitempty"`
	EmailCommitersOnAutomatedDeployFailure *bool      `json:"email_committers_on_automated_deploy_failure,omitempty"`
	StaticEmailsOnAutomatedDeployFailure   *string    `json:"static_emails_on_automated_deploy_failure,omitempty"`
	JenkinsJobNames                        *string    `json:"jenkins_job_names,omitempty"`
	NextStageIds                           []*string  `json:"next_stage_ids,omitempty"`
	NoCodeDeployed                         *bool      `json:"no_code_deployed,omitempty"`
	DockerBinaryPluginEnabled              *bool      `json:"docker_binary_plugin_enabled,omitempty"`
	IsTemplate                             *bool      `json:"is_template,omitempty"`
	NotifyAirbrake                         *bool      `json:"notify_airbrake,omitempty"`
	TamplateStageID                        *int       `json:"template_stage_id,omitempty"`
	JenkinsEmailCommitters                 *bool      `json:"jenkins_email_committers,omitempty"`
	Kubernetes                             *bool      `json:"kubernetes,omitempty"`
	RunInParallel                          *bool      `json:"run_in_parallel,omitempty"`
	JenkinsBuildParams                     *bool      `json:"jenkins_build_params,omitempty"`
	CancelQueuedDeploys                    *bool      `json:"cancel_queued_deploys,omitempty"`
	NoReferenceSelection                   *bool      `json:"no_reference_selection,omitempty"`
	PeriodicalDeploy                       *bool      `json:"periodical_deploy,omitempty"`
	BuildsInEnvironment                    *bool      `json:"builds_in_environment,omitempty"`
	BlockOnGCRVulnerabilities              *bool      `json:"block_on_gcr_vulnerabilities,omitempty"`
	NotifyAssertible                       *bool      `json:"notify_assertible,omitempty"`
	CreatedAt                              *time.Time `json:"created_at,omitempty"`
	UpdatedAt                              *time.Time `json:"updated_at,omitempty"`
	DeleteAt                               *time.Time `json:"deleted_at,omitempty"`
}

// List returns all stages
func (service *StageService) List() ([]*Stage, *Call, error) {
	path := "/stages.json"
	method := "GET"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return nil, call, err
	}

	type StagesResponse struct {
		Stages []*Stage `json:"stages,omitempty"`
	}
	var stagesResponse StagesResponse
	err = call.Do(&stagesResponse)
	if err != nil {
		return nil, call, err
	}

	return stagesResponse.Stages, call, nil
}

// Get returns a single stage resource
func (service *StageService) Get(id int) (*Stage, *Call, error) {
	path := fmt.Sprintf("/stages/%d.json", id)
	method := "GET"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return nil, call, err
	}

	var stage *Stage
	err = call.Do(&stage)
	if err != nil {
		return nil, call, err
	}

	return stage, call, nil
}

// Upsert updates or creates a new stage resource
func (service *StageService) Upsert(stage *Stage) (*Stage, *Call, error) {
	bytesArray, _ := json.Marshal(stage)

	var path string
	var method string

	if stage.CreatedAt != nil {
		path = fmt.Sprintf("/stages/%d.json", *stage.ID)
		method = "PUT"
	} else {
		path = "/stages.json"
		method = "POST"
	}

	call, err := service.s.NewCall(method, path, nil, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return nil, call, err
	}

	err = call.Do(&stage)
	if err != nil {
		return nil, call, err
	}

	return stage, call, nil
}

// Delete deletes a single stage resource
func (service *StageService) Delete(id int) (*Call, error) {
	path := fmt.Sprintf("/stages/%d.json", id)
	method := "DELETE"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return call, err
	}

	return call, call.Do(nil)
}
