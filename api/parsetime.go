package api

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ToNow    = regexp.MustCompile(`now-[0-9]+[mhdMyw]$`)
	PerNow   = regexp.MustCompile(`now/[dMyw]$`)
	Now      = regexp.MustCompile(`now$`)
	PerToNow = regexp.MustCompile(`now-[0-9]+[dMyw]/[dMyw]`)
)

var (
	ERRRangeFromat = errors.New("Range Fromat Error")
)

func boolRegexp(str string, re *regexp.Regexp) bool {
	return re.MatchString(str)
}

func parseToNow(from string, to string) (time.Time, time.Time, error) {
	var f, t time.Time
	if to != "now" {
		return f, t, ERRRangeFromat
	}
	t = time.Now()
	var subD time.Duration
	trim := from[4:]
	if strings.Contains(trim, "m") {
		subD = time.Second * 60
	} else if strings.Contains(trim, "h") {
		subD = time.Second * 60 * 60
	} else if strings.Contains(trim, "d") {
		subD = time.Second * 60 * 60 * 24
	} else if strings.Contains(trim, "M") {
		//ToDo: Calcurate exact days
		subD = time.Second * 60 * 60 * 24 * 30
	} else if strings.Contains(trim, "y") {
		subD = time.Second * 60 * 60 * 24 * 365
	} else if strings.Contains(trim, "w") {
		subD = time.Second * 60 * 60 * 24 * 7
	} else {
		return f, t, ERRRangeFromat
	}
	num, err := strconv.Atoi(trim[:len(trim)-1])
	if err != nil {
		return f, t, ERRRangeFromat
	}
	f = time.Now().Add(-time.Duration(num) * subD)

	return f, t, nil
}

func parsePerNow(from string, to string) (time.Time, time.Time, error) {
	var f, t time.Time
	now := time.Now()
	trim := from[4:]
	if trim == "d" {
		f = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	} else if trim == "M" {
		f = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else if trim == "y" {
		f = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	} else if trim == "w" {
		_, thisWeek := now.ISOWeek()
		beDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		f = beDay.AddDate(0, 0, -thisWeek)
	} else {
		return f, t, ERRRangeFromat
	}

	if to == "now" {
		t = now
		return f, t, nil
	}

	trim = to[4:]
	if trim == "d" {
		beDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		t = beDay.AddDate(0, 0, 1)
	} else if trim == "M" {
		beMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		t = beMonth.AddDate(0, 1, 0)
	} else if trim == "y" {
		beYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		t = beYear.AddDate(1, 0, 0)
	} else if trim == "w" {
		_, thisWeek := now.ISOWeek()
		beDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		t = beDay.AddDate(0, 0, 7-thisWeek)
	} else {
		return f, t, ERRRangeFromat
	}
	return f, t, nil
}

func parsePerToNow(from string, to string) (time.Time, time.Time, error) {
	var f, t time.Time
	now := time.Now()
	if from != to {
		return f, t, ERRRangeFromat
	}

	trim := from[4:]
	num, err := strconv.Atoi(trim[:len(trim)-3])
	if err != nil {
		return f, t, ERRRangeFromat
	}

	if strings.Contains(trim, "d") {
		beDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		t = beDay.AddDate(0, 0, -(num - 1))
		f = t.AddDate(0, 0, -1)
	} else if strings.Contains(trim, "w") {
		_, thisWeek := now.ISOWeek()
		beDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		f = beDay.AddDate(0, 0, -(thisWeek + (7 * num)))
		t = beDay.AddDate(0, 0, -(thisWeek + (7 * (num - 1))))
	} else if strings.Contains(trim, "M") {
		beMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		t = beMonth.AddDate(0, -(num - 1), 0)
		f = t.AddDate(0, -1, 0)
	} else if strings.Contains(trim, "y") {
		beYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		t = beYear.AddDate(-(num - 1), 0, 0)
		f = t.AddDate(-1, 0, 0)
	} else {
		return f, t, ERRRangeFromat
	}
	return f, t, nil
}

func parseISODate(from string, to string) (time.Time, time.Time, error) {
	var f, t time.Time
	layout := "2006-01-02T15:04:05.000Z"
	f, err := time.Parse(layout, from)
	if err != nil {
		return f, t, err
	}
	t, err = time.Parse(layout, to)
	if err != nil {
		return f, t, err
	}
	return f, t, nil
}
