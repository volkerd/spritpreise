package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/volkerd/spritpreise/pkg/common"
	"log"
	"time"
)

var rootCmd *cobra.Command

func Last[E any](s []E) E {
	if len(s) == 0 {
		var zero E
		return zero
	}
	return s[len(s)-1]
}

func init() {
	rootCmd = &cobra.Command{
		Use:   "spritpreise",
		Short: "Kopiere Spritpreise von tankerkoenig.de in InfluxDB",
		Long:  `Kopiere Spritpreise von tankerkoenig.de in InfluxDB`,
		Run: func(cmd *cobra.Command, args []string) {
			dateFlag, err := cmd.Flags().GetString("nach")
			if err != nil {
				log.Fatal(err)
			}
			_, _ = getFileList(dateFlag)
		},
	}
	rootCmd.PersistentFlags().StringP("nach", "s", "2025-03-01", "Erster Tag")
}

func Exec() {

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func getFileList(dateFlag string) (stationFileList []File, priceFileList []File) {
	afterDate, err := time.Parse("2006-01-02", dateFlag)
	if err != nil {
		log.Fatal(err)
	}
	stations := stationFiles(afterDate)
	fmt.Printf("%v stations found, start %v, end %v\n", len(stations), stations[0], Last(stations))
	prices := priceFiles(afterDate)
	fmt.Printf("%v prices found, start %v, end %v\n", len(prices), prices[0], Last(prices))
	cutOffDate, _ := time.Parse("2006-01-02", common.CUT_OFF_DATE)
	if afterDate.After(cutOffDate) && len(stations) != len(prices) {
		log.Fatalf("number of stations (%v) must be equal to number of prices (%v)", len(stations), len(prices))
	}
	return stations, prices
}
