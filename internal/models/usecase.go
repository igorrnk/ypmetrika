package models

type Usecase interface {
	Update(metric ServerMetric)
	Value(metric ServerMetric) (ServerMetric, bool)
	All() []ServerMetric
}
