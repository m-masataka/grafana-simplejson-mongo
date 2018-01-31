package mongodb

type MongodbQuery struct {
	Start   int64                    `json:"start"`
	End     int64                    `json:"end"`
	Queries []map[string]interface{} `json:"query"`
}
