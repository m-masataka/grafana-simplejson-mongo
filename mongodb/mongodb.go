package mongodb

import (
	"errors"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SessionProvider struct {
	Session *mgo.Session
}

func NewSession(host string) SessionProvider {
	var sp SessionProvider
	session, err := mgo.Dial("mongodb://" + host)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	sp.Session = session
	return sp
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

var (
	ERRNilDataPoint = errors.New("NilDataPoint")
)

func parseDate(v bson.M, num int) (time.Time, error) {
	var year, month, day, hour, minute, second, milisec int
	log.Println(v)
	for i := 0; i <= num; i++ {
		switch i {
		case 0:
			if v["_id"].(bson.M)["year"] == nil {
				return time.Time{}, ERRNilDataPoint
			}
			year = v["_id"].(bson.M)["year"].(int)
		case 1:
			if v["_id"].(bson.M)["month"] == nil {
				return time.Time{}, ERRNilDataPoint
			}
			month = v["_id"].(bson.M)["month"].(int)
		case 2:
			if v["_id"].(bson.M)["day"] == nil {
				return time.Time{}, ERRNilDataPoint
			}
			day = v["_id"].(bson.M)["day"].(int)
		case 3:
			if v["_id"].(bson.M)["hour"] == nil {
				return time.Time{}, ERRNilDataPoint
			}
			hour = v["_id"].(bson.M)["hour"].(int)
		case 4:
			if v["_id"].(bson.M)["minute"] == nil {
				return time.Time{}, ERRNilDataPoint
			}
			minute = v["_id"].(bson.M)["minute"].(int)
		case 5:
			if v["_id"].(bson.M)["second"] == nil {
				return time.Time{}, ERRNilDataPoint
			}
			second = v["_id"].(bson.M)["second"].(int)
		case 6:
			if v["_id"].(bson.M)["milisecond"] == nil {
				return time.Time{}, ERRNilDataPoint
			}
			milisec = v["_id"].(bson.M)["milisecond"].(int)
		}
	}
	return time.Date(year, time.Month(month), day, hour, minute, second, milisec, time.UTC), nil
}
