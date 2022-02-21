package main

import (
	"encoding/csv"
	"math/rand"
	"os"
)

type Exchange int

var (
	NYSE   Exchange = 0
	NASDAQ Exchange = 1
	AMEX   Exchange = 2
)

type Ticker struct {
	Symbol   string   `json:"symbol"`
	Name     string   `json:"name"`
	Exchange Exchange `json:"exchange"`
	Price    float64  `json:"price"`
}

func LoadTickers(filePath string) (map[string]*Ticker, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*Ticker)
	for _, record := range records {
		symbol := record[0]
		result[symbol] = &Ticker{
			Symbol: record[0],
			Name:   record[1],
			Price:  0.0,
		}

		switch record[3] {
		case "NYSE":
			result[symbol].Exchange = NYSE
		case "NASDAQ":
			result[symbol].Exchange = NASDAQ
		case "AMEX":
			result[symbol].Exchange = AMEX
		}
	}

	return result, nil
}

func (t *Ticker) UpdatePrice() {
	currentPrice := t.Price
	if currentPrice == 0.0 {
		// this is a new ticker. come up with a random price between 1 and 1000
		currentPrice = float64(rand.Intn(1000))
	} else {
		// we have an existing price. decide randomly if it should go up or down
		direction := rand.Intn(3)
		if direction == 0 {
			// go up
			currentPrice += float64(rand.Intn(10))
		} else if direction == 1 {
			// go down
			currentPrice -= float64(rand.Intn(10))
		} else {
			// stay the same
		}
	}

	// never go above 1000 or below 1
	if currentPrice > 1000 {
		currentPrice = 1000
	}
	if currentPrice < 1 {
		currentPrice = 1
	}

	t.Price = currentPrice
}
