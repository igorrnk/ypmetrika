package models

import (
	"errors"
	"fmt"
	"strconv"
)

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
	Name   string
	Type   MetricType
	Value  Value
	Source SourceType
}

func (metric Metric) ValueS() string {
	var s string
	switch metric.Type {
	case Gauge:
		s = fmt.Sprint(metric.Value.Gauge)
	case Counter:
		s = fmt.Sprint(metric.Value.Counter)
	}
	return s
}

type MetricType int

func (d MetricType) String() string {
	return [...]string{"gauge", "counter"}[d]
}

const (
	Gauge MetricType = iota
	Counter
)

type SourceType int

const (
	RuntimeSource SourceType = iota
	CounterSource
	RandomSource
)

func ToMetricType(s string) (MetricType, error) {
	if metricType, ok := map[string]MetricType{
		"gauge":   Gauge,
		"counter": Counter,
	}[s]; ok {
		return metricType, nil
	} else {
		return 0, errors.New("models.ToMetricType: converting \"%s\": invalid syntax")
	}
}

func ToValue(s string, metricType MetricType) (Value, error) {
	switch metricType {
	case Gauge:
		gauge, err := strconv.ParseFloat(s, 64)
		return Value{Gauge: gauge}, err
	case Counter:
		counter, err := strconv.ParseInt(s, 10, 64)
		return Value{Counter: counter}, err
	}
	return Value{}, fmt.Errorf("models.ToValue: converting %v, %v: invalid metricType", s, metricType)
}
