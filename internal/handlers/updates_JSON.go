package handlers

import (
	"encoding/json"
	"errors"
	"github.com/igorrnk/ypmetrika/internal/models"
	"io"
	"log"
	"net/http"
)

func (h Handler) UpdatesJSONFn(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler.UpdateJSONHandleFn: URL = %v\n", r.URL)
	metrics := make([]*models.Metric, 0)
	var body []byte
	var err error
	if body, err = io.ReadAll(r.Body); err != nil {
		log.Printf("Handler.UpdateJSONHandleFn: body ReadAll error: %v\n", err)
		return
	}
	log.Printf("Handler.UpdateJSONHandleFn: Body = %v\n", string(body))
	if err = json.Unmarshal(body, &metrics); err != nil {
		log.Printf("Handler.UpdateJSONHandleFn: Unmarshal error: %v\n", err)
		return
	}
	err = h.Server.Updates(metrics)
	if errors.Is(err, models.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
	if err != nil {
		log.Printf("Handler.UpdateJSONHandleFn: Server UpdateValue Metric error: %v\n", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("Handler.UpdateJSONHandleFn: Marshal error: %v\n", err)
		return
	}
	w.Write(data)
}
