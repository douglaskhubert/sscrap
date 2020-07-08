package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sscrap/internal/scraper"
)

func Run() {
	if len(os.Args) == 1 {
		usage()
		os.Exit(0)
	}
	stockName := os.Args[1]
	data, err := scraper.Scrap(stockName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	j, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(j))
}

func usage() {
	fmt.Println("Usage: \"" + os.Args[0] + " <STOCK_CODE>\". Try PETR3, VALE3, etc.")
}
