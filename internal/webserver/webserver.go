package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sscrap/internal/scraper"
	"strings"
)

func Listen() {
	fmt.Println("running as http server...")
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	p := ":" + port
	err := http.ListenAndServe(p, nil)
	if err != nil {
		log.Fatal("[HTTP_ERROR]: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	stockName := strings.ReplaceAll(r.URL.String(), "/", "")
	stockName = strings.ToUpper(stockName)
	log.Printf("fetching information about: %v\n", stockName)
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
