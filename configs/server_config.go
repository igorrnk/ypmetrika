package configs

type ServerConfig struct {
	AddressServer string
	NameHTMLFile  string
}

var DefaultServerConfig ServerConfig = ServerConfig{
	AddressServer: "127.0.0.1:8080",
	NameHTMLFile:  "./web/metrics.html",
}

func InitServerConfig() ServerConfig {
	return DefaultServerConfig
}
