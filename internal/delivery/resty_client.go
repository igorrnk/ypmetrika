package delivery

import (
	"encoding/json"
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
	return &RestyClient{
		Client:        resty.New(),
		AddressServer: config.AddressServer,
	}
}

func (client RestyClient) Post(metric *models.Metric) {

	url := fmt.Sprintf("http://%s/update/%s/%s/%s",
		client.AddressServer, metric.Type, metric.Name, metric.Value)
	resp, err := client.Client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)
	if err != nil {
		log.Println(err)
	}
	log.Printf("POST %v Status: %v", url, resp.Status())
}

func (client RestyClient) PostJSON(metric *models.Metric) {

	url := fmt.Sprintf("http://%s/update/",
		client.AddressServer)
	body, err := json.Marshal(metric)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := client.Client.R().
		SetBody(body).
		SetHeader("Content-Type", "application/json").
		Post(url)
	if err != nil {
		log.Println(err)
	}
	log.Printf("POST %v Status: %v", url, resp.Status())
}
