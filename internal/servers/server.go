package servers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/igorrnk/ypmetrika/configs"
	"github.com/igorrnk/ypmetrika/internal/handlers"
	"github.com/igorrnk/ypmetrika/internal/models"
	"github.com/igorrnk/ypmetrika/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	Config     configs.ServerConfig
	Repository storage.ServerRepository
	Router     chi.Router
}

func NewServer(config configs.ServerConfig) (*Server, error) {
	newServer := &Server{
		Config:     configs.InitServerConfig(),
		Repository: storage.NewServerMemoryStorage(),
	}
	newServer.Router = chi.NewRouter()
	h := handlers.NewHandler(config, newServer)
	newServer.Router.Get("/", h.HandleFn)
	newServer.Router.Get("/value/{typeMetric}/{nameMetric}", h.ValueHandleFn)
	newServer.Router.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", h.UpdateHandleFn)

	return newServer, nil
}

func (server *Server) Run() error {
	log.Println("Server is running.")
	s := &http.Server{Addr: server.Config.AddressServer,
		Handler: server.Router}
	go func() {
		err := s.ListenAndServe()
		log.Fatal(err)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	ctx, shutdown := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdown()

	return s.Shutdown(ctx)

}

func (server *Server) Update(metric models.ServerMetric) error {
	err := server.Repository.Write(metric)
	if err != nil {
		log.Printf("Metric %v hasn't been updated.", metric.Name)
		return err
	}
	log.Printf("Metric %v has been updated.", metric.Name)
	return nil
}

func (server *Server) Value(metric models.ServerMetric) (models.ServerMetric, bool) {
	return server.Repository.Read(metric)
}

func (server *Server) All() []models.ServerMetric {
	return server.Repository.ReadAll()
}
