package service_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"training/internal/pkg/service"
)

type data struct {
	search   string
	location string
	expected answer
}

type answer struct {
	H3Key string   `json:"h3index"`
	Word  []string `json:"words"`
}

var postTests = []data{
	{search: "محله ونک پیتزا", location: "53.4,10.6", expected: answer{H3Key: "871f02503ffffff", Word: []string{"محله", "ونک", "پیتزا"}}},
	{search: "محله ونک پیتزا", location: "53.4,10.6", expected: answer{H3Key: "871f02503ffffff", Word: []string{"محله", "ونک", "پیتزا", "محله", "ونک", "پیتزا"}}},
	{search: "مرزداران ناهید", location: "53.2,22.3", expected: answer{H3Key: "871f51c50ffffff", Word: []string{"مرزداران", "ناهید"}}},
}

func TestLoadConfig(t *testing.T) {
	arg := service.LoadConfig()
	assert.NotEmpty(t, arg, "Failed to load configs from yaml file.")
}

func TestEngine(t *testing.T) {
	arg := service.LoadConfig()
	arg.DB.Client.FlushDB()
	engine := service.Engine(arg)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "", nil)
	assert.Equal(t, nil, err, "building request: %v", err)
	engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNotFound, recorder.Code, "bad status: %v", recorder.Body.String())
}

func TestPostAlbums(t *testing.T) {
	arg := service.LoadConfig()
	arg.DB.Client.FlushDB()
	engine := service.Engine(arg)
	for _, row := range postTests {
		body, _ := json.Marshal(map[string]string{"search": row.search, "location": row.location})
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("POST", "/albums", bytes.NewBuffer(body))
		assert.Equal(t, nil, err, "building request: %v", err)
		engine.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusCreated, recorder.Code, "bad status code: %v", recorder.Body.String())
		expected, _ := json.MarshalIndent(row.expected, "", "    ")
		assert.Equal(t, string(expected), recorder.Body.String(), "bad response: %v", recorder.Body.String())
	}
}
func TestGetAlbums(t *testing.T) {
	arg := service.LoadConfig()
	arg.DB.Client.FlushDB()
	engine := service.Engine(arg)
	for _, row := range postTests {
		body, _ := json.Marshal(map[string]string{"search": row.search, "location": row.location})
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("POST", "/albums", bytes.NewBuffer(body))
		engine.ServeHTTP(recorder, request)

		recorder = httptest.NewRecorder()
		url := "/albums/" + row.expected.H3Key
		request, err = http.NewRequest("GET", url, nil)
		assert.Equal(t, nil, err, "building request: %v", err)
		engine.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusOK, recorder.Code, "bad status code: %v", recorder.Body.String())
		expected, _ := json.MarshalIndent(row.expected, "", "    ")
		assert.Equal(t, string(expected), recorder.Body.String(), "bad response: %v", recorder.Body.String())
	}
}
