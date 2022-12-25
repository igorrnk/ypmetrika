package servers

import (
	"github.com/igorrnk/ypmetrika/internal/handlers"
	"github.com/igorrnk/ypmetrika/internal/storage"
	"log"
	"net/http"
)

type Config struct {
	addressServer string
	nameCSVFile   string
}

var DefaultConfig Config = Config{
	addressServer: "127.0.0.1:8080",
	nameCSVFile:   "./configs/metrics.csv",
}

type Server struct {
	http.Server
	Config     Config
	Repository storage.Repositories
}

func NewServer() *Server {
	newServer := Server{
		Config:     DefaultConfig,
		Repository: &storage.MemStorage{},
	}
	serveMux := new(http.ServeMux)
	serveMux.Handle("/update/", handlers.UpdateHandler{Rep: newServer.Repository})
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})
	newServer.Addr = newServer.Config.addressServer
	newServer.Handler = serveMux

	return &newServer
}

func (server *Server) ListenAndServe() error {
	log.Println("Server is running.")
	return server.Server.ListenAndServe()
}
