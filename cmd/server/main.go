package main

import (
	"github.com/igorrnk/ypmetrika/internal/handler"
	"github.com/igorrnk/ypmetrika/internal/storage"
	"log"
	"net/http"
	"os"
)

const (
	addressServer = "127.0.0.1:8080"
)

func main() {
	logger := log.Default()
	logger.SetOutput(os.Stdout)
	log.Println("Server is running.")

	memStorage := new(storage.MemStorage)

	server := &http.Server{
		Addr: addressServer,
		Handler: handler.Handler{
			Rep: memStorage,
		},
	}
	log.Fatal(server.ListenAndServe())
}
