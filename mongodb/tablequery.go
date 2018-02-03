package mongodb

import (
	"log"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func (sp *SessionProvider) GetTableData(dbname string, collection string, timecol string, from time.Time, to time.Time) ([][]string, [][]string, error) {
	var keys [][]string
	var rows [][]string
	c := sp.Session.DB(dbname).C(collection)
	var results []bson.M

	var find bson.M
	if timecol != "" {
		find = bson.M{
			timecol: bson.M{
				"$gt": from,
				"$lt": to,
			},
		}
	} else {
		find = nil
	}

	err := c.Find(find).All(&results)
	if err != nil {
		log.Println(err)
		return keys, rows, err
	}
	if len(results) < 1 {
		return keys, rows, nil
	}
	for k, v := range results[0] {
		var key []string
		key = append(key, k)
		key = append(key, defineType(v))
		keys = append(keys, key)
	}
	for i := 0; i < len(results); i++ {
		var row []string
		for _, key := range keys {
			row = append(row, convertString(results[i][key[0]]))
		}
		rows = append(rows, row)
	}
	return keys, rows, nil
}

func defineType(v interface{}) string {
	var ret string
	switch v.(type) {
	case int:
		ret = "number"
	case time.Time:
		ret = "time"
	case bson.ObjectId:
		ret = "string"
	case string:
		ret = "string"
	}
	return ret
}

func convertString(v interface{}) string {
	var ret string
	switch v.(type) {
	case int:
		ret = strconv.Itoa(v.(int))
	case time.Time:
		ret = v.(time.Time).String()
	case bson.ObjectId:
		ret = v.(bson.ObjectId).String()
	case string:
		ret = v.(string)
	}
	return ret
}
