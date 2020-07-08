package main

import (
	"flag"
	"os"
	"sscrap/internal/cli"
	"sscrap/internal/webserver"
)

var httpMode bool

func init() {
	flag.BoolVar(&httpMode, "http", false, "When informed, app will run as a http server")
	flag.Parse()
}

func main() {
	if httpMode {
		webserver.Listen()
		os.Exit(0)
	}

	// cli mode
	cli.Run()
	os.Exit(0)
}
