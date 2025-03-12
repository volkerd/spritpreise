package cmd

import (
	"github.com/spf13/cobra"
	"log"
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
			processAll(dateFlag)
		},
	}
	rootCmd.PersistentFlags().StringP("nach", "n", "", "Lies Daten NACH diesem Datum")
}

func Exec() {

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}

}
