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
		http.Error(w, "Wrong metric type", http.StatusNotImplemented)
		return
	}

	metric, err1 := models.ToMetric(nameMetric, chi.URLParam(r, "valueMetric"), typeMetric)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	if err = h.Server.Update(metric); err != nil {
		http.Error(w, "Metric not found", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	log.Printf("Request %v has been handled.", r.RequestURI)
}
