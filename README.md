# StockScrap

Simple webscraper with a very creative name that fetches information about [B3](www.b3.com.br) stocks from [Status Invest](https://statusinvest.com.br).

## Installation

```sh
go build -o sscrap .
```


## Usage example

```sh
$ sscrap STCK3
{
  "id": "STCK3",
  "segmento": "Novo Mercado",
  "setor_atuacao": "Saúde",
  "fluxo_de_caixa": "224,23 M",
  "divida_bruta": "2.071,87 M",
  "divida_liquida_por_ebitda": "0,14",
  "lucro_liquido": {
    "ultimos_12meses": "851,85 M",
    "ultimos_5anos": {
      "2015": " 311,33 M",
      "2016": "456,49 M",
      "2017": "650,60 M",
      "2018": "788,33 M",
      "2019": "851,85 M"
    },
    "media_ultimos_5anos": "611.72"
  },
  "roe": {
    "ultimos_12meses": "10,92",
    "ultimos_5anos": {
      "2015": "93,76",
      "2016": "89,40",
      "2017": "137,84",
      "2018": "21,86",
      "2019": "11,73"
    },
    "media_ultimos_5anos": "70.92"
  }
}

```
