package delivery

import "github.com/go-resty/resty/v2"

type Job struct {
	Request  *resty.Request
	URL      string
	ResultCh chan *ResultJob
}

type ResultJob struct {
	Response *resty.Response
	Err      error
}

func NewJob(r *resty.Request, url string) *Job {
	return &Job{r, url, make(chan *ResultJob)}
}

func (job *Job) Result() (*resty.Response, error) {
	result := <-job.ResultCh
	close(job.ResultCh)
	return result.Response, result.Err
}
