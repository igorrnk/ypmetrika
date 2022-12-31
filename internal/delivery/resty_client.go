package delivery

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/igorrnk/ypmetrika/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"log"
)

type RestyClient struct {
	Client        *resty.Client
	AddressServer string
}

func NewRestyClient(config configs.AgentConfig) *RestyClient {
	client := resty.New()
	return &RestyClient{
		Client:        client,
		AddressServer: config.AddressServer,
	}
}

func (client RestyClient) Post(metric *models.AgentMetric) {
	url := fmt.Sprintf("http://%s/update/%s/%s/%s",
		client.AddressServer,
		metric.Type,
		metric.Name,
		metric.Value)
	resp, err := client.Client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)
	if err != nil {
		log.Println(err)
	}
	log.Printf("POST %v Status: %v", url, resp.Status())
}
