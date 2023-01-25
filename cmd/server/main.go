package main

import (
	"github.com/igorrnk/ypmetrika/configs"
	"github.com/igorrnk/ypmetrika/internal/servers"
	"log"
	"os"
)

func main() {
	logger := log.Default()
	logger.SetOutput(os.Stdout)

	config := configs.InitServerConfig()

	server, err := servers.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
