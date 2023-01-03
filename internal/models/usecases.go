package models

type ServerUsecase interface {
	Update(metric Metric) error
	Value(metric Metric) (Metric, bool)
	GetAll() []Metric
}

type Client interface {
	Post(*Metric)
	PostJSON(*Metric)
}

type Repository interface {
	Write(Metric) error
	Read(Metric) (Metric, bool)
	ReadAll() ([]Metric, error)
}
