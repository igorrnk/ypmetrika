package models

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetric_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Name    string
		Type    MetricType
		Gauge   float64
		Counter int64
		Source  SourceType
	}
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    Metric
	}{
		{
			name:   "Good Gauge",
			fields: fields{},
			args: args{
				bytes: []byte(`{"id":"Alloc","type":"gauge","value":123456.789}`),
			},
			wantErr: false,
			want: Metric{
				Name:   "Alloc",
				Type:   GaugeType,
				Gauge:  123456.789,
				Source: RuntimeSource,
			},
		},
		{
			name:   "Good Counter",
			fields: fields{},
			args: args{
				bytes: []byte(`{"id":"PollCount","type":"counter","delta":123}`),
			},
			wantErr: false,
			want: Metric{
				Name:    "PollCount",
				Type:    CounterType,
				Counter: 123,
				Source:  RuntimeSource,
			},
		},
		{
			name:   "Good Counter without value",
			fields: fields{},
			args: args{
				bytes: []byte(`{"id":"PollCount","type":"counter"}`),
			},
			wantErr: false,
			want: Metric{
				Name: "PollCount",
				Type: CounterType,
			},
		},
		{
			name:   "Bad Counter",
			fields: fields{},
			args: args{
				bytes: []byte(`{"id":"PollCount","type":"counter","delta":123.9,"value":0}`),
			},
			wantErr: true,
			want:    Metric{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := &Metric{}
			err := metric.UnmarshalJSON(tt.args.bytes)
			require.True(t, (err != nil) == tt.wantErr, "(err != nil)")
			assert.Equal(t, tt.want, *metric)
		})
	}
}

func TestMetric_MarshalJSON(t *testing.T) {
	type fields struct {
		Name    string
		Type    MetricType
		Gauge   float64
		Counter int64
		Source  SourceType
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Good Gauge",
			fields: fields{
				Name:   "Alloc",
				Type:   GaugeType,
				Gauge:  123456.789,
				Source: RuntimeSource,
			},
			want:    `{"id":"Alloc", "type":"gauge", "value":123456.789}`,
			wantErr: false,
		},
		{
			name: "Good Counter",
			fields: fields{
				Name:    "Alloc",
				Type:    CounterType,
				Counter: 123,
				Source:  RuntimeSource,
			},
			want:    `{"id":"Alloc","type":"counter","delta":123}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := Metric{
				Name:    tt.fields.Name,
				Type:    tt.fields.Type,
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
				Source:  tt.fields.Source,
			}
			got, err := metric.MarshalJSON()
			require.True(t, (err != nil) == tt.wantErr, "(err != nil)")
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}
