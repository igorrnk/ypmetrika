package delivery

import "github.com/igorrnk/ypmetrika/internal/models"

type Client interface {
	Post(*models.AgentMetric)
}
