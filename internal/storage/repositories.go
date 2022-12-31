package storage

import (
	"github.com/igorrnk/ypmetrika/internal/models"
)

type AgentRepository interface {
	Update(models.AgentMetric)
	Fill(fileName string) (err error)
	ReadAll() []models.AgentMetric
}

type ServerRepository interface {
	Write(metric models.ServerMetric) error
	Read(metric models.ServerMetric) (models.ServerMetric, bool)
	ReadAll() []models.ServerMetric
}
