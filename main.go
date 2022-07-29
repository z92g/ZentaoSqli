package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
)

func main() {
	var host, file string

	zentao := &Zentao{}

	flag.StringVar(&host, "h", "", "ip")
	flag.StringVar(&file, "f", "", "filepath")
	flag.Parse()

	view := `
 ____  ____  _  _  ____    __    _____    ___  _____  __    ____ 
(_   )( ___)( \( )(_  _)  /__\  (  _  )  / __)(  _  )(  )  (_  _)
 / /_  )__)  )  (   )(   /(__)\  )(_)(   \__ \ )(_)(  )(__  _)(_ 
(____)(____)(_)\_) (__) (__)(__)(_____)  (___/(___/\\(____)(____)  by:Z92G`
	color.Cyan.Println(view)
	fmt.Println()

	if file != "" && host == "" {
		zentao.batchScan(file)
	} else {
		zentao.singleScan(host)
	}
}
