package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"math"
)

var rootCmd *cobra.Command

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
			pairs := getPairs(dateFlag)
			var stationMin = math.MaxInt32
			var stationMax int
			var priceSum int
			for _, pair := range pairs {
				stationCount, priceCount := processPair(pair)
				stationMin = minimum(stationMin, stationCount)
				stationMax = maximum(stationMax, stationCount)
				priceSum += priceCount
			}
			fmt.Printf("%d prices found\n", priceSum)
			fmt.Printf("%d to %d stations found\n", stationMin, stationMax)
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
