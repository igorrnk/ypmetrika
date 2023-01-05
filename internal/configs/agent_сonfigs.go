package configs

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type AgentConfig struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	AddressServer  string        `env:"ADDRESS"`
}

func (config AgentConfig) String() string {
	return fmt.Sprintf("ADDRESS = %v; POLL_INTERVAL = %v; REPORT_INTERVAL = %v",
		config.AddressServer,
		config.PollInterval,
		config.ReportInterval)
}

var DefaultAgentConfig AgentConfig = AgentConfig{
	PollInterval:   2 * time.Second,
	ReportInterval: 10 * time.Second,
	AddressServer:  "127.0.0.1:8080",
}

func InitAgentConfig() AgentConfig {
	config := DefaultAgentConfig
	err := env.Parse(&config)
	if err != nil {
		log.Printf("configs.InitAgentConfig: error: %v", err)
	}
	log.Printf("Initial agent configuration: %s\n", config)
	return config
}
