package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type DataScraped struct {
	ID            string   `json:"id"`
	Segment       string   `json:"segmento"`
	Industry      string   `json:"setor_atuacao"`
	Cash          string   `json:"caixa"`
	GrossDebit    string   `json:"divida_bruta"`
	DebitByEBITDA string   `json:"divida_liquida_por_ebitda"`
	Profit        Profit   `json:"lucro_liquido"`
	ROE           ROE      `json:"roe"`
	Years         []string `json:"-"`
}

type ROE struct {
	LastTwelveMonths string         `json:"ultimos_12meses"`
	LastFiveYears    map[int]string `json:"ultimos_5anos"`
	AvgLastFiveYears string         `json:"media_ultimos_5anos"`
}

type Profit struct {
	LastTwelveMonths string            `json:"ultimos_12meses"`
	LastFiveYears    map[string]string `json:"ultimos_5anos"`
	AvgLastFiveYears string            `json:"media_ultimos_5anos"`
}

func (d *DataScraped) IsNil() bool {
	ok := d.Segment == ""
	ok1 := d.Industry == ""
	ok2 := d.Cash == ""
	return ok && ok1 && ok2
}

func (d *DataScraped) CalculateProfitAverage() error {
	sum := 0.0
	for _, v := range d.Profit.LastFiveYears {
		vv, err := sanitizeProfitValue(v)
		if err != nil {
			return err
		}
		vf, err := strconv.ParseFloat(vv, 64)
		if err != nil {
			return err
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
		v1 := strings.ReplaceAll(v, ",", ".")
		v2 := strings.ReplaceAll(v1, " ", "")
		vf, err := strconv.ParseFloat(v2, 64)
		if err != nil {
			return err
		}
		sum = vf + sum
	}
	avg := sum / 5.0
	d.ROE.AvgLastFiveYears = fmt.Sprintf("%.2f", avg)
	return nil
}

func sanitizeProfitValue(s string) (string, error) {
	re, err := regexp.Compile(`[A-Z]`)
	if err != nil {
		return "", err
	}
	ss := re.ReplaceAllString(s, "")
	ss = strings.ReplaceAll(ss, ".", "")
	ss = strings.ReplaceAll(ss, " ", "")
	ss = strings.ReplaceAll(ss, ",", ".")
	return ss, nil
}
