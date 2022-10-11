package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Base struct {
	PostAPI          string `yaml:"post_api"`
	GetAPI           string `yaml:"get_api"`
	BaseUrl          string `yaml:"base_url"`
	Resolution       int    `yaml:"resolution"`
	DataBaseUrl      string `yaml:"db_url"`
	DataBasePassword string `yaml:"db_password"`
	DB               int    `yaml:"db"`
}

func (b *Base) Load(filePath string) error {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	yamlFile = []byte(os.ExpandEnv(string(yamlFile)))

	err = yaml.Unmarshal(yamlFile, b)
	if err != nil {
		return err
	}

	return nil
}
