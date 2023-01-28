package configs

import (
	"flag"
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
	return fmt.Sprintf("Address = %v; Poll interval = %v; Report interval = %v",
		config.AddressServer,
		config.PollInterval,
		config.ReportInterval)
}

var DefaultAgentConfig AgentConfig = AgentConfig{
	PollInterval:   2 * time.Second,
	ReportInterval: 10 * time.Second,
	AddressServer:  "http://127.0.0.1:8080",
}

func InitAgentConfig() AgentConfig {
	config := DefaultAgentConfig
	addressServer := flag.String("a", "http://127.0.0.1:8080", "The address of the server")
	var err error
	flag.Func("p", "The poll interval", func(flagValue string) error {
		if config.PollInterval, err = time.ParseDuration(flagValue); err != nil {
			return err
		}
		return nil
	})
	flag.Func("r", "The report interval", func(flagValue string) error {
		if config.ReportInterval, err = time.ParseDuration(flagValue); err != nil {
			return err
		}
		return nil
	})
	flag.Parse()
	config.AddressServer = *addressServer
	err = env.Parse(&config)
	if err != nil {
		log.Printf("configs.InitAgentConfig: error: %v", err)
	}
	log.Printf("Initial agent configuration: %s\n", config)
	return config
}
