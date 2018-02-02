package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/m-masataka/grafana-simplejson-mongo/mongodb"
)

type TSQuery struct {
	DB         string
	Collection string
	Col        string
	TimeCol    string
	From       time.Time
	To         time.Time
	Interval   string
}

func checkRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func reqSearch(w http.ResponseWriter, r *http.Request) {
	log.Println("Search Query")
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

func (conf *Config) reqQuery(w http.ResponseWriter, r *http.Request) {
	log.Println("Time Series Query")
	var result TimeSeriesQuery
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	sp := mongodb.NewSession(conf.MongoHost)
	var q TSQuery
	err := q.parseRangeRaw(result.RangeRaw.From, result.RangeRaw.To)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q.Interval = result.Interval
	var resp []TimeSeriesResponse
	for i, v := range result.Targets {
		resp = append(resp, TimeSeriesResponse{Target: v.Target})
		err := q.parseTarget(v.Target)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp[i].DataPoint, err = sp.GetTimeSeriesData(q.DB, q.Collection, q.Col, q.TimeCol, q.From, q.To, q.Interval)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func (q *TSQuery) parseTarget(target string) error {
	res := strings.Split(target, ".")
	if len(res) < 4 {
		return ERRFormat
	}
	q.DB = res[0]
	q.Collection = res[1]
	q.Col = res[2]
	q.TimeCol = res[3]
	return nil
}

var (
	ERRTimePerDay = errors.New("TimePerDay")
	ERRFormat     = errors.New("Time format does not match")
	ERRSameTo     = errors.New("To is same as From")
)

func (q *TSQuery) parseRangeRaw(from string, to string) error {
	t, t2, err := parseTime(from)
	if err == ERRTimePerDay {
		now := time.Now()
		if strings.Contains(from, "d") {
			q.From = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			q.To = q.From.AddDate(0, 0, 1)
		} else if strings.Contains(from, "M") {
			q.From = time.Date(now.Year(), now.Month(), 1, 1, 0, 0, 0, time.UTC)
			q.To = q.From.AddDate(0, 1, 0)
		} else if strings.Contains(from, "y") {
			q.From = time.Date(now.Year(), 1, 1, 1, 0, 0, 0, time.UTC)
			q.To = q.From.AddDate(1, 0, 0)
		} else if strings.Contains(from, "w") {
			_, thisWeek := now.ISOWeek()
			thisDay := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, time.UTC)
			q.From = thisDay.AddDate(0, 0, -thisWeek)
			q.To = thisDay.AddDate(0, 0, 7-thisWeek)
		} else {
			return ERRFormat
		}
		return nil
	} else if err == ERRSameTo {
		q.From = t
		q.To = t2
		return nil
	} else if err != nil {
		return err
	}
	q.From = t

	t, _, err = parseTime(to)
	if err != nil {
		return err
	}
	q.To = t
	return nil
}

func parseTime(str string) (time.Time, time.Time, error) {
	if str == "now" {
		return time.Now(), time.Time{}, nil
	}

	if strings.Contains(str, "Z") {
		layout := "2006-01-02T15:04:05.000Z"
		t, err := time.Parse(layout, str)
		return t, time.Time{}, err
	}

	if strings.Contains(str, "-") {
		v := strings.Split(str, "-")
		now := time.Now()
		var d time.Duration
		var err error
		if strings.Contains(v[1], "d") {
			var trim string
			var subtime time.Duration
			if strings.Contains(v[1], "/") {
				trim = strings.TrimRight(v[1], "/d")
				beginningOfTheDay := time.Date(now.Year(), now.Month(), now.Day(), 1, 1, 1, 1, time.UTC)
				subtime = time.Now().Sub(beginningOfTheDay)
			} else {
				trim = v[1]
			}
			num, err := strconv.Atoi(strings.TrimRight(trim, "d"))
			if err != nil {
				return time.Time{}, time.Time{}, err
			}
			d = (time.Duration(num) * time.Hour * 24) - subtime
			return now.Add(-d), now.Add(-subtime), ERRSameTo
		} else if strings.Contains(v[1], "M") {
			var trim string
			var subtime time.Duration
			if strings.Contains(v[1], "/") {
				trim = strings.TrimRight(v[1], "/M")
				beginningOfTheMonth := time.Date(now.Year(), now.Month(), 1, 1, 1, 1, 1, time.UTC)
				subtime = time.Now().Sub(beginningOfTheMonth)
			} else {
				trim = v[1]
			}
			num, err := strconv.Atoi(strings.TrimRight(trim, "M"))
			if err != nil {
				return time.Time{}, time.Time{}, err
			}
			d = (time.Duration(num) * time.Hour * 24 * 30) - subtime
			return now.Add(-d), now.Add(-subtime), ERRSameTo
		} else if strings.Contains(v[1], "y") {
			var trim string
			var subtime time.Duration
			if strings.Contains(v[1], "/") {
				trim = strings.TrimRight(v[1], "/y")
				beginningOfTheYear := time.Date(now.Year(), 1, 1, 1, 1, 1, 1, time.UTC)
				subtime = time.Now().Sub(beginningOfTheYear)
			} else {
				trim = v[1]
			}
			num, err := strconv.Atoi(strings.TrimRight(trim, "y"))
			if err != nil {
				return time.Time{}, time.Time{}, err
			}
			d = (time.Duration(num) * time.Hour * 24 * 365) - subtime
			return now.Add(-d), now.Add(-subtime), ERRSameTo
		} else {
			d, err = time.ParseDuration(v[1])
			if err != nil {
				return time.Time{}, time.Time{}, err
			}
			return now.Add(-d), time.Time{}, nil
		}
	}

	if strings.Contains(str, "/") {
		return time.Time{}, time.Time{}, ERRTimePerDay
	}

	return time.Time{}, time.Time{}, ERRFormat
}
