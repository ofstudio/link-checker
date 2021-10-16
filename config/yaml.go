package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

type Profile struct {
	Id       string
	Title    string        `yaml:"title"`
	HomePage string        `yaml:"home_page"`
	Delay    time.Duration `yaml:"delay"`
}

type config struct {
	Profiles map[string]Profile `yaml:"profiles"`
}

func readYaml(filepath string) (*config, error) {
	var c config
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
