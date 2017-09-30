package model

import "github.com/ptrv/go-gpx"

type Feature interface {
	CsvLine() string
}

type Model interface {
	Features([]gpx.Wpt) []Feature
	Cut([]gpx.Wpt, []Feature) []gpx.Wpt
}
