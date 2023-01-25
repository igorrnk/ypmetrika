package servers

import (
	"github.com/igorrnk/ypmetrika/internal/models"
	"github.com/stretchr/testify/mock"
)

type ServerMock struct {
	mock.Mock
}

func (mock *ServerMock) Update(metric models.Metric) error {
	args := mock.Called(metric)
	return args.Error(0)
}

func (mock *ServerMock) Value(metric models.Metric) (models.Metric, bool) {
	args := mock.Called(metric)
	return args.Get(0).(models.Metric), args.Get(1).(bool)

}

func (mock *ServerMock) All() []models.Metric {
	return nil
}
