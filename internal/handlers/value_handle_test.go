package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"github.com/igorrnk/ypmetrika/internal/test"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler_ValueHandleFn(t *testing.T) {
	type fields struct {
		Config *configs.ServerConfig
		Server models.ServerUsecase
	}
	type mockArgs struct {
		arg0 string
		arg1 any
		arg2 any
		ret0 any
		ret1 any
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name     string
		request  string
		mockArgs mockArgs
		fields   fields
		want     want
	}{
		{
			name:    "ValueGaugeAlloc",
			request: "/value/gauge/Alloc",
			fields: fields{
				Config: &configs.DefaultSC,
				Server: &test.ServerMock{},
			},
			mockArgs: mockArgs{
				arg0: "Value",
				arg1: &models.Metric{
					Name: "Alloc",
					Type: models.GaugeType,
				},
				ret0: &models.Metric{
					Name:  "Alloc",
					Type:  models.GaugeType,
					Gauge: 123456.789,
				},
				ret1: nil,
			},
			want: want{
				code:        http.StatusOK,
				response:    "123456.789",
				contentType: "text/plain",
			},
		},
		{
			name:    "ValueCounterPollCount",
			request: "/value/counter/PollCount",
			fields: fields{
				Config: &configs.DefaultSC,
				Server: &test.ServerMock{},
			},
			mockArgs: mockArgs{
				arg0: "Value",
				arg1: &models.Metric{
					Name: "PollCount",
					Type: models.CounterType,
				},
				ret0: &models.Metric{
					Name:    "PollCount",
					Type:    models.CounterType,
					Counter: 1234,
				},
				ret1: nil,
			},
			want: want{
				code:        http.StatusOK,
				response:    "1234",
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.SetOutput(os.Stdout)
			h := Handler{
				Config: tt.fields.Config,
				Server: tt.fields.Server,
			}
			serverMock := tt.fields.Server.(*test.ServerMock)
			serverMock.On(tt.mockArgs.arg0, tt.mockArgs.arg1).Return(tt.mockArgs.ret0, tt.mockArgs.ret1)

			router := chi.NewRouter()
			router.Get("/value/{typeMetric}/{nameMetric}", h.ValueHandleFn)

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			res := w.Result()
			body, _ := io.ReadAll(res.Body)
			res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, string(body))

		})
	}
}
