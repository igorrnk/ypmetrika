package configs

import "time"

var DefaultAC AgentConfig = AgentConfig{
	AddressServer:  "127.0.0.1:8080",
	PollInterval:   2 * time.Second,
	ReportInterval: 10 * time.Second,
	Key:            "",
}

var DefaultSC = ServerConfig{
	AddressServer: "127.0.0.1:8080",
	StoreInterval: 30 * time.Second,
	StoreFileName: "/tmp/devops-metrics-db.json",
	RestoreData:   true,
	Key:           "",
	DBConnect:     "postgres://username:password@localhost:5432/database_name",
	NameHTMLFile:  "./web/metrics.html",
}
