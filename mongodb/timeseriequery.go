package mongodb

import (
	"log"
	"time"
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

func (sp *SessionProvider) GetTimeSeriesData(dbname string, collection string, col string, timecol string, from time.Time, to time.Time, intervalMs int) ([][]float64, error) {
	var res [][]float64
	var strflag bool
	c := sp.Session.DB(dbname).C(collection)
	var results []bson.M
	var judge bson.M
	err := c.Find(nil).One(&judge)
	if err != nil {
		return res, err
	}
	if reflect.TypeOf(judge[timecol]).Kind()  == reflect.String {
		strflag = true
	}else {
		strflag = false
	}
	pipeline := BuildTimeSeriesPipe(col, timecol, from, to, intervalMs, strflag)
	err = c.Pipe(pipeline).All(&results)
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


func BuildTimeSeriesPipe(col string, timecol string, from time.Time, to time.Time, intervalMs int, strflag bool) ([]bson.M) {
	var trange bson.M
	if strflag {
		trange = bson.M{ timecol: bson.M{"$gte": from.Format("20060102150405"), "$lte": to.Format("20060102150405")}}
	} else {
		trange = bson.M{ timecol: bson.M{"$gte": from, "$lte": to}}
	}
	pipeline := []bson.M{
		{ "$match": trange },
		{
			"$group": bson.M{
				"_id": buildTimeBson(timecol, intervalMs, strflag) ,
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
func buildTimeBson(timecol string, intervalMs int, strflag bool) bson.M {
	var ret bson.M
	var TimeCol interface{}
	if strflag {
		TimeCol = bson.M{ "$dateFromString": bson.M{ "dateString": "$" + timecol}}
	} else {
		TimeCol = "$" + timecol
	}
	if 86400000 <= intervalMs && intervalMs < 2629800000 {
		uni := "$day"
		return bson.M{
			"year": bson.M{
				"$year": TimeCol,
			},
			"month": bson.M{
				"$month": TimeCol,
			},
			"day": buildInterval(timecol, intervalMs, uni, 86400000, strflag),
		}
	} else if 3600000 <= intervalMs && intervalMs< 86400000 {
		uni := "$hour"
		return bson.M{
			"year": bson.M{
				"$year": TimeCol,
			},
			"month": bson.M{
				"$month": TimeCol,
			},
			"day": bson.M{
				"$dayOfMonth": TimeCol,
			},
			"hour": buildInterval(timecol, intervalMs, uni, 3600000, strflag),
		}
	} else if 60000 <= intervalMs && intervalMs< 3600000 {
		uni := "$minute"
		return bson.M{
			"year": bson.M{
				"$year": TimeCol,
			},
			"month": bson.M{
				"$month": TimeCol,
			},
			"day": bson.M{
				"$dayOfMonth": TimeCol,
			},
			"hour": bson.M{
				"$hour": TimeCol,
			},
			"minute": buildInterval(timecol, intervalMs, uni, 60000, strflag),
		}
	} else {
		uni := "$second"
		return bson.M{
			"year": bson.M{
				"$year": TimeCol,
			},
			"month": bson.M{
				"$month": TimeCol,
			},
			"day": bson.M{
				"$dayOfMonth": TimeCol,
			},
			"hour": bson.M{
				"$hour": TimeCol,
			},
			"minute": bson.M{
				"$minute": TimeCol,
			},
			"interval": buildInterval(timecol, intervalMs, uni, 1000, strflag),
		}
	}
	return ret
}

func buildInterval( timecol string, intervalMs int, uni string, ms int, strflag bool) bson.M {
	type list []interface{}
	interval := intervalMs/ms
	if interval < 1{
		interval = 1
	}
	var TimeCol bson.M
	if strflag {
		TimeCol = bson.M{ "$dateFromString": bson.M{ "dateString": "$" + timecol}}
	} else {
		TimeCol = bson.M{ "$dateFromString": bson.M{ "dateString": "$" + timecol}}
	}
	mod := list{ bson.M{uni: TimeCol}, interval}
	sub := list{ bson.M{uni: TimeCol}, bson.M{"$mod": mod } }
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

