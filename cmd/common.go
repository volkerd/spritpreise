package cmd

import (
	"fmt"
	"os"
	"time"
)

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
	STATION_PATH = "stations"
	PRICE_PATH   = "prices"
	CUT_OFF_DATE = "2019-01-23"
	STATION      = "station"
	PRICE        = "price"
)

var (
	zipcodeFilter  map[string]bool
	basePath       string
	influxdbBucket string
	influxdbOrg    string
	influxdbToken  string
	influxdbUrl    string
)

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
	basePath = os.Getenv("SPRITPREISE_BASE_PATH")
	influxdbBucket = os.Getenv("SPRITPREISE_INFLUXDB_BUCKET")
	influxdbOrg = os.Getenv("SPRITPREISE_INFLUXDB_ORG")
	influxdbToken = os.Getenv("SPRITPREISE_INFLUXDB_TOKEN")
	influxdbUrl = os.Getenv("SPRITPREISE_INFLUXDB_URL")

	zipcodeFilter = map[string]bool{
		"71364": true, // Winnenden
		"71522": true, // Backnang
		"73614": true, // Schorndorf
		"73650": true, // Winterbach
		"73663": true, // Berglen
	}

	fmt.Printf("copy data from %s to bucket %s\n", basePath, influxdbBucket)
}
