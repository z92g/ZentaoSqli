package main

import (
	"flag"
	"fmt"
)

var host, file string

func main() {
	zentao := NewZentaoSQLI()

	flag.StringVar(&host, "h", "", "ip")
	flag.StringVar(&file, "f", "", "filepath")
	flag.Parse()

	view := `
 ____  ____  _  _  ____    __    _____    ___  _____  __    ____ 
(_   )( ___)( \( )(_  _)  /__\  (  _  )  / __)(  _  )(  )  (_  _)
 / /_  )__)  )  (   )(   /(__)\  )(_)(   \__ \ )(_)(  )(__  _)(_ 
(____)(____)(_)\_) (__) (__)(__)(_____)  (___/(___/\\(____)(____)  by:Z92G`
	fmt.Println(view)
	fmt.Println()

	if file != "" && host == "" {
		zentao.batchSqliScan(file)
	} else {
		zentao.sqliScan(host)
	}
}
