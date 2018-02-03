package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Port      int
	MongoHost string
}

func httpServer(conf Config) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", checkRequest)
	r.HandleFunc("/search", conf.reqSearch)
	r.HandleFunc("/query", conf.reqQuery)
	return r
}

func StartHTTPServer(conf Config, errChan chan error) {
	go func() {
		p := fmt.Sprintf(":%d", conf.Port)
		errChan <- http.ListenAndServe(p, httpServer(conf))
	}()
}
