package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Port int
	Host string
}

func httpServer() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", checkRequest)
	r.HandleFunc("/search", reqSearch)
	r.HandleFunc("/query", reqQuery)
	return r
}

func StartHTTPServer(config Config, errChan chan error) {
	go func() {
		p := fmt.Sprintf(":%d", config.Port)
		errChan <- http.ListenAndServe(p, httpServer())
	}()
}
