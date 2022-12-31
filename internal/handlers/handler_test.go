package handlers

import (
	"github.com/igorrnk/ypmetrika/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_UpdateHandleFn(t *testing.T) {
	type fields struct {
		Config  configs.ServerConfig
		Usecase models.Usecase
	}
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		request string
		fields  fields

		want want
	}{
		{
			name:    "Invalid type",
			request: "/update/unknown/testCounter/100",
			fields: fields{
				Config: configs.DefaultServerConfig,
			},
			want: want{

				code: http.StatusNotImplemented,
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := Handler{
				Config: tt.fields.Config,
			}
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			hf := http.HandlerFunc(h.UpdateHandleFn)
			// запускаем сервер
			hf.ServeHTTP(w, request)
			res := w.Result()
			res.Body.Close()
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
		})
	}
}
