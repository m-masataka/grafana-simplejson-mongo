package mongodb

import (
	"errors"
	"log"
	"strconv"
	"strings"
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

func (sp *SessionProvider) GetTimeSeriesData(dbname string, collection string, col string, time string, from time.Time, to time.Time, interval string) ([][]float64, error) {
	var res [][]float64
	c := sp.Session.DB(dbname).C(collection)
	var results []bson.M
	pipeline, err := BuildTimeSeriesPipe(col, time, from, to, interval)
	err = c.Pipe(pipeline).All(&results)
	if err != nil {
		return res, err
	}
	for _, v := range results {
		array := make([]float64, 2)
		num := 0
		if strings.Contains(interval, "h") {
			num = 3
		} else if strings.Contains(interval, "m") {
			num = 4
		} else if strings.Contains(interval, "s") {
			num = 5
		}
		date, err := parseDate(v, num)
		if err == ERRNilDataPoint {
			log.Println("Contain invalid data")
			continue
		} else if err != nil {
			return res, err
		}
		array[0] = convertFloat(v[col])
		array[1] = convertFloat(date)
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

var (
	ERRNilDataPoint = errors.New("NilDataPoint")
)

func parseDate(v bson.M, num int) (time.Time, error) {
	var year, month, day, hour, minute, second, milisec int
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

func BuildTimeSeriesPipe(col string, timecol string, from time.Time, to time.Time, interval string) ([]bson.M, error) {
	var uni string
	var num int
	var err error
	if strings.Contains(interval, "s") {
		uni = "$second"
		num, err = strconv.Atoi(strings.TrimRight(interval, "s"))
		if err != nil {
			return []bson.M{}, err
		}
	} else if strings.Contains(interval, "m") {
		uni = "$minute"
		num, err = strconv.Atoi(strings.TrimRight(interval, "m"))
		if err != nil {
			return []bson.M{}, err
		}
	} else if strings.Contains(interval, "h") {
		uni = "$hour"
		num, err = strconv.Atoi(strings.TrimRight(interval, "h"))
		if err != nil {
			return []bson.M{}, err
		}
	}
	trange := bson.M{timecol: bson.M{"$gte": from, "$lte": to}}
	timeBson := buildTimeBson(uni, timecol, num)
	pipeline := []bson.M{
		{
			"$match": trange,
		},
		{
			"$group": bson.M{
				"_id": timeBson,
				col: bson.M{
					"$avg": "$" + col,
				},
			},
		},
		{
			"$sort": bson.M{
				"_id": 1,
			},
		},
	}
	return pipeline, nil
}

func buildTimeBson(uni string, timecol string, interval int) bson.M {
	mod := []interface{}{bson.M{uni: "$" + timecol}, interval}
	submod := []interface{}{bson.M{uni: "$" + timecol}, bson.M{"$mod": mod}}
	var ret bson.M
	switch uni {
	case "$hour":
		ret = bson.M{
			"year": bson.M{
				"$year": "$" + timecol,
			},
			"month": bson.M{
				"$month": "$" + timecol,
			},
			"day": bson.M{
				"$dayOfMonth": "$" + timecol,
			},
			"hour": bson.M{
				"$subtract": submod,
			},
		}
	case "$minute":
		ret = bson.M{
			"year": bson.M{
				"$year": "$" + timecol,
			},
			"month": bson.M{
				"$month": "$" + timecol,
			},
			"day": bson.M{
				"$dayOfMonth": "$" + timecol,
			},
			"hour": bson.M{
				"$hour": "$" + timecol,
			},
			"minute": bson.M{
				"$subtract": submod,
			},
		}
	case "$second":
		ret = bson.M{
			"year": bson.M{
				"$year": "$" + timecol,
			},
			"month": bson.M{
				"$month": "$" + timecol,
			},
			"day": bson.M{
				"$dayOfMonth": "$" + timecol,
			},
			"hour": bson.M{
				"$hour": "$" + timecol,
			},
			"minute": bson.M{
				"$minute": "$" + timecol,
			},
			"second": bson.M{
				"$subtract": submod,
			},
		}
	case "$milisecond":
		ret = bson.M{
			"year": bson.M{
				"$year": "$" + timecol,
			},
			"month": bson.M{
				"$month": "$" + timecol,
			},
			"day": bson.M{
				"$dayOfMonth": "$" + timecol,
			},
			"hour": bson.M{
				"$hour": "$" + timecol,
			},
			"minute": bson.M{
				"$minute": "$" + timecol,
			},
			"second": bson.M{
				"$second": "$" + timecol,
			},
			"millisecond": bson.M{
				"$subtract": submod,
			},
		}

	}
	return ret
}
