package models

type ServerUsecase interface {
	UpdateValue(metric Metric) (Metric, error)
	Update(metric Metric) error
	Value(metric Metric) (Metric, error)
	GetAll() ([]Metric, error)
}

type Client interface {
	Post(*Metric)
	PostJSON(*Metric)
}

type Repository interface {
	Write(Metric) error
	Read(Metric) (Metric, error)
	ReadAll() ([]Metric, error)
}
