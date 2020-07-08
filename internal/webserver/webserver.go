package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sscrap/internal/scraper"
	"strings"
)

func Listen() {
	fmt.Println("running as http server...")
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("[HTTP_ERROR]: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	stockName := strings.ReplaceAll(r.URL.String(), "/", "")
	log.Println(stockName)
	data, err := scraper.Scrap(stockName)
	if err != nil {
		log.Println("[ERROR] ", err)
	}
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("[ERROR] ", err)
	}
	fmt.Fprintf(w, string(j))
}
