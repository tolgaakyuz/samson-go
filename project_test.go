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

func ExampleProjectService_List() {
	client := New("token")

	projects, _, err := client.Projects.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, project := range projects {
		fmt.Println(*project.ID, *project.Name)
	}
}

func ExampleProjectService_Get() {
	client := New("token")

	projects, _, err := client.Projects.List()
	if err != nil {
		log.Fatal(err)
	}

	project, _, err := client.Projects.Get(*projects[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*project.Name)
}

func ExampleProjectService_Upsert_create() {
	client := New("token")

	project := &Project{
		Name:        String("test project"),
		Description: String("test project description"),
	}

	project, _, err := client.Projects.Upsert(project)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleProjectService_Upsert_update() {
	ExampleProjectService_Upsert_create()

	client := New("token")

	project := &Project{
		Name:        String("test project"),
		Description: String("test project description - updated"),
	}

	project, _, err := client.Projects.Upsert(project)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleProjectService_Delete() {
	client := New("token")

	_, err := client.Projects.Delete(2)
	if err != nil {
		log.Fatal(err)
	}
}

func TestProjectServiceList(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/projects.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("projects.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	projects, call, err := client.Projects.List()
	assert.Nil(err)
	assert.Equal(len(projects), 2)
	assert.IsType(&Call{}, call)
}

func TestProjectServiceList_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	projects, call, err := client.Projects.List()
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(projects)
}

func TestProjectServiceList_fail_2(t *testing.T) {
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

	projects, call, err := client.Projects.List()
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(projects)
}

func TestProjectServiceGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/projects/2.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("project.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	project, call, err := client.Projects.Get(2)
	assert.Nil(err)
	assert.Equal(*project.ID, 2)
	assert.IsType(&Call{}, call)
}

func TestProjectServiceGet_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	projects, call, err := client.Projects.Get(2)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(projects)
}

func TestProjectServiceGet_fail_2(t *testing.T) {
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

	project, call, err := client.Projects.Get(2)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(project)
}

func TestProjectServiceCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.URL.Path, "/projects.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal(payload["name"], "name")
		assert.Equal(payload["description"], "description")

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("project.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	project := &Project{
		Name:        String("name"),
		Description: String("description"),
	}

	project, call, err := client.Projects.Upsert(project)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
	assert.NotNil(project.CreatedAt)
}

func TestProjectServiceUpdate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("GET", r.Method)
		assert.Equal("/projects/2.json", r.URL.Path)
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("project.json"))
	})

	server := httptest.NewServer(handler)
	client = New(token)
	client.BaseURL = server.URL

	project, _, err := client.Projects.Get(2)
	assert.Nil(err)

	server.Close()

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.URL.Path, "/projects/2.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal(payload["name"], "Example-kubernetes")
		assert.Equal(payload["description"], "updated description")

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("project.json"))
	})

	server = httptest.NewServer(handler)
	client.BaseURL = server.URL
	defer server.Close()

	project.Description = String("updated description")

	project, call, err := client.Projects.Upsert(project)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
	assert.NotNil(project.CreatedAt)
}

func TestProjectServiceUpsert_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	project := &Project{
		Name:        String("name"),
		Description: String("description"),
	}

	project, call, err := client.Projects.Upsert(project)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(project)
}

func TestProjectServiceUpsert_fail_2(t *testing.T) {
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

	project := &Project{
		Name:        String("name"),
		Description: String("description"),
	}

	project, call, err := client.Projects.Upsert(project)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(project)
}

func TestProjectServiceDelete(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.URL.Path, "/projects/2.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, "")
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	call, err := client.Projects.Delete(2)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
}

func TestProjectServiceDelete_fail(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	call, err := client.Projects.Delete(2)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
}
