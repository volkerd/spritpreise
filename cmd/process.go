package cmd

import (
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
)

func processPair(pair FilePair) (int, int) {
	stations := readStations(pair.station.path)
	stationFilter := createStationFilter(stations)
	prices := readPrices(pair.price.path, stationFilter)
	fmt.Printf("%v: %d prices found for %d stations\n", pair.price.date, len(prices), len(stations))
	writeData(prices, stationFilter)
	return len(stations), len(prices)
}

func writeData(prices []Price, stationFilter map[string]Station) {
	client := influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN)
	writeAPI := client.WriteAPI(INFLUXDB_ORG, INFLUXDB_BUCKET)
	//ctx := context.Background()
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
	client.Close()
}
