package servers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/igorrnk/ypmetrika/configs"
	"github.com/igorrnk/ypmetrika/internal/handlers"
	"github.com/igorrnk/ypmetrika/internal/models"
	"github.com/igorrnk/ypmetrika/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"time"
)

type Server struct {
	Config     configs.ServerConfig
	Repository models.Repository
	Router     chi.Router
}

func NewServer(config configs.ServerConfig) (*Server, error) {
	newServer := &Server{
		Config:     configs.InitServerConfig(),
		Repository: storage.New(),
	}
	newServer.Router = chi.NewRouter()
	newServer.Router.Use(middleware.Logger)
	h := handlers.NewHandler(config, newServer)
	newServer.Router.Get("/", h.HandleFn)
	newServer.Router.Get("/value/{typeMetric}/{nameMetric}", h.ValueHandleFn)
	newServer.Router.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", h.UpdateHandleFn)
	newServer.Router.Post("/update/", h.UpdateJSONHandleFn)
	newServer.Router.Post("/value/", h.ValueJSONHandleFn)

	return newServer, nil
}

func (server *Server) Run() error {
	log.Println("Server is running.")
	s := &http.Server{Addr: server.Config.AddressServer,
		Handler: server.Router}
	go func() {
		err := s.ListenAndServe()
		log.Println(err)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	ctx, shutdown := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdown()

	return s.Shutdown(ctx)

}

func (server *Server) Update(metric models.Metric) error {
	if metric.Type == models.CounterType {
		if oldMetric, ok := server.Repository.Read(metric); ok {
			metric.Value.Counter += oldMetric.Value.Counter
		}
	}
	err := server.Repository.Write(metric)
	if err != nil {
		log.Printf("Metric %v hasn't been updated.", metric.Name)
		return err
	}
	//log.Printf("Metric %v (%v) has been updated %v.", metric.Name, metric.Type, metric.Value)
	return nil
}

func (server *Server) Value(metric models.Metric) (models.Metric, bool) {
	return server.Repository.Read(metric)
}

func (server *Server) GetAll() []models.Metric {
	metrics, _ := server.Repository.ReadAll()
	sort.SliceStable(metrics, func(i, j int) bool { return metrics[i].Name < metrics[j].Name })
	sort.SliceStable(metrics, func(i, j int) bool { return metrics[i].Type < metrics[j].Type })
	return metrics
}
