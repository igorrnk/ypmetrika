package storage

import (
	"github.com/igorrnk/ypmetrika/internal/models"
	"reflect"
	"testing"
)

func TestMemStorage_Write(t *testing.T) {
	type fields struct {
		Metrics map[string]models.Metric
	}
	type args struct {
		metric models.Metric
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
				metric: models.Metric{
					Name:  "TestMetric",
					Type:  models.CounterType,
					Value: models.NewValue(0, 1000),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := NewAgentStorage()
			wantMemStorage := &MemoryStorage{metrics: map[string]*models.Metric{tt.args.metric.Name: &tt.args.metric}}
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
