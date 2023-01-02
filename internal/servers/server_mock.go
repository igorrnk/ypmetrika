package servers

import (
	"github.com/igorrnk/ypmetrika/internal/models"
	"github.com/stretchr/testify/mock"
)

type UsecaseMock struct {
	mock.Mock
}

func (m *UsecaseMock) Update(metric models.Metric) error {
	args := m.Called(metric)
	return args.Error(0)
}

func (m *UsecaseMock) Value(metric models.Metric) (models.Metric, bool) {
	args := m.Called(metric)
	return args.Get(0).(models.Metric), args.Get(1).(bool)

}

func (m *UsecaseMock) All() []models.Metric {
	return nil
}
