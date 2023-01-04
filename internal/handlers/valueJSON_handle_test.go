package handlers

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/igorrnk/ypmetrika/configs"
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

func TestHandler_ValueJSONHandleFn(t *testing.T) {
	type fields struct {
		Config configs.ServerConfig
		Server models.ServerUsecase
	}
	type request struct {
		reauestURI  string
		contentType string
		body        []byte
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
		contentType string
		body        []byte
	}
	tests := []struct {
		name     string
		fields   fields
		request  request
		mockArgs mockArgs
		want     want
	}{
		{
			name: "ValueJSONGaugeAlloc",
			request: request{
				reauestURI:  "/value/",
				contentType: "application/json",
				body:        []byte(`{"id":"Alloc","type":"gauge"}`),
			},
			fields: fields{
				Config: configs.DefaultServerConfig,
				Server: &test.ServerMock{},
			},
			mockArgs: mockArgs{arg0: "Value",
				arg1: models.Metric{
					Name: "Alloc",
					Type: models.GaugeType,
				},
				ret0: models.Metric{
					Name:  "Alloc",
					Type:  models.GaugeType,
					Value: models.Value{Gauge: 123456.789},
				},
				ret1: true,
			},
			want: want{
				code:        http.StatusOK,
				body:        []byte(`{"id":"Alloc","type":"gauge","delta":0,"value":123456.789}`),
				contentType: "application/json",
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
			router.Get("/value/", h.ValueJSONHandleFn)

			request := httptest.NewRequest(http.MethodGet, tt.request.reauestURI, bytes.NewReader(tt.request.body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			res := w.Result()
			body, _ := io.ReadAll(res.Body)
			_ = res.Body.Close()

			assert := assert.New(t)
			assert.Equal(tt.want.code, res.StatusCode)
			assert.Equal(tt.want.contentType, res.Header.Get("Content-Type"))
			assert.JSONEq(string(tt.want.body), string(body))

		})
	}
}
