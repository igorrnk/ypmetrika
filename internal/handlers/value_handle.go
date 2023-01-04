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
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	metric, ok := h.Server.Value(models.Metric{Name: nameMetric, Type: typeMetric})
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprint(metric.Value)))
	if err != nil {
		log.Println(err)
	}
	log.Printf("Request %v has been handled.", r.RequestURI)
}
