package handlers

import (
	"log"
	"net/http"
)

func (h Handler) PingHandleFn(w http.ResponseWriter, r *http.Request) {
	err := h.Server.PingDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "db error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)

}
