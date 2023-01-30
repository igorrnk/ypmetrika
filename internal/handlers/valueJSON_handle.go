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
	metric := &models.Metric{}
	var body []byte
	var err error
	if body, err = io.ReadAll(r.Body); err != nil {
		log.Printf("Handler.ValueJSONHandleFn: body ReadAll error: %v\n", err)
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	log.Printf("Handler.ValueJSONHandleFn: Body = %v\n", string(body))
	if err = json.Unmarshal(body, &metric); err != nil {
		log.Printf("Handler.ValueJSONHandleFn: Unmarshal error: %v\n", err)
		http.Error(w, "Unable to decode body", http.StatusBadRequest)
		return
	}
	metric, err = h.Server.Value(metric)
	if err != nil {
		log.Printf("Handler.ValueJSONHandleFn: Server Value Metric hasn`t been found.\n")
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&metric); err != nil {
		http.Error(w, "unable to serialize metric", http.StatusInternalServerError)
		return
	}
}
