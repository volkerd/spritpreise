package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

func getPairs(afterDate string) (pairList []FilePair) {
	stationFiles, priceFiles := getFileList(afterDate)
	pairList = mergeFileLists(stationFiles, priceFiles)
	return pairList
}

func getFileList(dateFlag string) (stationFileList []File, priceFileList []File) {
	_, err := time.Parse("2006-01-02", dateFlag)
	if err != nil {
		log.Fatal(err)
	}
	stations := stationFiles(dateFlag)
	fmt.Printf("%v stations found, start %v, end %v\n", len(stations), First(stations), Last(stations))
	prices := priceFiles(dateFlag)
	fmt.Printf("%v prices found, start %v, end %v\n", len(prices), First(prices), Last(prices))
	if dateFlag > CUT_OFF_DATE && len(stations) != len(prices) {
		log.Fatalf("number of stations (%v) must be equal to number of prices (%v)", len(stations), len(prices))
	}
	return stations, prices
}

func stationFiles(afterDate string) (paths []File) {
	return traversFS(afterDate, STATION, STATION_PATH)
}

func priceFiles(afterDate string) (paths []File) {
	return traversFS(afterDate, PRICE, PRICE_PATH)
}

func traversFS(afterDate string, dataType string, root string) []File {
	fmt.Printf("TraversFS(%v,%v,%v)\n", afterDate, dataType, root)
	r, err := regexp.Compile("([0-9]{4}-[0-9]{2}-[0-9]{2})-([a-z]*)\\.")
	if err != nil {
		log.Fatal(err)
	}
	fsys := os.DirFS(basePath)
	paths := []File{}
	err = fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) == ".csv" {
			var date string
			m := r.FindStringSubmatch(filepath.Base(path))
			if len(m) == 3 {
				date = m[1]
			} else {
				date = CUT_OFF_DATE
			}
			if date > afterDate {
				absPath := filepath.Join(basePath, path)
				if err != nil {
					log.Fatal(err)
				}
				file := File{
					path:     absPath,
					date:     date,
					datatype: dataType}
				//fmt.Println(file)
				paths = append(paths, file)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("error walking through files: %v", err)
	}
	sort.SliceStable(paths, func(i, j int) bool { return paths[i].date < paths[j].date })
	return paths
}

func mergeFileLists(stationFileList []File, priceFileList []File) (pairList []FilePair) {
	var pm map[string]File
	pm = make(map[string]File)
	for _, station := range stationFileList {
		pm[station.date] = station
	}
	for _, price := range priceFileList {
		station, exists := pm[price.date]
		if !exists {
			if price.date > CUT_OFF_DATE {
				log.Fatalf("station %v does not exist", price.date)
			} else {
				station, exists = pm[CUT_OFF_DATE]
				if !exists {
					log.Fatalf("station %v does not exist", price.date)
				}
			}
		}
		pair := FilePair{station: station, price: price}
		pairList = append(pairList, pair)
	}
	fmt.Printf("%v pairs found, start %v, end %v\n", len(pairList), First(pairList), Last(pairList))

	return pairList
}
