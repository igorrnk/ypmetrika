package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Metric struct {
	Name    string `json:"id"`
	Type    MetricType
	Gauge   float64
	Counter int64
	Source  SourceType
	Hash    string
}

type JSONMetric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func (metric *Metric) UnmarshalJSON(bytes []byte) error {
	type MetricAlias Metric
	aliasValue := JSONMetric{}
	var err error
	if err = json.Unmarshal(bytes, &aliasValue); err != nil {
		return err
	}
	metric.Name = aliasValue.ID
	metric.Hash = aliasValue.Hash
	if metric.Type, err = ToMetricType(aliasValue.MType); err != nil {
		return err
	}
	if aliasValue.Delta != nil {
		metric.Counter = *aliasValue.Delta
	}
	if aliasValue.Value != nil {
		metric.Gauge = *aliasValue.Value
	}
	return nil
}

func (metric *Metric) MarshalJSON() ([]byte, error) {
	aliasValue := JSONMetric{
		ID:    metric.Name,
		MType: metric.Type.String(),
		Hash:  metric.Hash,
	}
	switch metric.Type {
	case GaugeType:
		aliasValue.Value = &metric.Gauge
	case CounterType:
		aliasValue.Delta = &metric.Counter
	}
	return json.Marshal(aliasValue)
}

func (metric *Metric) Value() string {
	var s string
	switch metric.Type {
	case GaugeType:
		s = fmt.Sprintf("%s", metric.Gauge)
	case CounterType:
		s = fmt.Sprintf("%s", metric.Counter)
	}
	return s
}

type MetricType int

func (d MetricType) String() string {
	return [...]string{"gauge", "counter"}[d]
}

const (
	GaugeType MetricType = iota
	CounterType
)

type SourceType int

const (
	RuntimeSource SourceType = iota
	CounterSource
	RandomSource
)

func ToMetricType(s string) (MetricType, error) {
	if metricType, ok := map[string]MetricType{
		"gauge":   GaugeType,
		"counter": CounterType,
	}[s]; ok {
		return metricType, nil
	} else {
		return 0, fmt.Errorf("wrong metric type: \"%s\"", s)
	}
}

func ToMetric(name string, value string, metricType MetricType) (*Metric, error) {
	metric := &Metric{
		Name: name,
		Type: metricType,
	}
	var err error
	switch metricType {
	case GaugeType:
		metric.Gauge, err = strconv.ParseFloat(value, 64)
	case CounterType:
		metric.Counter, err = strconv.ParseInt(value, 10, 64)
	}
	if err != nil {
		return nil, fmt.Errorf("wrong value metric: %v", err)
	}
	return metric, nil
}
