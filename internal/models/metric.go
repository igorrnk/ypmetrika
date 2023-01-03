package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Gauge float64
type Counter int64

type Value struct {
	Gauge   float64
	Counter int64
}

func (value Value) String() string {
	switch {
	case value.Gauge != 0:
		return fmt.Sprint(value.Gauge)
	case value.Counter != 0:
		return fmt.Sprint(value.Counter)
	}
	return "0"
}

type Metric struct {
	Name   string `json:"id"`
	Type   MetricType
	Value  Value
	Source SourceType
}

type JSONMetric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (metric *Metric) UnmarshalJSON(bytes []byte) error {
	type MetricAlias Metric
	aliasValue := JSONMetric{}
	var err error
	if err = json.Unmarshal(bytes, &aliasValue); err != nil {
		return err
	}
	metric.Name = aliasValue.ID
	if metric.Type, err = ToMetricType(aliasValue.MType); err != nil {
		return err
	}
	metric.Value.Counter = *aliasValue.Delta
	metric.Value.Gauge = *aliasValue.Value
	return nil
}

func (metric Metric) MarshalJSON() ([]byte, error) {
	type MetricAlias Metric
	aliasValue := JSONMetric{
		ID:    metric.Name,
		MType: metric.Type.String(),
		Delta: &metric.Value.Counter,
		Value: &metric.Value.Gauge,
	}
	return json.Marshal(aliasValue)
}

func (metric Metric) ValueS() string {
	var s string
	switch metric.Type {
	case GaugeType:
		s = fmt.Sprint(metric.Value.Gauge)
	case CounterType:
		s = fmt.Sprint(metric.Value.Counter)
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
		return 0, errors.New("models.ToMetricType: converting \"%s\": invalid syntax")
	}
}

func ToValue(s string, metricType MetricType) (Value, error) {
	switch metricType {
	case GaugeType:
		gauge, err := strconv.ParseFloat(s, 64)
		return Value{Gauge: gauge}, err
	case CounterType:
		counter, err := strconv.ParseInt(s, 10, 64)
		return Value{Counter: counter}, err
	}
	return Value{}, fmt.Errorf("models.ToValue: converting %v, %v: invalid metricType", s, metricType)
}
