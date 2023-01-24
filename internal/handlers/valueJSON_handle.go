package handlers

import (
	"encoding/json"
	"github.com/igorrnk/ypmetrika/internal/models"
	"io"
	"log"
	"net/http"
)

func (h Handler) ValueJSONHandleFn(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler.ValueJSONHandleFn: URL = %v\n", r.URL)
	metric := models.Metric{}
	var body []byte
	var err error
	if body, err = io.ReadAll(r.Body); err != nil {
		log.Printf("Handler.ValueJSONHandleFn: body ReadAll error: %v\n", err)
		return
	}
	if err = r.Body.Close(); err != nil {
		log.Printf("Handler.ValueJSONHandleFn: body Close error: %v\n", err)
	}
	log.Printf("Handler.ValueJSONHandleFn: Body = %v\n", string(body))
	if err = json.Unmarshal(body, &metric); err != nil {
		log.Printf("Handler.ValueJSONHandleFn: Unmarshal error: %v\n", err)
		return
	}
	metric, err = h.Server.Value(metric)
	if err != nil {
		log.Printf("Handler.ValueJSONHandleFn: Server Value Metric hasn`t been found.\n")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(metric)
	if err != nil {
		log.Printf("Handler.ValueJSONHandleFn: Marshal error: %v\n", err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Printf("Handler.ValueJSONHandleFn: Write error: %v\n", err)
	}
}
