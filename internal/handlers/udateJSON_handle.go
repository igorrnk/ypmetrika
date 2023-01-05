package handlers

import (
	"encoding/json"
	"github.com/igorrnk/ypmetrika/internal/models"
	"io"
	"log"
	"net/http"
)

func (h Handler) UpdateJSONHandleFn(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler.UpdateJSONHandleFn: URL = %v\n", r.URL)
	metric := models.Metric{}
	var body []byte
	var err error
	if body, err = io.ReadAll(r.Body); err != nil {
		log.Printf("Handler.UpdateJSONHandleFn: body ReadAll error: %v\n", err)
		return
	}
	if err = r.Body.Close(); err != nil {
		log.Printf("Handler.UpdateJSONHandleFn: body Close error: %v\n", err)
	}
	log.Printf("Handler.UpdateJSONHandleFn: Body = %v\n", string(body))
	if err = json.Unmarshal(body, &metric); err != nil {
		log.Printf("Handler.UpdateJSONHandleFn: Unmarshal error: %v\n", err)
		return
	}
	if metric, err = h.Server.UpdateValue(metric); err != nil {
		log.Printf("Handler.UpdateJSONHandleFn: Server UpdateValue Metric error: %v\n", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(metric)
	if err != nil {
		log.Printf("Handler.UpdateJSONHandleFn: Marshal error: %v\n", err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Printf("Handler.UpdateJSONHandleFn: Write error: %v\n", err)
	}
}
