package services

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/crypts"
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
	crypter    models.Crypter
}

func NewService(ctx context.Context, config *configs.ServerConfig) (*Service, error) {

	newServer := &Service{
		config:     config,
		context:    ctx,
		repository: storage.NewFileStorage(ctx, config),
		crypter:    crypts.NewCrypterSHA256(config.Key),
	}
	var err error

	if config.DBConnect == "" {
		newServer.repository = storage.NewFileStorage(ctx, config)
	} else {
		newServer.repository, err = storage.NewDBStorage(ctx, config)
	}
	if err != nil {
		return nil, err
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
	newServer.router.Get("/ping", h.PingHandleFn)
	newServer.router.Post("/updates/", h.UpdatesJSONFn)

	return newServer, nil
}

func (service *Service) Run() error {
	log.Println("Service is running.")
	service.httpServer = &http.Server{Addr: service.config.AddressServer,
		Handler: service.router}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := service.httpServer.Shutdown(service.context); err != nil {
			log.Printf("Service.Run: http.Service.Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	defer service.context.Done()
	if err := service.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	<-idleConnsClosed
	service.Close()
	log.Println("Service has been stopped.")
	return nil
}

func (service *Service) UpdateValue(metric *models.Metric) (*models.Metric, error) {
	var err error
	if err = service.crypter.CheckHash(metric); err != nil {
		return nil, err
	}
	if err = service.Update(metric); err != nil {
		return nil, err
	}
	metric, err = service.Value(metric)
	if err != nil {
		return nil, errors.New("Service.UpdateValue: wrong metric")
	}
	return metric, nil
}

func (service *Service) Update(metric *models.Metric) error {
	oldMetric := &models.Metric{
		Name: metric.Name,
		Type: metric.Type,
	}
	var err error
	if metric.Type == models.CounterType {
		oldMetric, err = service.repository.Read(oldMetric)

		if err != nil && !errors.Is(err, models.ErrNotFound) {
			return err
		}
		if errors.Is(err, models.ErrNotFound) {
			// do nothing
		} else {
			metric.Counter += oldMetric.Counter
		}
	}
	err = service.repository.Write(metric)
	if err != nil {
		log.Printf("Metric %v hasn't been updated.", metric.Name)
		return err
	}
	//log.Printf("Metric %v (%v) has been updated %v.", metric.Name, metric.Type, metric.Value)
	return nil
}

func (service *Service) Updates(metrics []*models.Metric) error {
	for _, metric := range metrics {
		err := service.Update(metric)
		if err != nil {

		}
	}
	return nil
}

func (service *Service) Value(metric *models.Metric) (*models.Metric, error) {
	var err error
	metric, err = service.repository.Read(metric)
	if err != nil {
		return nil, err
	}
	service.crypter.AddHash(metric)
	return metric, err
}

func (service *Service) GetAll() ([]models.Metric, error) {
	metrics, err := service.repository.ReadAll()
	if err != nil {
		return nil, err
	}
	sort.SliceStable(metrics, func(i, j int) bool { return metrics[i].Name < metrics[j].Name })
	sort.SliceStable(metrics, func(i, j int) bool { return metrics[i].Type < metrics[j].Type })
	return metrics, nil
}

func (service *Service) PingDB() error {
	if s, ok := service.repository.(*storage.DBStorage); ok {
		return s.Ping()
	}
	return models.ErrNotDB
}

func (service *Service) Close() {
	service.repository.Close()
}
