package configs

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type ServerConfig struct {
	AddressServer string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFileName string        `env:"STORE_FILE"`
	RestoreData   bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DBConnect     string        `env:"DATABASE_DSN"`
	DBDriverName  string
	NameHTMLFile  string
}

func InitServerConfig() (*ServerConfig, error) {
	addressServer := flag.String("a", DefaultSC.AddressServer, "The address of the server")
	storeFileName := flag.String("f", DefaultSC.StoreFileName, "The path of the data file")
	restoreData := flag.Bool("r", DefaultSC.RestoreData, "Restore from the data file")
	storeInterval := flag.Duration("i", DefaultSC.StoreInterval, "The store interval")
	key := flag.String("k", DefaultSC.Key, "The crypt key")
	dbConnect := flag.String("d", DefaultSC.DBConnect, "path to DB")
	flag.Parse()
	config := &ServerConfig{
		AddressServer: *addressServer,
		StoreInterval: *storeInterval,
		StoreFileName: *storeFileName,
		RestoreData:   *restoreData,
		Key:           *key,
		DBConnect:     *dbConnect,
		DBDriverName:  DefaultSC.DBDriverName,
		NameHTMLFile:  DefaultSC.NameHTMLFile,
	}

	err := env.Parse(config)
	if err != nil {
		log.Printf("configs.InitServerConfig: error: %v", err)
		return nil, err
	}
	log.Printf("Initial server configuration: %+v\n", config)
	return config, nil
}
