package storage

import (
	"github.com/igorrnk/ypmetrika/internal/metrics"
	"reflect"
	"testing"
)

func TestMemStorage_Write(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]float64
		CounterMetrics map[string]int64
	}
	type args struct {
		metric *metrics.Metric
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantFields fields
	}{
		{
			name: "Good test #1",
			args: args{
				metric: &metrics.Metric{
					Name:   "TestMetric",
					Type:   "counter",
					Value:  "100",
					Source: "",
				},
			},
			wantErr: false,
			wantFields: fields{
				CounterMetrics: map[string]int64{
					"TestMetric": int64(100)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := &MemStorage{
				GaugeMetrics:   tt.fields.GaugeMetrics,
				CounterMetrics: tt.fields.CounterMetrics,
			}
			wantMemStorage := &MemStorage{
				GaugeMetrics:   tt.wantFields.GaugeMetrics,
				CounterMetrics: tt.wantFields.CounterMetrics,
			}
			err := memStorage.Write(tt.args.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if reflect.DeepEqual(memStorage, wantMemStorage) {
				t.Errorf("Write() memStorage = %v, wantMemStorage %v", memStorage, wantMemStorage)
			}

		})
	}
}
