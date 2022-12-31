package models

type Usecase interface {
	Update(metric ServerMetric) error
	Value(metric ServerMetric) (ServerMetric, bool)
	All() []ServerMetric
}
