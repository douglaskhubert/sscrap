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
  "segmento": "Nível 2",
  "setor_atuacao": "Petróleo. Gás e Biocombustíveis",
  "fluxo_de_caixa": "29.729,00 M",
  "divida_bruta": "463.916,00 M",
  "divida_liquida_por_ebitda": "4,12",
  "lucro_liquido_ultimos_12meses": "40.970,00 M",
  "roe": "-5,53"
}

```
