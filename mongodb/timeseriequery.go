package mongodb

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func (sp *SessionProvider) GetTimeSeriesData(dbname string, collection string, col string, timecol string, from time.Time, to time.Time, intervalMs int) ([][]float64, error) {
	var res [][]float64
	c := sp.Session.DB(dbname).C(collection)
	var results []bson.M
	pipeline := BuildTimeSeriesPipe(col, timecol, from, to, intervalMs)
	err := c.Pipe(pipeline).All(&results)
	if err != nil {
		return res, err
	}
	for _, v := range results {
		array := make([]float64, 2)
		date, err := parseIdtoDate(v, intervalMs)
		if err == ERRNilDataPoint {
			log.Println("Contain invalid data")
			continue
		} else if err != nil {
			return res, err
		}
		array[0] = convertFloat(v["value"])
		array[1] = convertFloat(date)
		res = append(res, array)
	}
	return res, nil
}


func BuildTimeSeriesPipe(col string, timecol string, from time.Time, to time.Time, intervalMs int) ([]bson.M) {
	trange := bson.M{timecol: bson.M{"$gte": from, "$lte": to}}
	pipeline := []bson.M{
		{ "$match": trange },
		{
			"$group": bson.M{
				"_id": buildTimeBson(timecol, intervalMs) ,
				"value": bson.M{"$avg": "$" + col},
			},
		},
		{
			"$sort": bson.M{
				"_id": 1,
			},
		},
	}
	return pipeline
}
//{ "$sort": bson.M{ "_id": 1}},
func buildTimeBson(timecol string, intervalMs int) bson.M {
	var ret bson.M
	if 86400000 <= intervalMs && intervalMs < 2629800000 {
		uni := "$day"
		return bson.M{
			"year": bson.M{
				"$year": "$" + timecol,
			},
			"month": bson.M{
				"$month": "$" + timecol,
			},
			"day": buildInterval(timecol, intervalMs, uni, 86400000),
		}
	} else if 3600000 <= intervalMs && intervalMs< 86400000 {
		uni := "$hour"
		return bson.M{
			"year": bson.M{
				"$year": "$" + timecol,
			},
			"month": bson.M{
				"$month": "$" + timecol,
			},
			"day": bson.M{
				"$dayOfMonth": "$" + timecol,
			},
			"hour": buildInterval(timecol, intervalMs, uni, 3600000),
		}
	} else if 60000 <= intervalMs && intervalMs< 3600000 {
		uni := "$minute"
		return bson.M{
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
			"minute": buildInterval(timecol, intervalMs, uni, 60000),
		}
	} else {
		uni := "$second"
		return bson.M{
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
			"interval": buildInterval(timecol, intervalMs, uni, 1000),
		}
	}
	return ret
}

func buildInterval( timecol string, intervalMs int, uni string, ms int) bson.M {
	type list []interface{}
	interval := intervalMs/ms
	if interval < 1{
		interval = 1
	}
	mod := list{ bson.M{uni: "$" + timecol}, interval}
	sub := list{ bson.M{uni: "$" + timecol}, bson.M{"$mod": mod } }
	return bson.M{"$subtract": sub}
}

func parseIdtoDate(v bson.M, intervalMs int) (time.Time, error){
	var year, month, day, hour, minute, second int
	if v["_id"].(bson.M)["year"] == nil {
		log.Println("1")
		return time.Time{}, ERRNilDataPoint
	}
	year = v["_id"].(bson.M)["year"].(int)
	if v["_id"].(bson.M)["month"] == nil {
		log.Println("2")
		return time.Time{}, ERRNilDataPoint
	}
	month = v["_id"].(bson.M)["month"].(int)
	if intervalMs < 2629800000 {
		if v["_id"].(bson.M)["day"] == nil {
			log.Println("3")
			return time.Time{}, ERRNilDataPoint
		}
		if intervalMs >= 86400000 {
			goto fin
		}
		day = v["_id"].(bson.M)["day"].(int)
	}
	if intervalMs < 86400000 {
		if v["_id"].(bson.M)["hour"] == nil {
			log.Println("4")
			return time.Time{}, ERRNilDataPoint
		}
		if intervalMs >= 3600000{
			goto fin
		}
		hour = v["_id"].(bson.M)["hour"].(int)
	}
	if intervalMs < 3600000 {
		if v["_id"].(bson.M)["minute"] == nil {
			log.Println("5")
			return time.Time{}, ERRNilDataPoint
		}
		if intervalMs >= 60000 {
			goto fin
		}
		minute = v["_id"].(bson.M)["minute"].(int)
	}
	if intervalMs < 60000 {
		if v["_id"].(bson.M)["interval"] == nil {
			log.Println("6")
			return time.Time{}, ERRNilDataPoint
		}
		second = v["_id"].(bson.M)["interval"].(int)
	}
	fin:
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC), nil
}

