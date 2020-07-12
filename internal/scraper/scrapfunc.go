package scraper

import (
	"log"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type ScrapFunc func(n *html.Node, d DataScraped) DataScraped

const (
	HEADER_CASH            = "Saldo Final de Caixa e Equivalentes - (R$)"
	HEADER_GROSS_DEBIT     = "Dívida Bruta - (R$)"
	HEADER_DEBIT_BY_EBITDA = "Dívida Líquida/Ebitda"
	HEADER_PROFIT          = "Lucro Líquido - (R$)"
	HEADER_ROE             = "ROE - (%)"
)

var scrapFuncByHeader = map[string]ScrapFunc{
	HEADER_CASH:            cashFunc,
	HEADER_GROSS_DEBIT:     grossDebitFunc,
	HEADER_DEBIT_BY_EBITDA: debitByEBITDAFunc,
	HEADER_PROFIT:          profitFunc,
	HEADER_ROE:             roeFunc,
}

var _ = ScrapFunc(cashFunc)
var _ = ScrapFunc(grossDebitFunc)
var _ = ScrapFunc(debitByEBITDAFunc)
var _ = ScrapFunc(profitFunc)
var _ = ScrapFunc(roeFunc)

func cashFunc(n *html.Node, d DataScraped) DataScraped {
	t := n.Parent.Parent.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild
	d.Cash = strings.ReplaceAll(renderNode(t), "\n", "")
	return d
}

func grossDebitFunc(n *html.Node, d DataScraped) DataScraped {
	tt := n.Parent.Parent.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild
	d.GrossDebit = strings.ReplaceAll(renderNode(tt), "\n", "")
	return d
}

func debitByEBITDAFunc(n *html.Node, d DataScraped) DataScraped {
	tt := n.Parent.Parent.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild
	d.DebitByEBITDA = strings.ReplaceAll(renderNode(tt), "\n", "")
	return d
}

func profitFunc(n *html.Node, d DataScraped) DataScraped {
	elem12Months := n.Parent.Parent.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild
	d.Profit.LastTwelveMonths = strings.ReplaceAll(renderNode(elem12Months), "\n", "")
	td12Months := elem12Months.Parent.Parent
	for _, a := range td12Months.Attr {
		if a.Key == "class" && strings.Contains(a.Val, "lastTwelveMonths") {
			// proceed
			spanYearMinusOne,
				spanYearMinusTwo,
				spanYearMinusThree,
				spanYearMinusFour,
				spanYearMinusFive := getLast5YearValues(td12Months)

			d.Profit.LastFiveYears = map[string]string{
				d.Years[0]: strings.ReplaceAll(renderNode(spanYearMinusOne), "\n", ""),
				d.Years[1]: strings.ReplaceAll(renderNode(spanYearMinusTwo), "\n", ""),
				d.Years[2]: strings.ReplaceAll(renderNode(spanYearMinusThree), "\n", ""),
				d.Years[3]: strings.ReplaceAll(renderNode(spanYearMinusFour), "\n", ""),
				d.Years[4]: strings.ReplaceAll(renderNode(spanYearMinusFive), "\n", ""),
			}
			err := d.CalculateProfitAverage()
			if err != nil {
				log.Fatalf("failed at calculating last 5 years profit average, error: %v", err)
			}
		}
	}
	return d
}

func roeFunc(n *html.Node, d DataScraped) DataScraped {
	spanLastTwelveMonths := n.Parent.Parent.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild
	d.ROE.LastTwelveMonths = strings.ReplaceAll(renderNode(spanLastTwelveMonths), "\n", "")
	tdLastTwelveMonths := spanLastTwelveMonths.Parent.Parent
	for _, a := range tdLastTwelveMonths.Attr {
		if a.Key == "class" && strings.Contains(a.Val, "lastTwelveMonths") {
			spanYearMinusOne,
				spanYearMinusTwo,
				spanYearMinusThree,
				spanYearMinusFour,
				spanYearMinusFive := getLast5YearValues(tdLastTwelveMonths)
			year := time.Now().Year()
			d.ROE.LastFiveYears = map[int]string{
				year - 1: sanitizeFloat(renderNode(spanYearMinusOne)),
				year - 2: sanitizeFloat(renderNode(spanYearMinusTwo)),
				year - 3: sanitizeFloat(renderNode(spanYearMinusThree)),
				year - 4: sanitizeFloat(renderNode(spanYearMinusFour)),
				year - 5: sanitizeFloat(renderNode(spanYearMinusFive)),
			}
			err := d.CalculateROEAverage()
			if err != nil {
				log.Printf("failed at calculating last 5 years ROE average, error: %v\n", err)
			}
		}
	}
	return d
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

func sanitizeFloat(s string) string {
	ss := strings.ReplaceAll(s, "\n", "")
	sss := strings.ReplaceAll(ss, " ", "")
	return sss
}
