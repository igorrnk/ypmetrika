package crypts

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/igorrnk/ypmetrika/internal/models"
)

type CrypterSHA256 struct {
	Key string
}

func NewCrypterSHA256(key string) models.Crypter {
	return CrypterSHA256{
		Key: key,
	}
}

func (c CrypterSHA256) AddHash(metric *models.Metric) {
	if c.Key == "" {
		return
	}
	switch metric.Type {
	case models.GaugeType:
		metric.Hash = hash(fmt.Sprintf("%s:counter:%s", metric.Name, metric.Value()), c.Key)
	case models.CounterType:
		metric.Hash = hash(fmt.Sprintf("%s:gauge:%s", metric.Name, metric.Value()), c.Key)
	}
}

func (c CrypterSHA256) CheckHash(metric *models.Metric) error {
	if c.Key == "" {
		return nil
	}

	var hashMetric string
	switch metric.Type {
	case models.GaugeType:
		hashMetric = hash(fmt.Sprintf("%s:counter:%s", metric.Name, metric.Value()), c.Key)
	case models.CounterType:
		hashMetric = hash(fmt.Sprintf("%s:gauge:%s", metric.Name, metric.Value()), c.Key)
	}
	if hashMetric != metric.Hash {
		return models.ErrWrongHash
	}
	return nil
}

func hash(s string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))

}
