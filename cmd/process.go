package cmd

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	"math"
	"time"
)

func processAll(dateFlag string) {
	client := influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN)
	defer client.Close()
	queryAPI := client.QueryAPI(INFLUXDB_ORG)
	writeAPI := client.WriteAPI(INFLUXDB_ORG, INFLUXDB_BUCKET)
	defer writeAPI.Flush()

	if dateFlag == "" {
		dateFlag = findLatestEntry(queryAPI)
	}

	pairs := getPairs(dateFlag)
	var stationMin = math.MaxInt32
	var stationMax int
	var priceSum int
	for _, pair := range pairs {
		stationCount, priceCount := processPair(writeAPI, pair)
		stationMin = minimum(stationMin, stationCount)
		stationMax = maximum(stationMax, stationCount)
		priceSum += priceCount
	}
	fmt.Printf("%d prices found\n", priceSum)
	fmt.Printf("%d to %d stations found\n", stationMin, stationMax)
}

// schaut 30 Tage zurueck, liefert immer einen Wert
func findLatestEntry(queryAPI api.QueryAPI) string {
	const DAYS_BACK = 30
	fluxQuery := fmt.Sprintf(
		`from(bucket:"%s")
		|> range(start: -%dd)
		|> last()`,
		INFLUXDB_BUCKET, DAYS_BACK)
	result, err := queryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		log.Fatalf("Error querying influxdb: %v", err)
	}
	latestDate := time.Now().Add(-DAYS_BACK * 24 * time.Hour)
	for result.Next() {
		date := result.Record().Time()
		if date.After(latestDate) {
			latestDate = date
		}
	}
	if result.Err() != nil {
		log.Fatalf("Error querying influxdb: %v", err)
	}
	fmt.Printf("found latest entry at %s\n", latestDate.Format("2006-01-02"))
	return latestDate.Format("2006-01-02")
}

func processPair(writeAPI api.WriteAPI, pair FilePair) (int, int) {
	stations := readStations(pair.station.path)
	stationFilter := createStationFilter(stations)
	prices := readPrices(pair.price.path, stationFilter)
	fmt.Printf("%v: %d prices found for %d stations\n", pair.price.date, len(prices), len(stations))
	writeData(writeAPI, prices, stationFilter)
	return len(stations), len(prices)
}

func writeData(writeAPI api.WriteAPI, prices []Price, stationFilter map[string]Station) {
	for _, price := range prices {
		station := stationFilter[price.station_uuid]
		p := influxdb2.NewPointWithMeasurement("price").
			AddTag("station", station.uuid).
			AddTag("brand", station.brand).
			AddTag("plz", station.post_code).
			SetTime(price.date)
		if price.dieselchange {
			p.AddField("diesel", price.diesel)
		}
		if price.e5change {
			p.AddField("e5", price.e5)
		}
		if price.e10change {
			p.AddField("e10", price.diesel)
		}
		writeAPI.WritePoint(p)
	}
	writeAPI.Flush()
}
