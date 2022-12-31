package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/igorrnk/ypmetrika/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"log"
	"net/http"
	"text/template"
)

type Handler struct {
	Config  configs.ServerConfig
	Usecase models.Usecase
}

func NewHandler(config configs.ServerConfig, usecase models.Usecase) *Handler {
	return &Handler{
		Config:  config,
		Usecase: usecase,
	}
}

func (h Handler) HandleFn(w http.ResponseWriter, r *http.Request) {
	metrics := h.Usecase.All()
	page := models.Page{
		Tittle: "All metrics",
	}
	list := ""
	for _, metric := range metrics {
		list += fmt.Sprintf(`<li> %v: %v </li>`, metric.Name, metric.Value) + "\n"
	}
	page.List = list

	t, _ := template.ParseFiles(h.Config.NameHTMLFile)
	t.Execute(w, page)

}

func (h Handler) UpdateHandleFn(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateHandlerFn(): running.")
	nameMetric := chi.URLParam(r, "nameMetric")
	typeMetric := chi.URLParam(r, "typeMetric")
	valueMetric := chi.URLParam(r, "valueMetric")
	if typeMetric != "gauge" && typeMetric != "counter" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	metric := models.ServerMetric{Name: nameMetric, Type: typeMetric, Value: valueMetric}
	if err := h.Usecase.Update(metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h Handler) ValueHandleFn(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")
	typeMetric := chi.URLParam(r, "typeMetric")
	metric, ok := h.Usecase.Value(models.ServerMetric{Name: nameMetric, Type: typeMetric})
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(metric.Value)))
}
