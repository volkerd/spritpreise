package cmd

import "time"

type File struct {
	path     string
	date     string
	datatype string
}

type FilePair struct {
	station File
	price   File
}

type StationFilter struct {
	zipcodes []string
}

type Station struct {
	uuid         string
	name         string
	brand        string
	street       string
	house_number string
	post_code    string
	city         string
	latitude     string
	longitude    string
}

type Price struct {
	date         time.Time
	station_uuid string
	diesel       float32
	e5           float32
	e10          float32
	dieselchange bool
	e5change     bool
	e10change    bool
}

const (
	BASE_PATH       = "/Users/volkerdemel/work/tanken/tankerkoenig-data"
	STATION_PATH    = "stations"
	PRICE_PATH      = "prices"
	CUT_OFF_DATE    = "2019-01-23"
	STATION         = "station"
	PRICE           = "price"
	INFLUXDB_BUCKET = "spritpreise_test"
	INFLUXDB_ORG    = "demelnet"
	INFLUXDB_TOKEN  = "FUzVuYQyM1OLGndyN9mxwIRmphWu53pgOhDCFbB_f9rJ7IL9RKI3mv9ftFBaPnIzAm5Tpjk5vUsKqo0-fLhlzg=="
	// Store the URL of your InfluxDB instance
	INFLUXDB_URL = "http://nass:8086"
)

var zipcodeFilter map[string]bool

func First[E any](s []E) E {
	if len(s) == 0 {
		var zero E
		return zero
	}
	return s[0]
}

func Last[E any](s []E) E {
	if len(s) == 0 {
		var zero E
		return zero
	}
	return s[len(s)-1]
}

func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maximum(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func init() {
	zipcodeFilter = map[string]bool{
		"71364": true, // Winnenden
		"71522": true, // Backnang
		"73614": true, // Schorndorf
		"73650": true, // Winterbach
		"73663": true, // Berglen
	}
}
