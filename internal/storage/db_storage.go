package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	context context.Context
	db      *sql.DB
}

func NewDBStorage(ctx context.Context, config *configs.ServerConfig) (models.Repository, error) {
	storage := &DBStorage{
		context: ctx,
	}
	var err error
	storage.db, err = sql.Open(config.DBDriverName, config.DBConnect)
	if err != nil {
		log.Printf("NewDBStorage: failed open database: %v", err)
		return nil, err
	}

	_, err = storage.db.Exec("CREATE TABLE IF NOT EXISTS metrics (name varchar(30) UNIQUE," +
		" typem smallint, gauge float, counter bigint)")
	if err != nil {
		log.Printf("NewDBStorage: failed create table: %v", err)
		return nil, err
	}
	return storage, nil
}

func (storage *DBStorage) Write(metric *models.Metric) error {
	var err error
	switch metric.Type {
	case models.GaugeType:
		_, err = storage.db.Exec("INSERT INTO metrics (name, typem, gauge, counter) VALUES ($1,$2,$3,$4)"+
			" ON CONFLICT (name) DO UPDATE SET gauge = EXCLUDED.gauge",
			metric.Name, metric.Type, metric.Gauge, metric.Counter)
	case models.CounterType:
		_, err = storage.db.Exec("INSERT INTO metrics (name, typem, gauge, counter) VALUES ($1,$2,$3,$4) "+
			"ON CONFLICT (name) DO UPDATE SET counter = EXCLUDED.counter",
			metric.Name, metric.Type, metric.Gauge, metric.Counter)
	}

	if err != nil {
		log.Printf("DBStorage.Write: %v", err)
	}
	return err
}
func (storage *DBStorage) Read(metric *models.Metric) (*models.Metric, error) {
	var err error
	switch metric.Type {
	case models.GaugeType:
		err = storage.db.QueryRow("SELECT gauge FROM metrics WHERE name = $1", metric.Name).Scan(&metric.Gauge)
	case models.CounterType:
		err = storage.db.QueryRow("SELECT counter FROM metrics WHERE name = $1", metric.Name).Scan(&metric.Counter)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		log.Printf("DBStorage.Read: %v", err)
	}
	return metric, err
}
func (storage *DBStorage) ReadAll() ([]models.Metric, error) {
	metrics := make([]models.Metric, 0)
	rows, err := storage.db.Query("SELECT * FROM metrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var m models.Metric
		err = rows.Scan(&m.Name, &m.Type, &m.Gauge, &m.Counter)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, m)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (storage *DBStorage) Ping() error {
	return storage.db.Ping()
}
func (storage *DBStorage) Close() {
	storage.db.Close()
}
