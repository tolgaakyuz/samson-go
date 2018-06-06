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

func ExampleEnvironmentService_List() {
	client := New("token")

	environments, _, err := client.Environments.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, environment := range environments {
		fmt.Println(*environment.Name)
	}
}

func ExampleEnvironmentService_Get() {
	client := New("token")

	environments, _, err := client.Environments.List()
	if err != nil {
		log.Fatal(err)
	}

	environment, _, err := client.Environments.Get(*environments[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*environment.Name)
}

func ExampleEnvironmentService_Upsert_create() {
	client := New("token")

	environment := &Environment{
		Name:       String("staging"),
		Production: String("1"),
	}

	environment, _, err := client.Environments.Upsert(environment)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleEnvironmentService_Upsert_update() {
	ExampleEnvironmentService_Upsert_create()

	client := New("token")

	environment := &Environment{
		Name:       String("staging"),
		Production: String("1"),
	}

	environment, _, err := client.Environments.Upsert(environment)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleEnvironmentService_Delete() {
	client := New("token")

	_, err := client.Environments.Delete(1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestEnvironmentServiceList(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/environments.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("environments.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	environments, call, err := client.Environments.List()
	assert.Nil(err)
	assert.Equal(len(environments), 2)
	assert.IsType(&Call{}, call)
}

func TestEnvironmentServiceList_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	environments, call, err := client.Environments.List()
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(environments)
}

func TestEnvironmentServiceList_fail_2(t *testing.T) {
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

	environments, call, err := client.Environments.List()
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(environments)
}

func TestEnvironmentServiceGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/environments/1.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("environment_prod.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	environment, call, err := client.Environments.Get(1)
	assert.Nil(err)
	assert.Equal(*environment.ID, 1)
	assert.Equal(*environment.Name, "production")
	assert.Equal(*environment.Production, "1")
	assert.IsType(&Call{}, call)
}

func TestEnvironmentServiceGet_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	environments, call, err := client.Environments.Get(5)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(environments)
}

func TestEnvironmentServiceGet_fail_2(t *testing.T) {
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

	environment, call, err := client.Environments.Get(2)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(environment)
}

func TestEnvironmentServiceCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.URL.Path, "/environments.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal(payload["name"], "staging")
		assert.Equal(payload["production"], "1")

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("environment_staging.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	environment := &Environment{
		Name:       String("staging"),
		Production: String("1"),
	}

	environment, call, err := client.Environments.Upsert(environment)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
}

func TestEnvironmentServiceUpdate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("GET", r.Method)
		assert.Equal("/environments/1.json", r.URL.Path)
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("environment_prod.json"))
	})

	server := httptest.NewServer(handler)
	client = New(token)
	client.BaseURL = server.URL

	environment, _, err := client.Environments.Get(1)
	assert.Nil(err)

	server.Close()

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.URL.Path, "/environments/1.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal(payload["name"], "preview")
		assert.Equal(payload["production"], "0")

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("environment_staging.json"))
	})

	server = httptest.NewServer(handler)
	client.BaseURL = server.URL
	defer server.Close()

	environment.Name = String("preview")
	environment.Production = String("0")

	environment, call, err := client.Environments.Upsert(environment)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
}

func TestEnvironmentServiceUpsert_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	environment := &Environment{
		Name:       String("staging"),
		Production: String("1"),
	}

	environment, call, err := client.Environments.Upsert(environment)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(environment)
}

func TestEnvironmentServiceUpsert_fail_2(t *testing.T) {
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

	environment := &Environment{
		Name:       String("staging"),
		Production: String("1"),
	}

	environment, call, err := client.Environments.Upsert(environment)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(environment)
}

func TestEnvironmentServiceDelete(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.URL.Path, "/environments/2.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, "")
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	call, err := client.Environments.Delete(2)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
}

func TestEnvironmentServiceDelete_fail(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	call, err := client.Environments.Delete(2)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
}
