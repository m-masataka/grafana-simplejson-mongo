package mongodb

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	host = "localhost"
	port = 27017
)

type SessionProvider struct {
	Session *mgo.Session
}

func NewSession() SessionProvider {
	var sp SessionProvider
	session, err := mgo.Dial("mongodb://" + host)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	sp.Session = session
	return sp
}

func (sp *SessionProvider) GetTimeSeriesData(dbname string, collection string, col string, time string) ([][]float64, error) {
	var res [][]float64
	c := sp.Session.DB(dbname).C(collection)
	var results []bson.M
	sel := bson.M{"_id": 0, col: 1, time: 1}
	err := c.Find(nil).Select(sel).Sort("-timestamp").All(&results)
	if err != nil {
		return res, err
	}
	for _, v := range results {
		array := make([]float64, 2)
		array[0] = convertFloat(v[col])
		array[1] = convertFloat(v[time])
		res = append(res, array)
	}
	return res, nil
}

func convertFloat(v interface{}) float64 {
	var r float64
	switch v.(type) {
	case int:
		r = float64(v.(int))
	case float64:
		r = v.(float64)
	case time.Time:
		r = float64(v.(time.Time).UnixNano() / int64(time.Millisecond))
	}
	return r
}
