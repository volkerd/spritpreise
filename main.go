package main

import (
	"github.com/dimiro1/banner"
	"github.com/volkerd/spritpreise/cmd"
	"os"
)

func main() {
	cmd.Exec()
}

func init() {
	templ := `{{ "spritpreise" }} {{ .GoVersion }}  {{ .GOOS }}/{{ .GOARCH }} {{ .Now "2.1.2006 15:04" }}
`
	banner.InitString(os.Stdout, true, false, templ)
}
