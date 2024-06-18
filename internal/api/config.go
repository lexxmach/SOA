package api

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port uint16 `json:"port"`

	Title   string `json:"title"`
	Version string `json:"version"`

	DB struct {
		// If set to true, db will bocked
		Mock bool `json:"mock"`

		Host string `json:"host"`
		Port uint16 `json:"port"`
	} `json:"db"`

	JWTSecret string `json:"jwt"`
}

func MustGetConfig(path string) *Config {
	plan, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(plan, &config)
	if err != nil {
		panic(err)
	}

	return &config
}
