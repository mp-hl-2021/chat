package main

import (
	"github.com/mp-hl-2021/chat/api"
	"net/http"
	"time"
)

func main() {
	router := api.NewRouter()

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

