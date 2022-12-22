package handler

import (
	"github.com/igorrnk/ypmetrika/internal/metrics"
	"github.com/igorrnk/ypmetrika/internal/storage"
	"log"
	"net/http"
)

type Handler struct {
	Rep storage.Repositories
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Server gets not POST method.")
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		log.Println("Server gets not text/plain content-type.")
		http.Error(w, "Wrong Content-Type", http.StatusBadRequest)
		return
	}

	log.Printf("URL path = %v", r.RequestURI)
	metric := new(metrics.Metric)
	err := metric.URLtoMetric(r.RequestURI)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	err = h.Rep.Write(metric)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Printf("Metric %v %v = %v has been written.", metric.Name, metric.Type, metric.Value)
	w.WriteHeader(http.StatusOK)
}
