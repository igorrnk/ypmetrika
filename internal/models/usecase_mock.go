package models

import "github.com/stretchr/testify/mock"

type UsecaseMock struct {
	mock.Mock
}

func (m *UsecaseMock) Update(metric ServerMetric) error {
	args := m.Called(metric)
	return args.Error(0)
}

func (m *UsecaseMock) Value(metric ServerMetric) (ServerMetric, bool) {
	args := m.Called(metric)
	return args.Get(0).(ServerMetric), args.Get(1).(bool)

}

func (m *UsecaseMock) All() []ServerMetric {
	return nil
}
