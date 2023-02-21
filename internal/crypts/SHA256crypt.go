package crypts

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/igorrnk/ypmetrika/internal/models"
	"log"
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
	metric.Hash = hash(createString(metric), c.Key)
}

func (c CrypterSHA256) CheckHash(metric *models.Metric) error {
	if c.Key == "" {
		return nil
	}
	hashMetric := hash(createString(metric), c.Key)
	if hashMetric != metric.Hash {
		log.Printf("wrong metric: recieve %s, must %s", metric.Hash, hashMetric)
		return models.ErrWrongHash
	}
	return nil
}

func createString(metric *models.Metric) string {
	var s string
	switch metric.Type {
	case models.GaugeType:
		s = fmt.Sprintf("%s:gauge:%f", metric.Name, metric.Gauge)
	case models.CounterType:
		s = fmt.Sprintf("%s:counter:%d", metric.Name, metric.Counter)
	}
	return s
}

func hash(s string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
