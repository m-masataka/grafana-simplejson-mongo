package main

import (
	"log"

	"github.com/m-masataka/grafana-simplejson-mongo/api"
)

func main() {
	conf := api.Config{
		Port:      8080,
		MongoHost: "localhost",
	}
	errs := make(chan error, 2)
	api.StartHTTPServer(conf, errs)
	log.Println("start")
	for {
		err := <-errs
		log.Println(err)
	}
}
