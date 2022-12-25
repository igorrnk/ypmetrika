package main

import (
	"github.com/igorrnk/ypmetrika/internal/servers"
	"log"
	"os"
)

func main() {
	logger := log.Default()
	logger.SetOutput(os.Stdout)

	server := servers.NewServer()
	log.Fatal(server.ListenAndServe())
}
