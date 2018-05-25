package api

import (
	"regexp"
	"strings"
)

var TimeSeriesColumn = regexp.MustCompile(`{.+?}`)

func TimeSeriesColumnRegexp(str string) []string {
	var ret []string
	data := TimeSeriesColumn.FindString(str)
	data = strings.TrimRight(data, "}")
	data = strings.TrimLeft(data, "{")
	ret = strings.Split(data, ",")
	return ret
}
