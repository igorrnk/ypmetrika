package metrics

import (
	"reflect"
	"testing"
)

func TestMetrics_Fill(t *testing.T) {
	ms := Metrics{}
	err := ms.Fill()
	if err != nil {
		t.Log(err)
	}
	t.Logf("Metrics: %v", len(ms.Metrics))
	t.Logf("ms.count: %v", ms.UpdateCount)
}

func TestMetrics_Update(t *testing.T) {
	ms := Metrics{}
	err := ms.Fill()
	if err != nil {
		t.Log(err)
	}
	_ = ms.Update()
	t.Logf("metric: %v %v\n", ms.Metrics[0].Name, ms.Metrics[0].Value)
}

func TestMetric_URLtoMetric(t *testing.T) {
	type fields struct {
		Name   string
		Type   string
		Value  string
		Source string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantFields fields
	}{
		{
			name: "Good path #1",
			args: args{
				path: "/update/gauge/RandomValue/1727040455672546632",
			},
			wantErr: false,
			wantFields: fields{
				Name:   "RandomValue",
				Type:   "gauge",
				Value:  "1727040455672546632",
				Source: "",
			},
		},
		{
			name: "Wrong path #1",
			args: args{
				path: "/update/gauge/RandomValue/",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := &Metric{
				Name:   tt.fields.Name,
				Type:   tt.fields.Type,
				Value:  tt.fields.Value,
				Source: tt.fields.Source,
			}
			metricWant := &Metric{
				tt.wantFields.Name,
				tt.wantFields.Type,
				tt.wantFields.Value,
				tt.wantFields.Source,
			}
			err := metric.URLtoMetric(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLtoMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(metric, metricWant) {
				t.Errorf("URLtoMetric() metric = %v, wantErr %v", metric, metricWant)
			}
		})
	}
}
