package api

import (
	"regexp"
	"strings"
)

var TimeSeriesColumn = regexp.MustCompile(`data{.+?}`)

func TimeSeriesColumnRegexp(str string) []string {
	var ret []string
	data := TimeSeriesColumn.FindString(str)
	data = strings.TrimRight(data, "}")
	data = strings.TrimLeft(data, "data{")
	ret = strings.Split(data, ",")
	return ret
}
