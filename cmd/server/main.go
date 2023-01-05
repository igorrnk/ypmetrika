package main

import (
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/servers"
	"log"
	"os"
)

func main() {

	logFile, _ := os.OpenFile("./log/serverLog.log", os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0777)
	log.SetOutput(logFile)
	//log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	config := configs.InitServerConfig()
	server, err := servers.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
