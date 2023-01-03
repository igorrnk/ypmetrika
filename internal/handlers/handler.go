package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/igorrnk/ypmetrika/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"io"
	"log"
	"net/http"
	"text/template"
)

type Handler struct {
	Config configs.ServerConfig
	Server models.ServerUsecase
}

func NewHandler(config configs.ServerConfig, usecase models.ServerUsecase) *Handler {
	return &Handler{
		Config: config,
		Server: usecase,
	}
}

func (h Handler) HandleFn(w http.ResponseWriter, r *http.Request) {
	page := models.Page{
		Tittle: "GetAll metrics",
		List:   h.Server.GetAll(),
	}

	t, _ := template.ParseFiles(h.Config.NameHTMLFile)
	err := t.Execute(w, page)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Request %v has been handled.", r.RequestURI)

}

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
	w.WriteHeader(http.StatusOK)
	log.Printf("Request %v has been handled.", r.RequestURI)
}

func (h Handler) UpdateJSONHandleFn(w http.ResponseWriter, r *http.Request) {
	metric := models.Metric{}
	var body []byte
	var err error
	if body, err = io.ReadAll(r.Body); err != nil {
		log.Println(err)
		return
	}
	if err = r.Body.Close(); err != nil {
		log.Println(err)
	}
	if err = json.Unmarshal(body, &metric); err != nil {
		log.Println(err)
		return
	}
	if err = h.Server.Update(metric); err != nil {
		log.Println(err)
		return
	}
	metric, ok := h.Server.Value(metric)
	if !ok {
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(metric)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) ValueJSONHandleFn(w http.ResponseWriter, r *http.Request) {
	metric := models.Metric{}
	var body []byte
	var err error
	if body, err = io.ReadAll(r.Body); err != nil {
		log.Println(err)
		return
	}
	if err = r.Body.Close(); err != nil {
		log.Println(err)
	}
	if err = json.Unmarshal(body, &metric); err != nil {
		log.Println(err)
		return
	}
	metric, ok := h.Server.Value(metric)
	if !ok {
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(metric)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

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
