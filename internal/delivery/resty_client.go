package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"log"
)

type RestyClient struct {
	Client        *resty.Client
	AddressServer string
}

func NewRestyClient(config *configs.AgentConfig) *RestyClient {
	return &RestyClient{
		Client:        resty.New(),
		AddressServer: config.AddressServer,
	}
}

func (client RestyClient) Post(metric *models.Metric) {

	url := fmt.Sprintf("%s/update/%s/%s/%s",
		client.AddressServer, metric.Type, metric.Name, metric.Value())
	resp, err := client.Client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)
	if err != nil {
		log.Println(err)
	}
	log.Printf("POST %v Status: %v", url, resp.Status())
}

func (client RestyClient) PostJSON(metric *models.Metric) {

	url := fmt.Sprintf("%s/update/",
		client.AddressServer)
	body, err := json.Marshal(metric)
	if err != nil {
		log.Printf("client.PostJSON: error: %v\n", err)
		return
	}
	resp, err := client.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
	if err != nil {
		log.Printf("client.PostJSON: error: %v\n", err)
	}
	log.Printf("client.PostJSON: URL = %v\n", url)
	log.Printf("client.PostJSON: BODY = %v\n", string(body))
	if resp != nil {
		log.Printf("POST %v Status: %v\n", url, resp.Status())
	} else {
		log.Printf("POST %v Status: no response\n", url)
	}
}
