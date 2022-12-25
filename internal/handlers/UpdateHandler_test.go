package handlers

import (
	"github.com/igorrnk/ypmetrika/internal/storage"
	"net/http"
	"testing"
)

func TestUpdateHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		Rep storage.Repositories
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := UpdateHandler{
				Rep: tt.fields.Rep,
			}
			h.ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}
