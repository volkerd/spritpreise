package cmd

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func readStations(path string) (stations []Station) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("error opening file %s: %v", path, err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("error closing file %s: %v", path, err)
		}
	}(f)
	r := csv.NewReader(f)
	_, err = r.Read()
	if err != nil {
		log.Fatalf("error reading file %s: %v", path, err)
	}
	//fmt.Println(header)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error reading file %s: %v", path, err)
		}
		station := Station{
			uuid:         record[0],
			name:         record[1],
			brand:        record[2],
			street:       record[3],
			house_number: record[4],
			post_code:    record[5],
			city:         record[6],
			latitude:     record[7],
			longitude:    record[8],
		}
		_, isPresent := zipcodeFilter[station.post_code]
		if isPresent {
			//fmt.Println(record)
			stations = append(stations, station)
		}
	}
	//fmt.Printf("%d stations found in %s\n", len(stations), path)
	return stations
}

func createStationFilter(stations []Station) (stationFilter map[string]Station) {
	stationFilter = make(map[string]Station)
	for _, station := range stations {
		stationFilter[station.uuid] = station
	}
	return stationFilter
}

func readPrices(path string, stationFilter map[string]Station) (prices []Price) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("error opening file %s: %v", path, err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("error closing file %s: %v", path, err)
		}
	}(f)
	r := csv.NewReader(f)
	_, err = r.Read()
	if err != nil {
		log.Fatalf("error reading header from file %s: %v", path, err)
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error reading file %s: %v", path, err)
		}
		price := Price{
			date:         parseTime(record[0]),
			station_uuid: record[1],
			diesel:       parseFloat(record[2]),
			e5:           parseFloat(record[3]),
			e10:          parseFloat(record[4]),
			dieselchange: parseBool(record[5]),
			e5change:     parseBool(record[6]),
			e10change:    parseBool(record[7]),
		}
		_, isPresent := stationFilter[price.station_uuid]
		if isPresent {
			//fmt.Println(record)
			prices = append(prices, price)
		}
	}
	//fmt.Printf("%d prices found in %s\n", len(prices), path)
	return prices
}

func parseTime(aString string) time.Time {
	zeit, err := time.Parse("2006-01-02 15:04:05-07", aString)
	if err != nil {
		log.Fatalf("time.Parse(%s) caused %v", aString, err)
	}
	return zeit
}

func parseFloat(aString string) float32 {
	float, err := strconv.ParseFloat(aString, 32)
	if err != nil {
		log.Fatalf("strconv.ParseFloat(%s) caused %v", aString, err)
	}
	return float32(float)
}

func parseBool(aString string) bool {
	return aString != "0"
}
