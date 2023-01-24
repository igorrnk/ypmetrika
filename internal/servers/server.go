package servers

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/handlers"
	"github.com/igorrnk/ypmetrika/internal/middleware"
	"github.com/igorrnk/ypmetrika/internal/models"
	"github.com/igorrnk/ypmetrika/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
)

type Server struct {
	config     configs.ServerConfig
	httpServer *http.Server
	repository models.Repository
	router     chi.Router
	context    context.Context
}

func NewServer(ctx context.Context, config configs.ServerConfig) (*Server, error) {

	newServer := &Server{
		config:     config,
		context:    ctx,
		repository: storage.NewServerStorage(ctx, config),
	}
	newServer.router = chi.NewRouter()
	//newServer.Router.Use(middleware.Logger)

	//newServer.router.Use(middleware.Compress(5, "text/html", "text/json"))
	newServer.router.Use(middleware.Compress)
	h := handlers.NewHandler(config, newServer)
	newServer.router.Get("/", h.HandleFn)
	newServer.router.Get("/value/{typeMetric}/{nameMetric}", h.ValueHandleFn)
	newServer.router.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", h.UpdateHandleFn)
	newServer.router.Post("/update/", h.UpdateJSONHandleFn)
	newServer.router.Post("/value/", h.ValueJSONHandleFn)

	return newServer, nil
}

func (server *Server) Run() error {
	log.Println("Server is running.")
	server.httpServer = &http.Server{Addr: server.config.AddressServer,
		Handler: server.router}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := server.httpServer.Shutdown(server.context); err != nil {
			log.Printf("Server.Run: http.Server.Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	defer server.context.Done()

	if err := server.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	<-idleConnsClosed
	log.Println("Server has been stopped.")
	return nil
}

func (server *Server) UpdateValue(metric models.Metric) (models.Metric, error) {
	if err := server.Update(metric); err != nil {
		return models.Metric{}, err
	}
	var err error
	metric, err = server.Value(metric)
	if err != nil {
		return models.Metric{}, errors.New("Server.UpdateValue: wrong metric")
	}
	return metric, nil
}

func (server *Server) Update(metric models.Metric) error {
	if metric.Type == models.CounterType {
		if oldMetric, err := server.repository.Read(metric); err == nil {
			*metric.Value.Counter += *oldMetric.Value.Counter
		}
	}
	err := server.repository.Write(metric)
	if err != nil {
		log.Printf("Metric %v hasn't been updated.", metric.Name)
		return err
	}
	//log.Printf("Metric %v (%v) has been updated %v.", metric.Name, metric.Type, metric.Value)
	return nil
}

func (server *Server) Value(metric models.Metric) (models.Metric, error) {
	return server.repository.Read(metric)
}

func (server *Server) GetAll() ([]models.Metric, error) {
	metrics, err := server.repository.ReadAll()
	if err != nil {
		return nil, err
	}
	sort.SliceStable(metrics, func(i, j int) bool { return metrics[i].Name < metrics[j].Name })
	sort.SliceStable(metrics, func(i, j int) bool { return metrics[i].Type < metrics[j].Type })
	return metrics, nil
}
