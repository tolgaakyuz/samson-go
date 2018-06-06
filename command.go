package samson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CommandService service
type CommandService service

// Command model
type Command struct {
	ID        *int    `json:"id,omitempty"`
	Command   *string `json:"command,omitempty"`
	ProjectID *string `json:"project_id,omitempty"`
}

// IsGlobal returns whether the command is defined as global or not
// Global commands are available for all projects
func (c *Command) IsGlobal() bool {
	return c.ProjectID == nil
}

// List returns all commands
func (service *CommandService) List() ([]*Command, *Call, error) {
	path := "/commands.json"
	method := "GET"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return nil, call, err
	}

	type response struct {
		Commands []*Command `json:"commands,omitempty"`
	}
	var res response
	err = call.Do(&res)
	if err != nil {
		return nil, call, err
	}

	return res.Commands, call, nil
}

// Get returns a single command resource
func (service *CommandService) Get(id int) (*Command, *Call, error) {
	path := fmt.Sprintf("/commands/%d.json", id)
	method := "GET"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return nil, call, err
	}

	var command Command
	err = call.Do(&command)
	if err != nil {
		return nil, call, err
	}

	return &command, call, nil
}

// Upsert updates or creates a new command resource
func (service *CommandService) Upsert(command *Command) (*Command, *Call, error) {
	bytesArray, _ := json.Marshal(command)

	var path string
	var method string

	if command.ID != nil {
		path = fmt.Sprintf("/commands/%d.json", *command.ID)
		method = "PUT"
	} else {
		path = "/commands.json"
		method = "POST"
	}

	call, err := service.s.NewCall(method, path, nil, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return nil, call, err
	}

	err = call.Do(&command)
	if err != nil {
		return nil, call, err
	}

	return command, call, nil
}

// Delete deletes a sinlge command resource
func (service *CommandService) Delete(id int) (*Call, error) {
	path := fmt.Sprintf("/commands/%d.json", id)
	method := "DELETE"

	call, err := service.s.NewCall(method, path, nil, nil, nil)
	if err != nil {
		return call, err
	}

	return call, call.Do(nil)
}
