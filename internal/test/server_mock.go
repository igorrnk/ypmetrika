package test

import (
	"github.com/igorrnk/ypmetrika/internal/models"
	"github.com/stretchr/testify/mock"
	"log"
)

type ServerMock struct {
	mock.Mock
}

func (mock *ServerMock) UpdateValue(metric models.Metric) (models.Metric, error) {
	args := mock.Called(metric)
	return args.Get(0).(models.Metric), args.Error(1)
}

func (mock *ServerMock) Update(metric models.Metric) error {
	args := mock.Called(metric)
	return args.Error(0)
}

func (mock *ServerMock) Value(metric models.Metric) (models.Metric, bool) {
	log.Printf("Called ServerMock.Value(%v)\n", metric)
	args := mock.Called(metric)
	return args.Get(0).(models.Metric), args.Get(1).(bool)
}

func (mock *ServerMock) GetAll() []models.Metric {
	args := mock.Called()
	return args.Get(0).([]models.Metric)
}
