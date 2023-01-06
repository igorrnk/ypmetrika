package configs

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type ServerConfig struct {
	AddressServer string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFileName string        `env:"STORE_FILE"`
	RestoreData   bool          `env:"RESTORE"`
	NameHTMLFile  string
}

func (config ServerConfig) String() string {
	return fmt.Sprintf("ADDRESS = %v", config.AddressServer)
}

var DefaultServerConfig = ServerConfig{
	AddressServer: "127.0.0.1:8080",
	StoreInterval: 30 * time.Second,
	StoreFileName: "/tmp/devops-metrics-db.json",
	RestoreData:   true,
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
