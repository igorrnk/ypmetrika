package agents

import (
	"github.com/igorrnk/ypmetrika/internal/metrics"
	"reflect"
	"testing"
)

func TestNewAgent(t *testing.T) {
	tests := []struct {
		name string
		want Agent
	}{
		{
			name: "TestNewAgent #1",
			want: Agent{
				Config:  DefaultConfig,
				Metrics: &metrics.Metrics{},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAgent()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAgent() = %v, want %v", got, tt.want)
			}
			//t.Logf("NewAgent() = %v, want %v", got, tt.want)
		})
	}
}
