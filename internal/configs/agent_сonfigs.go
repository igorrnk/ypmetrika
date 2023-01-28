package configs

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type AgentConfig struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	AddressServer  string        `env:"ADDRESS"`
}

func InitAgentConfig() (*AgentConfig, error) {
	pollInterval := flag.Duration("p", DefaultAC.PollInterval, "The poll interval")
	reportInterval := flag.Duration("r", DefaultAC.ReportInterval, "The report interval")
	addressServer := flag.String("a", DefaultAC.AddressServer, "The address of the server")
	flag.Parse()
	agentConfig := &AgentConfig{
		PollInterval:   *pollInterval,
		ReportInterval: *reportInterval,
		AddressServer:  *addressServer,
	}
	err := env.Parse(agentConfig)
	if err != nil {
		log.Printf("configs.InitAgentConfig: error: %v", err)
		return nil, err
	}
	log.Printf("Initial agent configuration: %+v\n", agentConfig)
	return agentConfig, nil
}
