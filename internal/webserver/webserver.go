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
	var resp map[string]interface{}
	stockName := strings.ReplaceAll(r.URL.String(), "/", "")
	stockName = strings.ToUpper(stockName)
	log.Printf("fetching information about: %v\n", stockName)
	data, err := scraper.Scrap(stockName)
	if err != nil {
		resp = map[string]interface{}{
			"success": false,
			"error":   "Failed at fetching stock information. You can retry, if error persists, contact me!",
		}
	} else {
		resp = map[string]interface{}{
			"success": true,
			"data":    data,
		}
	}
	j, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Println("[ERROR] ", err)
	}
	fmt.Fprintf(w, string(j))
}

type Response struct {
	Success bool
	Data    string
}
