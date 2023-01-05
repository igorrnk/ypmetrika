package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/igorrnk/ypmetrika/internal/models"
	"log"
	"net/http"
)

func (h Handler) UpdateHandleFn(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")
	typeMetric, err := models.ToMetricType(chi.URLParam(r, "typeMetric"))
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	valueMetric, err := models.ToValue(chi.URLParam(r, "valueMetric"), typeMetric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metric := models.Metric{Name: nameMetric, Type: typeMetric, Value: valueMetric}
	if err = h.Server.Update(metric); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	log.Printf("Request %v has been handled.", r.RequestURI)
}
