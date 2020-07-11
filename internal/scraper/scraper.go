package scraper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const BaseEndpoint = "https://statusinvest.com.br/acoes/"

var DataScrapedNil = DataScraped{}

func Scrap(stockName string) (DataScraped, error) {
	endpoint := BaseEndpoint + stockName
	resp, err := http.Get(endpoint)
	if err != nil {
		return DataScrapedNil, errors.New(
			fmt.Sprintf("failed at getting: %v, error: %v", endpoint, err))
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return DataScraped{}, errors.New(fmt.Sprintf(
			"failed at reading response body: %v", err))
	}

	var data DataScraped

	data.ID = strings.ToUpper(stockName)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			if data.Segment == "" {
				data.Segment = tryScrapSegment(n)
			}
		}
		if n.Type == html.ElementNode && n.Data == "a" {
			if data.Industry == "" {
				data.Industry = tryScrapIndustry(n)
			}
		}

		if n.Type == html.ElementNode && n.Data == "span" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "d-block" {
					t := n.FirstChild
					header := strings.ReplaceAll(renderNode(t), "\n", "")

					if fn, ok := scrapFuncByHeader[header]; ok && fn != nil {
						data = fn(t, data)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if data.IsNil() {
		return DataScrapedNil, errors.New("couldn't find stock information")
	}

	return data, err
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func tryScrapSegment(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "title" && a.Val == "Segmento de listagem na B3" {
			t := n.Parent.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild
			return renderNode(t)
		}
	}
	return ""
}

func tryScrapIndustry(n *html.Node) string {
	for _, a := range n.Attr {
		searchFor := "Ver outras empresas do setor "
		if a.Key == "title" && strings.Contains(a.Val, searchFor) {
			industry := strings.ReplaceAll(a.Val, searchFor, "")
			return strings.ReplaceAll(industry, "'", "")
		}
	}
	return ""
}
