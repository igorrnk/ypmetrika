package handlers

import (
	"log"
	"net/http"
)

func (h Handler) PingHandleFn(w http.ResponseWriter, r *http.Request) {
	err := h.Server.PingDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "The database isn't connected.", http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("The database is connected."))
	log.Printf("Request %v has been handled.", r.RequestURI)
}
