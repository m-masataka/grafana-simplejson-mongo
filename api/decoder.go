package api

type TimeSeriesQuery struct {
	PanelID int `json:"panelId"`
	Range   struct {
		From string `json:"from"`
		To   string `json:"to"`
		Raw  struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"raw"`
	} `json:"range"`
	RangeRaw struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"rangeRaw"`
	Interval   string `json:"interval"`
	IntervalMs int    `json:"intervalMs"`
	Targets    []struct {
		Target string `json:"target"`
		RefID  string `json:"refId"`
		Type   string `json:"type"`
	} `json:"targets"`
	Format        string `json:"format"`
	MaxDataPoints int    `json:"maxDataPoints"`
}

type TimeSeriesResponse struct {
	Target    string      `json:"target"`
	DataPoint [][]float64 `json:"datapoints"`
}

type SearchRequest struct {
	Target string `json:"target"`
}
