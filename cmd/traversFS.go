package cmd

import (
	"fmt"
	"github.com/volkerd/spritpreise/pkg/common"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

type File struct {
	path     string
	date     string
	datatype string
}

func priceFiles(afterDate time.Time) (paths []File) {
	return traversFS(afterDate, common.PRICE, common.PRICE_PATH)
}

func stationFiles(afterDate time.Time) (paths []File) {
	return traversFS(afterDate, common.STATION, common.STATION_PATH)
}

func traversFS(afterDate time.Time, dataType string, root string) (paths []File) {
	fmt.Printf("TraversFS(%v,%v,%v)\n", afterDate.Format("2006-01-02"), dataType, root)
	afterDay := afterDate.Format("2006-01-02")
	r, err := regexp.Compile("([0-9]{4}-[0-9]{2}-[0-9]{2})-([a-z]*)\\.")
	if err != nil {
		log.Fatal(err)
	}
	fsys := os.DirFS(common.BASE_PATH)
	fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) == ".csv" {
			var date string
			m := r.FindStringSubmatch(filepath.Base(path))
			if len(m) == 3 {
				date = m[1]
			} else {
				date = common.CUT_OFF_DATE
			}
			if date > afterDay {
				file := File{
					path:     path,
					date:     date,
					datatype: dataType}
				//fmt.Println(file)
				paths = append(paths, file)
			}
		}

		return nil
	})
	sort.SliceStable(paths, func(i, j int) bool { return paths[i].date < paths[j].date })
	return paths
}
