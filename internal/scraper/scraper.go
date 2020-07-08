package scraper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

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
						elem12Months := t.Parent.Parent.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild
						data.Profit.LastTwelveMonths = strings.ReplaceAll(renderNode(elem12Months), "\n", "")
						td12Months := elem12Months.Parent.Parent
						for _, a := range td12Months.Attr {
							if a.Key == "class" && strings.Contains(a.Val, "lastTwelveMonths") {
								// proceed
								spanYearMinusOne,
									spanYearMinusTwo,
									spanYearMinusThree,
									spanYearMinusFour,
									spanYearMinusFive := getLast5YearValues(td12Months)

								year := time.Now().Year()

								data.Profit.LastFiveYears = map[int]string{
									year - 1: strings.ReplaceAll(renderNode(spanYearMinusOne), "\n", ""),
									year - 2: strings.ReplaceAll(renderNode(spanYearMinusTwo), "\n", ""),
									year - 3: strings.ReplaceAll(renderNode(spanYearMinusThree), "\n", ""),
									year - 4: strings.ReplaceAll(renderNode(spanYearMinusFour), "\n", ""),
									year - 5: strings.ReplaceAll(renderNode(spanYearMinusFive), "\n", ""),
								}
								err := data.CalculateProfitAverage()
								if err != nil {
									log.Fatalf("failed at calculating last 5 years profit average, error: %v", err)
								}
							}
						}
						break
					}
					if header == "ROE - (%)" {
						spanLastTwelveMonths := t.Parent.Parent.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild
						data.ROE.LastTwelveMonths = strings.ReplaceAll(renderNode(spanLastTwelveMonths), "\n", "")
						tdLastTwelveMonths := spanLastTwelveMonths.Parent.Parent
						for _, a := range tdLastTwelveMonths.Attr {
							if a.Key == "class" && strings.Contains(a.Val, "lastTwelveMonths") {
								spanYearMinusOne,
									spanYearMinusTwo,
									spanYearMinusThree,
									spanYearMinusFour,
									spanYearMinusFive := getLast5YearValues(tdLastTwelveMonths)
								year := time.Now().Year()
								data.ROE.LastFiveYears = map[int]string{
									year - 1: strings.ReplaceAll(renderNode(spanYearMinusOne), "\n", ""),
									year - 2: strings.ReplaceAll(renderNode(spanYearMinusTwo), "\n", ""),
									year - 3: strings.ReplaceAll(renderNode(spanYearMinusThree), "\n", ""),
									year - 4: strings.ReplaceAll(renderNode(spanYearMinusFour), "\n", ""),
									year - 5: strings.ReplaceAll(renderNode(spanYearMinusFive), "\n", ""),
								}
								err := data.CalculateROEAverage()
								if err != nil {
									log.Fatalf("failed at calculating last 5 years ROE average, error: %v", err)
								}
							}
						}
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
	return data, nil
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

type DataScraped struct {
	ID            string `json:"id"`
	Segment       string `json:"segmento"`
	Industry      string `json:"setor_atuacao"`
	CashFlow      string `json:"fluxo_de_caixa"`
	GrossDebit    string `json:"divida_bruta"`
	DebitByEBITDA string `json:"divida_liquida_por_ebitda"`
	Profit        Profit `json:"lucro_liquido"`
	ROE           ROE    `json:"roe"`
}

type ROE struct {
	LastTwelveMonths string         `json:"ultimos_12meses"`
	LastFiveYears    map[int]string `json:"ultimos_5anos"`
	AvgLastFiveYears string         `json:"media_ultimos_5anos"`
}

type Profit struct {
	LastTwelveMonths string         `json:"ultimos_12meses"`
	LastFiveYears    map[int]string `json:"ultimos_5anos"`
	AvgLastFiveYears string         `json:"media_ultimos_5anos"`
}

func (d *DataScraped) CalculateProfitAverage() error {
	sum := 0.0
	for _, v := range d.Profit.LastFiveYears {
		vv := sanitizeProfitValue(v)
		vf, err := strconv.ParseFloat(vv, 64)
		if err != nil {
			log.Fatal(err)
		}
		sum = vf + sum
	}
	avg := sum / 5.0
	d.Profit.AvgLastFiveYears = fmt.Sprintf("%.2f", avg)
	return nil
}

func (d *DataScraped) CalculateROEAverage() error {
	sum := 0.0
	for _, v := range d.ROE.LastFiveYears {
		vv := strings.ReplaceAll(v, " ", "")
		vv = strings.ReplaceAll(v, ",", ".")
		vf, err := strconv.ParseFloat(vv, 64)
		if err != nil {
			log.Fatal(err)
		}
		sum = vf + sum
	}
	avg := sum / 5.0
	d.ROE.AvgLastFiveYears = fmt.Sprintf("%.2f", avg)
	return nil
}

func sanitizeProfitValue(s string) string {
	re, err := regexp.Compile(`[A-Z]`)
	if err != nil {
		log.Fatal(err)
	}
	ss := re.ReplaceAllString(s, "")
	ss = strings.ReplaceAll(ss, ".", "")
	ss = strings.ReplaceAll(ss, " ", "")
	ss = strings.ReplaceAll(ss, ",", ".")
	return ss
}
func getLast5YearValues(n *html.Node) (minusOne, minusTwo, minusThree, minusFour, minusFive *html.Node) {
	tdYearMinusOne := n.NextSibling.NextSibling.NextSibling.NextSibling
	minusOne = tdYearMinusOne.FirstChild.NextSibling.FirstChild

	tdYearMinusTwo := tdYearMinusOne.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling
	minusTwo = tdYearMinusTwo.FirstChild.NextSibling.FirstChild

	tdYearMinusThree := tdYearMinusTwo.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling
	minusThree = tdYearMinusThree.FirstChild.NextSibling.FirstChild

	tdYearMinusFour := tdYearMinusThree.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling
	minusFour = tdYearMinusFour.FirstChild.NextSibling.FirstChild

	tdYearMinusFive := tdYearMinusFour.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling
	minusFive = tdYearMinusFive.FirstChild.NextSibling.FirstChild
	return
}
