package database_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"training/internal/pkg/database"
	"training/internal/pkg/service"
)

type Data struct {
	key   string
	value []string
}

var addTests = []Data{
	{key: "87257528effffff", value: []string{"محله", "ونک", "پیتزا"}},
	{key: "87257328effffff", value: []string{"محله", "ونک", "پیتزا"}},
	{key: "87267528effffff", value: []string{"محله", "ونک", "پیتزا"}},
	{key: "87157528effffff", value: []string{"محله", "ونک", "پیتزا"}},
}

func TestNewDatabase(t *testing.T) {
	arg := service.LoadConfig()
	_, err := database.NewDatabase(arg.Conf)
	assert.Equal(t, nil, err, "Connection to Database failed: %v", err)
}

func TestAdd(t *testing.T) {
	arg := service.LoadConfig()
	arg.DB.Client.FlushDB()
	for _, data := range addTests {
		err := arg.DB.Add(data.key, data.value)
		assert.Equal(t, nil, err, "Failed to add to database!")
	}
	arg.DB.Client.FlushDB()
}

func TestGet(t *testing.T) {
	arg := service.LoadConfig()
	arg.DB.Client.FlushDB()
	for _, data := range addTests {
		err := arg.DB.Add(data.key, data.value)
		if err != nil {
			fmt.Println("TEST:could not add to database.", err)
		}
	}
	for _, data := range addTests {
		value := arg.DB.Get(data.key)
		assert.Equal(t, data.value, value, "Failed to add to database!")
	}
	arg.DB.Client.FlushDB()
}
