package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/igorrnk/ypmetrika/internal/models"
	"log"
	"net/http"
)

func (h Handler) ValueHandleFn(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")
	typeMetric, err := models.ToMetricType(chi.URLParam(r, "typeMetric"))
	if err != nil {
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}
	metric, ok := h.Server.Value(models.Metric{Name: nameMetric, Type: typeMetric})
	if !ok {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprint(metric.Value)))
	log.Printf("Request %v has been handled.", r.RequestURI)
}
