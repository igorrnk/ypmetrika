package main

import (
	"github.com/igorrnk/ypmetrika/internal/handlers"
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
	serveMux := new(http.ServeMux)
	serveMux.Handle("/update/", handlers.UpdateHandler{Rep: memStorage})
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})
	server := &http.Server{
		Addr:    addressServer,
		Handler: serveMux,
	}
	log.Fatal(server.ListenAndServe())
}
