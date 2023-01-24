package handlers

import (
	"github.com/igorrnk/ypmetrika/internal/configs"
	"github.com/igorrnk/ypmetrika/internal/models"
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	Config configs.ServerConfig
	Server models.ServerUsecase
}

func NewHandler(config configs.ServerConfig, serverUsecase models.ServerUsecase) *Handler {
	return &Handler{
		Config: config,
		Server: serverUsecase,
	}
}

func (h Handler) HandleFn(w http.ResponseWriter, r *http.Request) {
	page := models.Page{
		Tittle: "GetAll metrics",
	}

	list, err1 := h.Server.GetAll()
	if err1 != nil {
		log.Println(err1)
	}
	page.List = list

	w.Header().Add("Content-Type", "text/html")
	t, err := template.ParseFiles(h.Config.NameHTMLFile)
	err = t.Execute(w, page)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Request %v has been handled.", r.RequestURI)
}
