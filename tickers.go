package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Exchange int

var (
	NYSE   Exchange = 0
	NASDAQ Exchange = 1
)

type Ticker struct {
	Symbol        string   `json:"symbol"`
	Name          string   `json:"name"`
	Exchange      Exchange `json:"exchange"`
	Price         float64  `json:"price"`
	NetChange     float64  `json:"net_change"`
	PercentChange float64  `json:"percent_change"`
	MarketCap     float64  `json:"market_cap"`
	Country       string   `json:"country"`
	IPOYear       string   `json:"ipo_year"`
	Volume        float64  `json:"volume"`
	Sector        string   `json:"sector"`
	Industry      string   `json:"industry"`
	TotalStock    int64    `json:"total_stock"`
}

func LoadStartPrices(filePath string, exchange Exchange, tickers map[string]*Ticker) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	fmt.Printf("Loading start prices from %s\n", filePath)
	fmt.Printf("Found %d tickers\n", len(records))

	for _, record := range records {
		symbol := record[0]

		// do we have this ticker?
		if _, ok := tickers[symbol]; ok {
			fmt.Printf("Skipping ticker %s\n", symbol)
			continue
		}

		ticker := &Ticker{
			Symbol:   symbol,
			Exchange: exchange,
			Name:     record[1],
		}

		// update price
		priceInDollars := record[2]
		// drop the dollar sign and strip spaces
		priceInDollars = priceInDollars[1:]
		priceInDollars = strings.Replace(priceInDollars, " ", "", -1)
		priceInDollars = strings.Replace(priceInDollars, ",", "", -1)

		price, err := strconv.ParseFloat(priceInDollars, 64)
		if err != nil {
			fmt.Printf("error parsing price %s for ticker %s\n", priceInDollars, symbol)
			continue
		}
		ticker.Price = price

		// update net change
		netChange, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			fmt.Printf("error parsing net change %s for ticker %s\n", record[3], symbol)
			continue
		}
		ticker.NetChange = netChange

		// update percent change
		percentChangeString := record[4]
		if percentChangeString != "" {
			// drop the percent sign and strip spaces
			percentChangeString = percentChangeString[:len(percentChangeString)-1]
			percentChangeString = strings.Replace(percentChangeString, " ", "", -1)
			percentChange, err := strconv.ParseFloat(percentChangeString, 64)
			if err != nil {
				fmt.Printf("error parsing percent change %s for ticker %s\n", percentChangeString, symbol)
				continue
			}
			ticker.PercentChange = percentChange
		}

		// update market cap
		marketCap, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			continue
		}
		ticker.MarketCap = marketCap

		// update country
		ticker.Country = record[6]

		// update ipo year
		ticker.IPOYear = record[7]

		// update volume
		volume, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			fmt.Printf("error parsing volume %s for ticker %s\n", record[8], symbol)
			continue
		}
		ticker.Volume = volume

		// update sector
		ticker.Sector = record[9]

		// update industry
		ticker.Industry = record[10]

		ticker.TotalStock = int64(ticker.MarketCap / ticker.Price)

		tickers[symbol] = ticker
	}

	return nil
}

func (t *Ticker) UpdatePrice() {
	oldPrice := t.Price
	currentPrice := t.Price
	direction := rand.Float64()
	if direction <= 0.2 {
		// go up by a number between 0 and 5%
		change := rand.Float64() * 2
		currentPrice = math.Round((oldPrice+change)*100) / 100
	} else if direction <= 0.4 {
		// go down by a number between 0 and 5%
		change := rand.Float64() * 2
		currentPrice = math.Round((oldPrice-change)*100) / 100
	} else {
		// stay the same
		return
	}

	// never go above below 0.01
	if currentPrice < 0.01 {
		currentPrice = 0.01
	}

	// calculate net change
	netChange := math.Round((currentPrice-oldPrice)*100) / 100
	percentChange := math.Round((netChange/currentPrice)*10000) / 100

	t.Price = currentPrice
	t.NetChange = netChange
	t.PercentChange = percentChange
	t.MarketCap = math.Round(currentPrice * float64(t.TotalStock))
}
