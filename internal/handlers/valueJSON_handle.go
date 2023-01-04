package handlers

import (
	"encoding/json"
	"github.com/igorrnk/ypmetrika/internal/models"
	"io"
	"log"
	"net/http"
)

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
