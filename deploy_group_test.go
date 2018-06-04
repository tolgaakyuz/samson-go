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

func ExampleDeployGroupService_List() {
	client := New("token")

	deployGroups, _, err := client.DeployGroups.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, dg := range deployGroups {
		fmt.Println(*dg.ID, *dg.Name)
	}
}

func ExampleDeployGroupService_Get() {
	client := New("token")

	deployGroups, _, err := client.DeployGroups.List()
	if err != nil {
		log.Fatal(err)
	}

	dg, _, err := client.DeployGroups.Get(*deployGroups[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*dg.Name)
}

func ExampleDeployGroupService_Upsert_create() {
	client := New("token")

	dg := &DeployGroup{
		Name: String("test deploy group"),
	}

	dg, _, err := client.DeployGroups.Upsert(dg)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleDeployGroupService_Upsert_update() {
	client := New("token")

	dg := &DeployGroup{
		Name: String("test deploy groups"),
	}

	dg, _, err := client.DeployGroups.Upsert(dg)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleDeployGroupService_Delete() {
	client := New("token")

	_, err := client.DeployGroups.Delete(1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestDeployGroupServiceList(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/deploy_groups.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("deploy_groups.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	deployGroups, call, err := client.DeployGroups.List()
	assert.Nil(err)
	assert.Equal(len(deployGroups), 2)
	assert.IsType(&Call{}, call)
}

func TestDeployGroupServiceList_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	deployGroups, call, err := client.DeployGroups.List()
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(deployGroups)
}

func TestDeployGroupServiceList_fail_2(t *testing.T) {
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

	deployGroups, call, err := client.DeployGroups.List()
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(deployGroups)
}

func TestDeployGroupServiceGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, "/deploy_groups/3.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("deploy_group.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	dg, call, err := client.DeployGroups.Get(3)
	assert.Nil(err)
	assert.Equal(*dg.ID, 3)
	assert.IsType(&Call{}, call)
}

func TestDeployGroupServiceGet_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	dg, call, err := client.DeployGroups.Get(2)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(dg)
}

func TestDeployGroupServiceGet_fail_2(t *testing.T) {
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

	dg, call, err := client.DeployGroups.Get(2)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(dg)
}

func TestDeployGroupServiceCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "POST")
		assert.Equal(r.URL.Path, "/deploy_groups.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal(payload["name"], "name")

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("deploy_group.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	dg := &DeployGroup{
		Name: String("name"),
	}

	dg, call, err := client.DeployGroups.Upsert(dg)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
	assert.NotNil(dg.CreatedAt)
}

func TestDeployGroupServiceUpdate(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("GET", r.Method)
		assert.Equal("/deploy_groups/3.json", r.URL.Path)
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("deploy_group.json"))
	})

	server := httptest.NewServer(handler)
	client = New(token)
	client.BaseURL = server.URL

	dg, _, err := client.DeployGroups.Get(3)
	assert.Nil(err)

	server.Close()

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "PUT")
		assert.Equal(r.URL.Path, "/deploy_groups/3.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.Nil(err)
		assert.Equal(payload["name"], "updated name")

		w.WriteHeader(200)
		fmt.Fprintln(w, readTestData("deploy_group.json"))
	})

	server = httptest.NewServer(handler)
	client.BaseURL = server.URL
	defer server.Close()

	dg.Name = String("updated name")

	dg, call, err := client.DeployGroups.Upsert(dg)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
	assert.NotNil(dg.CreatedAt)
}

func TestDeployGroupServiceUpsert_fail_1(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	dg := &DeployGroup{
		Name: String("name"),
	}

	dg, call, err := client.DeployGroups.Upsert(dg)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(dg)
}

func TestDeployGroupServiceUpsert_fail_2(t *testing.T) {
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

	dg := &DeployGroup{
		Name: String("name"),
	}

	dg, call, err := client.DeployGroups.Upsert(dg)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
	assert.Nil(dg)
}

func TestDeployGroupServiceDelete(t *testing.T) {
	var err error
	assert := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "DELETE")
		assert.Equal(r.URL.Path, "/deploy_groups/3.json")
		assert.Equal("Bearer "+token, r.Header.Get("Authorization"))
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(200)
		fmt.Fprintln(w, "")
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client = New(token)
	client.BaseURL = server.URL

	call, err := client.DeployGroups.Delete(3)
	assert.Nil(err)
	assert.IsType(&Call{}, call)
}

func TestDeployGroupServiceDelete_fail(t *testing.T) {
	var err error
	assert := assert.New(t)

	client = New(token)
	client.BaseURL = "^http://localhost"

	call, err := client.DeployGroups.Delete(3)
	assert.NotNil(err)
	assert.IsType(&Call{}, call)
}
