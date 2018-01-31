package api

import (
	"encoding/json"
	"net/http"

	"github.com/grafana-simplejson-mongo/mongodb"
)

func checkRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func reqSearch(w http.ResponseWriter, r *http.Request) {
	/*fake*/
	var result SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var resp []string
	for i := 0; i < 5; i++ {
		resp = append(resp, result.Target+string(i))
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(bytes)
}

func reqQuery(w http.ResponseWriter, r *http.Request) {
	var result TimeSeriesQuery
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	sp := mongodb.NewSession()
	var resp []TimeSeriesResponse
	for i, v := range result.Targets {
		var err error
		resp = append(resp, TimeSeriesResponse{Target: v.Target})
		resp[i].DataPoint, err = sp.GetTimeSeriesData("fluentd", v.Target, v.Target, "time")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(bytes)
}

func addTarget(tqs []TimeSeriesResponse, name string) {
	target := TimeSeriesResponse{Target: name}
	tqs = append(tqs, target)
}

func (tq *TimeSeriesResponse) addDataPoint(x float64, y float64) {
	array := make([]float64, 2)
	array[0] = x
	array[1] = y
	tq.DataPoint = append(tq.DataPoint, array)
}
