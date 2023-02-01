package handlers

import (
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	Config *configs.ServerConfig
	Server models.ServerUsecase
}

func NewHandler(config *configs.ServerConfig, serverUsecase models.ServerUsecase) *Handler {
	return &Handler{
		Config: config,
		Server: serverUsecase,
	}
}

func (h Handler) HandleFn(w http.ResponseWriter, r *http.Request) {
	list, err := h.Server.GetAll()
	if err != nil {
		http.Error(w, "Metric reading error", http.StatusInternalServerError)
	}
	page := models.Page{
		Tittle: "GetAll metrics",
		List:   list,
	}
	w.Header().Add("Content-Type", "text/html")
	//w.Header().Add("Content-Encoding", "gzip")
	t, err := template.ParseFiles(h.Config.NameHTMLFile)
	if err != nil {
		log.Println(err)
		http.Error(w, "Wrong html template", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, page)
	if err != nil {
		log.Println(err)
		http.Error(w, "Wrong html template", http.StatusInternalServerError)
	}
	log.Printf("Request %v has been handled.", r.RequestURI)
}
