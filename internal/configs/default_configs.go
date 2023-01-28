package configs

import "time"

var DefaultAC AgentConfig = AgentConfig{
	PollInterval:   2 * time.Second,
	ReportInterval: 10 * time.Second,
	AddressServer:  "http://127.0.0.1:8080",
}

var DefaultSC = ServerConfig{
	AddressServer: "127.0.0.1:8080",
	StoreInterval: 30 * time.Second,
	StoreFileName: "/tmp/devops-metrics-db.json",
	RestoreData:   true,
	NameHTMLFile:  "./web/metrics.html",
}
