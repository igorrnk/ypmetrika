package configs

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
)

type ServerConfig struct {
	AddressServer string `env:"ADDRESS"`
	NameHTMLFile  string
}

func (config ServerConfig) String() string {
	return fmt.Sprintf("ADDRESS = %v", config.AddressServer)
}

var DefaultServerConfig = ServerConfig{
	AddressServer: "127.0.0.1:8080",
	NameHTMLFile:  "./web/metrics.html",
}

func InitServerConfig() ServerConfig {
	config := DefaultServerConfig
	err := env.Parse(&config)
	if err != nil {
		log.Printf("configs.InitServerConfig: error: %v", err)
	}
	log.Printf("Initial server configuration: %s\n", config)
	return config
}
