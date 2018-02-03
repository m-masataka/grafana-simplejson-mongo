package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func (conf *Config) reqSearch(w http.ResponseWriter, r *http.Request) {
	log.Println("Search Query")
	/*fake*/
	var result SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var resp []string
	for i := 0; i < 5; i++ {
		resp = append(resp, "OK"+string(i))
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(bytes)
}
