package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"log"
	"net/http"
)

type RestyClient struct {
	Client        *resty.Client
	AddressServer string
	jobCh         chan *Job
}

func NewRestyClient(config *configs.AgentConfig) *RestyClient {
	client := &RestyClient{
		Client:        resty.New(),
		AddressServer: config.AddressServer,
		jobCh:         make(chan *Job),
	}
	client.runWorkers(config.Limit)
	return client
}

func (client *RestyClient) runWorkers(count int) {
	for i := 0; i < count; i++ {
		go func() {
			for job := range client.jobCh {
				resp, err := func(r *resty.Request, url string) (*resty.Response, error) {
					return r.Post(url)
				}(job.Request, job.URL)
				job.ResultCh <- &ResultJob{resp, err}
			}
		}()
	}
}

func (client *RestyClient) Close() {
	close(client.jobCh)
}

func (client *RestyClient) AddPostRequest(r *resty.Request, url string) (*resty.Response, error) {
	job := NewJob(r, url)
	client.jobCh <- job
	return job.Result()
}

func (client *RestyClient) Post(metric *models.Metric) error {
	url := fmt.Sprintf("%s/update/%s/%s/%s",
		client.AddressServer, metric.Type, metric.Name, metric.Value())
	resp, err := client.AddPostRequest(
		client.Client.R().SetHeader("Content-Type", "text/plain"),
		url)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode() != http.StatusOK {
			return models.ErrNotReport
		}
	}
	return nil
}

func (client *RestyClient) PostJSON(metric *models.Metric) error {

	url := fmt.Sprintf("%s/update/",
		client.AddressServer)
	body, err := json.Marshal(metric)
	if err != nil {
		log.Printf("client.PostJSON: error: %v\n", err)
		return nil
	}
	resp, err := client.AddPostRequest(
		client.Client.R().SetHeader("Content-Type", "application/json").SetBody(body),
		url)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode() != http.StatusOK {
			return models.ErrNotReport
		}
	}
	return nil
}

func (client *RestyClient) PostMetrics(metrics []models.Metric) error {
	url := fmt.Sprintf("%s/updates/",
		client.AddressServer)
	body, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("client.PostJSON: error: %v\n", err)
		return err
	}
	resp, err := client.AddPostRequest(
		client.Client.R().SetHeader("Content-Type", "application/json").SetBody(body),
		url)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode() != http.StatusOK {
			return models.ErrNotReport
		}
	}
	return nil
}
