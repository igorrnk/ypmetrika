package configs

import "time"

type AgentConfig struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	AddressServer  string
	NameCSVFile    string
}

var DefaultAgentConfig AgentConfig = AgentConfig{
	PollInterval:   2 * time.Second,
	ReportInterval: 10 * time.Second,
	AddressServer:  "127.0.0.1:8080",
}

func InitAgentConfig() AgentConfig {
	return DefaultAgentConfig
}
