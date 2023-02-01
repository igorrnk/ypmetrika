package services

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

type Service struct {
	config     *configs.ServerConfig
	httpServer *http.Server
	repository models.Repository
	router     chi.Router
	context    context.Context
}

func NewService(ctx context.Context, config *configs.ServerConfig) (*Service, error) {

	newServer := &Service{
		config:     config,
		context:    ctx,
		repository: storage.NewFileStorage(ctx, config),
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

func (server *Service) Run() error {
	log.Println("Service is running.")
	server.httpServer = &http.Server{Addr: server.config.AddressServer,
		Handler: server.router}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := server.httpServer.Shutdown(server.context); err != nil {
			log.Printf("Service.Run: http.Service.Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	defer server.context.Done()
	if err := server.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	<-idleConnsClosed
	log.Println("Service has been stopped.")
	return nil
}

func (server *Service) UpdateValue(metric *models.Metric) (*models.Metric, error) {
	var err error
	if err = server.Update(metric); err != nil {
		return nil, err
	}
	metric, err = server.Value(metric)
	if err != nil {
		return nil, errors.New("Service.UpdateValue: wrong metric")
	}
	return metric, nil
}

func (server *Service) Update(metric *models.Metric) error {
	if metric.Type == models.CounterType {
		oldMetric, err := server.repository.Read(metric)
		if err != nil && !errors.Is(err, models.ErrNotFound) {
			return err
		}
		if errors.Is(err, models.ErrNotFound) {
			// do nothing
		} else {
			metric.Counter += oldMetric.Counter
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

func (server *Service) Value(metric *models.Metric) (*models.Metric, error) {
	return server.repository.Read(metric)
}

func (server *Service) GetAll() ([]models.Metric, error) {
	metrics, err := server.repository.ReadAll()
	if err != nil {
		return nil, err
	}
	sort.SliceStable(metrics, func(i, j int) bool { return metrics[i].Name < metrics[j].Name })
	sort.SliceStable(metrics, func(i, j int) bool { return metrics[i].Type < metrics[j].Type })
	return metrics, nil
}
