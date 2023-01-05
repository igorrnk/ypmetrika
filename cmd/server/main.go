package main

import (
	"github.com/igorrnk/ypmetrika/configs"
	"github.com/igorrnk/ypmetrika/internal/servers"
	"log"
	"os"
)

func main() {
	logger := log.Default()
	logFile, _ := os.OpenFile("./log/serverLog.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)
	logger.SetOutput(logFile)
	logger.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	config := configs.InitServerConfig()
	server, err := servers.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
