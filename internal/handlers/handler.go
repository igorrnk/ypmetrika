package handlers

import (
	"github.com/igorrnk/ypmetrika/configs"
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
		List:   h.Server.GetAll(),
	}

	t, _ := template.ParseFiles(h.Config.NameHTMLFile)
	err := t.Execute(w, page)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Request %v has been handled.", r.RequestURI)

}
