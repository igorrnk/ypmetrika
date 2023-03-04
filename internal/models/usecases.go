package models

type ServerUsecase interface {
	UpdateValue(metric *Metric) (*Metric, error)
	Update(metric *Metric) error
	Updates([]*Metric) error
	Value(metric *Metric) (*Metric, error)
	PingDB() error
	//GetAll returns slice of all metrics
	GetAll() ([]Metric, error)
	Close()
}

type Client interface {
	Post(*Metric) error
	PostJSON(*Metric) error
	PostMetrics([]Metric) error
	Close()
}

type Repository interface {
	Write(*Metric) error
	Read(*Metric) (*Metric, error)
	ReadAll() ([]Metric, error)
	Close()
}

type Crypter interface {
	AddHash(metric *Metric)
	CheckHash(metric *Metric) error
}
