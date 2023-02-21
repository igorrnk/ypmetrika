package configs

import "time"

var DefaultAC AgentConfig = AgentConfig{
	AddressServer:  "127.0.0.1:8080",
	PollInterval:   2 * time.Second,
	ReportInterval: 10 * time.Second,
	Key:            "",
	Limit:          0,
}

var DefaultSC = ServerConfig{
	AddressServer: "127.0.0.1:8080",
	StoreInterval: 30 * time.Second,
	StoreFileName: "/tmp/devops-metrics-db.json",
	RestoreData:   true,
	Key:           "",
	DBConnect:     "",
	DBDriverName:  "pgx",
	NameHTMLFile:  "./web/metrics.html",
}
