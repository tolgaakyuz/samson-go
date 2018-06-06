package samson

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleCommandService_List() {
	client := New("token")

	commands, _, err := client.Commands.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, command := range commands {
		fmt.Println(*command.Command)
	}
}

func ExampleCommandService_Get() {
	client := New("token")

	commands, _, err := client.Commands.List()
	if err != nil {
		log.Fatal(err)
	}

	command, _, err := client.Commands.Get(*commands[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*command.Command)
}

func ExampleCommandService_Upsert_create() {
	client := New("token")

	command := &Command{
		Command:   String("go lint"),
		ProjectID: String("1"),
	}

	command, _, err := client.Commands.Upsert(command)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCommandService_Upsert_update() {
	ExampleCommandService_Upsert_create()

	client := New("token")

	command := &Command{
		Command:   String("go lint"),
		ProjectID: String("1"),
	}

	command, _, err := client.Commands.Upsert(command)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCommandService_Delete() {
	client := New("token")

	_, err := client.Commands.Delete(1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCommandServiceList(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/commands.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("commands.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	commands, call, err := client.Commands.List()
	assert.Nil(err)
	assert.Equal(len(commands), 3)
	assert.IsType(&Call{}, call)
}

func TestCommandServiceList_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	commands, call, err := client.Commands.List()
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(commands)
}

func TestCommandServiceList_fail_2(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, "malformed json response")
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	commands, call, err := client.Commands.List()
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(commands)
}

func TestCommandServiceGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/commands/1.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("command.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	command, call, err := client.Commands.Get(1)
	assert.Nil(err)
	assert.Equal(*command.ID, 1)
	assert.IsType(&Call{}, call)
}

func TestCommandServiceGet_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	commands, call, err := client.Commands.Get(5)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(commands)
}

func TestCommandServiceGet_fail_2(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, "malformed json response")
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	command, call, err := client.Commands.Get(2)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(command)
}

func TestCommandServiceCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.URL.Path, "/commands.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal(payload["command"], "test command")
		assert.Equal(payload["project_id"], "1")

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("command.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	command := &Command{
		Command:   String("test command"),
		ProjectID: String("1"),
	}

	command, call, err := client.Commands.Upsert(command)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
}

func TestCommandServiceCreate_global(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.URL.Path, "/commands.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal(payload["command"], "test command")
		assert.Equal(payload["project_id"], nil)

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("command_global.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	command := &Command{
		Command: String("test command"),
	}

	command, call, err := client.Commands.Upsert(command)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
	assert.True(command.IsGlobal())
}

func TestCommandServiceUpdate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("GET", r.Method)
		assert.Equal("/commands/1.json", r.URL.Path)
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("command.json"))
	})

	server := httptest.NewServer(handler)
	client = New(token)
	client.BaseURL = server.URL

	command, _, err := client.Commands.Get(1)
	assert.Nil(err)

	server.Close()

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.URL.Path, "/commands/1.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal(payload["command"], "updated command")
		assert.Equal(payload["project_id"], "1")

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("command.json"))
	})

	server = httptest.NewServer(handler)
	client.BaseURL = server.URL
	defer server.Close()

	command.Command = String("updated command")

	command, call, err := client.Commands.Upsert(command)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
}

func TestCommandServiceUpsert_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	command := &Command{
		Command:   String("test command"),
		ProjectID: String("1"),
	}

	command, call, err := client.Commands.Upsert(command)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(command)
}

func TestCommandServiceUpsert_fail_2(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, "malformed json response")
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	command := &Command{
		Command:   String("test command"),
		ProjectID: String("1"),
	}

	command, call, err := client.Commands.Upsert(command)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(command)
}

func TestCommandServiceDelete(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.URL.Path, "/commands/2.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, "")
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	call, err := client.Commands.Delete(2)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
}

func TestCommandServiceDelete_fail(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	call, err := client.Commands.Delete(2)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
}
