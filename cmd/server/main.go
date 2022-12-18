package main

import (
	"fmt"
	"log"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println("POST")
		fmt.Println(r.RequestURI)
	}

}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", HelloWorld)
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
