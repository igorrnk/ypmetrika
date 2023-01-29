package main

import (
	"context"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/servers"
	"log"
	"os"
)

func main() {

	//logFile, _ := os.OpenFile("./log/serverLog.log", os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	log.SetOutput(os.Stdout)
	//log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	config, err := configs.InitServerConfig()
	if err != nil {
		log.Fatal(err)
	}
	server, err1 := servers.NewServer(context.Background(), config)
	if err1 != nil {
		log.Fatal(err)
	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
