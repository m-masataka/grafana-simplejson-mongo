package mongodb

import (
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

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
				"$dayOfMonth": "$" + timecol,
			},
			"second": bson.M{
				"$second": submod,
			},
		}
	}
	return ret
}
