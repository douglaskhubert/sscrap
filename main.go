package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) == 1 {
		usage()
		os.Exit(0)
	}
	stockName := os.Args[1]
	endpoint := "https://statusinvest.com.br/acoes/" + stockName
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Fatalf("failed at getting: %v, error: %v", endpoint, err)
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("failed at reading response body: %v", err)
	}
	var data DataScraped
	data.ID = strings.ToUpper(stockName)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "title" && a.Val == "Segmento de listagem na B3" {
					t := n.Parent.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild
					data.Segment = renderNode(t)
					break
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				searchFor := "Ver outras empresas do setor "
				if a.Key == "title" && strings.Contains(a.Val, searchFor) {
					industry := strings.ReplaceAll(a.Val, searchFor, "")
					data.Industry = strings.ReplaceAll(industry, "'", "")
					break
				}
			}
		}

		if n.Type == html.ElementNode && n.Data == "span" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "d-block" {
					t := n.FirstChild
					header := strings.ReplaceAll(renderNode(t), "\n", "")
					if header == "Saldo Final de Caixa e Equivalentes - (R$)" {
						tt := t.Parent.Parent.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild
						data.CashFlow = strings.ReplaceAll(renderNode(tt), "\n", "")
						break
					}
					if header == "Dívida Bruta - (R$)" {
						tt := t.Parent.Parent.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild
						data.GrossDebit = strings.ReplaceAll(renderNode(tt), "\n", "")
						break
					}
					if header == "Dívida Líquida/Ebitda" {
						tt := t.Parent.Parent.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild
						data.DebitByEBITDA = strings.ReplaceAll(renderNode(tt), "\n", "")
						break
					}
					if header == "Lucro Líquido - (R$)" {
						tt := t.Parent.Parent.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild
						data.Last12MonthsProfit = strings.ReplaceAll(renderNode(tt), "\n", "")
						break
					}
					if header == "ROE - (%)" {
						tt := t.Parent.Parent.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild
						data.ROE = strings.ReplaceAll(renderNode(tt), "\n", "")
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	j, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(j))
}
func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

type DataScraped struct {
	ID                 string `json:"id"`
	Segment            string `json:"segmento"`
	Industry           string `json:"setor_atuacao"`
	CashFlow           string `json:"fluxo_de_caixa"`
	GrossDebit         string `json:"divida_bruta"`
	DebitByEBITDA      string `json:"divida_liquida_por_ebitda"`
	Last12MonthsProfit string `json:"lucro_liquido_ultimos_12meses"`
	ROE                string `json:"roe"`
}

func usage() {
	fmt.Println("Usage: \"" + os.Args[0] + " <STOCK_CODE>\". Try PETR3, VALE3, etc.")
}
