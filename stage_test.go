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

func ExampleStageService_List() {
	client := New("token")

	stages, _, err := client.Stages.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, stage := range stages {
		fmt.Println(*stage.ID, *stage.Name)
	}
}

func ExampleStageService_Get() {
	client := New("token")

	stages, _, err := client.Stages.List()
	if err != nil {
		log.Fatal(err)
	}

	stage, _, err := client.Stages.Get(*stages[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*stage.Name)
}

func ExampleStageService_Upsert_create() {
	client := New("token")

	stage := &Stage{
		Name: String("test stage"),
	}

	stage, _, err := client.Stages.Upsert(stage)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleStageService_Upsert_update() {
	client := New("token")

	stage := &Stage{
		Name: String("test stage"),
	}

	stage, _, err := client.Stages.Upsert(stage)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleStageService_Delete() {
	client := New("token")

	_, err := client.Stages.Delete(3)
	if err != nil {
		log.Fatal(err)
	}
}

func TestStageServiceList(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/stages.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("stages.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	stages, call, err := client.Stages.List()
	assert.Nil(err)
	assert.Equal(len(stages), 2)
	assert.IsType(&Call{}, call)
}

func TestStageServiceList_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	stages, call, err := client.Stages.List()
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(stages)
}

func TestStageServiceList_fail_2(t *testing.T) {
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

	stages, call, err := client.Stages.List()
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(stages)
}

func TestStageServiceGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/stages/3.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("stage.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	stage, call, err := client.Stages.Get(3)
	assert.Nil(err)
	assert.Equal(*stage.ID, 3)
	assert.IsType(&Call{}, call)
}

func TestStageServiceGet_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	stage, call, err := client.Stages.Get(2)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(stage)
}

func TestStageServiceGet_fail_2(t *testing.T) {
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

	stage, call, err := client.Stages.Get(2)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(stage)
}

func TestStageServiceCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.URL.Path, "/stages.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal(payload["name"], "name")

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("stage.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	stage := &Stage{
		Name: String("name"),
	}

	stage, call, err := client.Stages.Upsert(stage)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
	assert.NotNil(stage.CreatedAt)
}

func TestStageServiceUpdate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("GET", r.Method)
		assert.Equal("/stages/3.json", r.URL.Path)
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("stage.json"))
	})

	server := httptest.NewServer(handler)
	client = New(token)
	client.BaseURL = server.URL

	stage, _, err := client.Stages.Get(3)
	assert.Nil(err)

	server.Close()

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.URL.Path, "/stages/3.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal(payload["name"], "updated name")

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("stage.json"))
	})

	server = httptest.NewServer(handler)
	client.BaseURL = server.URL
	defer server.Close()

	stage.Name = String("updated name")

	stage, call, err := client.Stages.Upsert(stage)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
	assert.NotNil(stage.CreatedAt)
}

func TestStageServiceUpsert_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	stage := &Stage{
		Name: String("name"),
	}

	stage, call, err := client.Stages.Upsert(stage)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(stage)
}

func TestStageServiceUpsert_fail_2(t *testing.T) {
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

	stage := &Stage{
		Name: String("name"),
	}

	stage, call, err := client.Stages.Upsert(stage)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(stage)
}

func TestStageServiceDelete(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.URL.Path, "/stages/3.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, "")
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	call, err := client.Stages.Delete(3)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
}

func TestStageServiceDelete_fail(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	call, err := client.Stages.Delete(3)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
}
