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

		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		Port     uint16 `json:"port"`
	} `json:"db"`

	PostAddress string `json:"post_address"`

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
