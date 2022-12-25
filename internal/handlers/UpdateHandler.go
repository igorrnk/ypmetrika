package handlers

import (
	"github.com/igorrnk/ypmetrika/internal/metrics"
	"github.com/igorrnk/ypmetrika/internal/storage"
	"log"
	"net/http"
)

type UpdateHandler struct {
	Rep storage.Repositories
}

func (h UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Server gets not POST method.")
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("URL path = %v", r.RequestURI)
	metric := new(metrics.Metric)
	err := metric.URLtoMetric(r.RequestURI)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if metric.Type != "counter" && metric.Type != "gauge" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	err = h.Rep.Write(metric)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Metric %v %v = %v has been written.", metric.Name, metric.Type, metric.Value)
	w.WriteHeader(http.StatusOK)
}
