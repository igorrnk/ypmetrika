package models

type Crypter interface {
	AddHash(metric *Metric)
	CheckHash(metric *Metric) error
}
