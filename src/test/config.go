package test

import "sync"

type TestConfig struct {
	ShopServiceUrl string
}

var config *TestConfig

func NewConfig() *TestConfig {
	return &TestConfig{
		ShopServiceUrl: "localhost:8080",
	}
}

var get_config_once sync.Once

func GetConfig() TestConfig {
	get_config_once.Do(func() {
		config = NewConfig()
	})

	return *config
}
