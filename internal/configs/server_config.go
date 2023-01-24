package configs

import (
	"flag"
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
	addressServer := flag.String("a", "127.0.0.1:8080", "The address of the server")
	storeFileName := flag.String("f", "/tmp/devops-metrics-db.json", "The path of the data file")
	restoreData := flag.Bool("r", true, "Restore from the data file")

	var err error
	flag.Func("i", "The store interval", func(flagValue string) error {
		if config.StoreInterval, err = time.ParseDuration(flagValue); err != nil {
			return err
		}
		return nil
	})

	flag.Parse()
	config.AddressServer = *addressServer
	config.StoreFileName = *storeFileName
	config.RestoreData = *restoreData

	err = env.Parse(&config)
	if err != nil {
		log.Printf("configs.InitServerConfig: error: %v", err)
	}
	log.Printf("Initial server configuration: %s\n", config)
	return config
}
