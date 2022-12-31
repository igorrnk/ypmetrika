package storage

import (
	"github.com/igorrnk/ypmetrika/internal/models"
	"reflect"
	"testing"
)

func TestMemStorage_Write(t *testing.T) {
	type fields struct {
		GaugeMetrics   map[string]float64
		CounterMetrics map[string]int64
	}
	type args struct {
		metric models.ServerMetric
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantFields fields
	}{
		{
			name: "Good test #1",
			args: args{
				metric: models.ServerMetric{
					Name:  "TestMetric",
					Type:  "counter",
					Value: "100",
				},
			},

			wantErr: false,
			wantFields: fields{
				CounterMetrics: map[string]int64{
					"TestMetric": int64(100),
				},
				GaugeMetrics: make(map[string]float64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := NewServerMemoryStorage()
			wantMemStorage := &ServerMemoryStorage{
				GaugeMetrics:   tt.wantFields.GaugeMetrics,
				CounterMetrics: tt.wantFields.CounterMetrics,
			}
			err := memStorage.Write(tt.args.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(memStorage, wantMemStorage) {
				t.Errorf("WriteMetric(): memStorage = %v, wantMemStorage %v", memStorage, wantMemStorage)
			}

		})
	}
}
