package delivery

import (
	"context"
	"github.com/go-resty/resty/v2"
	configs2 "github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestRestyClient_PostJSON(t *testing.T) {
	type fields struct {
		Client        *resty.Client
		AddressServer string
	}
	type args struct {
		metric *models.Metric
	}
	type want struct {
		requestURI  string
		contentType string
		body        []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Gauge",
			fields: fields{
				Client:        resty.New(),
				AddressServer: configs2.DefaultAC.AddressServer,
			},
			args: args{
				metric: &models.Metric{
					Name:  "Alloc",
					Type:  models.GaugeType,
					Value: models.Value{Gauge: 123456.789},
				},
			},
			want: want{
				requestURI:  "/update/",
				contentType: "application/json",
				body:        []byte(`{"id":"Alloc","type":"gauge","value":123456.789}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := RestyClient{
				Client:        tt.fields.Client,
				AddressServer: tt.fields.AddressServer,
			}
			server := &http.Server{}
			var got want
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				got.requestURI = r.RequestURI
				got.body, _ = io.ReadAll(r.Body)
				got.contentType = r.Header.Get("Content-Type")
				r.Body.Close()
				w.WriteHeader(http.StatusOK)
			})
			go http.ListenAndServe(configs2.DefaultSC.AddressServer, nil)
			time.Sleep(1 * time.Second)
			client.PostJSON(tt.args.metric)
			time.Sleep(1 * time.Second)
			server.Shutdown(context.TODO())
			assert.Equal(t, tt.want, got)
		})
	}
}
