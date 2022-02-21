package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	Commit = "HEAD"
)

func main() {
	upSince := time.Now()
	tickersFile := "tickers.csv"
	// read the CSV file
	// use FAKESTOCK_TICKERS environment variable to override the default if present
	if os.Getenv("FAKESTOCK_TICKERS") != "" {
		tickersFile = os.Getenv("FAKESTOCK_TICKERS")
	}
	tickers, err := LoadTickers(tickersFile)
	if err != nil {
		panic(err)
	}

	quit := make(chan struct{})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Println("\nReceived an interrupt, stopping services...")
			quit <- struct{}{}
		}
	}()

	e := echo.New()
	e.GET("/tickers", func(c echo.Context) error {
		// return all tickers
		return c.JSON(http.StatusOK, tickers)
	})
	e.GET("/tickers/:symbol", func(c echo.Context) error {
		// return ticker for the given symbol
		symbol := c.Param("symbol")
		ticker, ok := tickers[symbol]
		if !ok {
			return c.JSON(http.StatusNotFound, "Ticker not found")
		}

		return c.JSON(http.StatusOK, ticker)
	})
	e.GET("/exchanges", func(c echo.Context) error {
		// return all exchanges
		return c.JSON(http.StatusOK, map[string]interface{}{
			"NYSE":   NYSE,
			"NASDAQ": NASDAQ,
			"AMEX":   AMEX,
		})
	})
	e.GET("/_ping", func(c echo.Context) error {
		// return the current commit hash
		return c.JSON(http.StatusOK, map[string]interface{}{
			"commit": Commit,
			"uptime": time.Since(upSince).String(),
		})
	})

	go func() {
		e.Logger.Fatal(e.Start("0.0.0.0:8080"))
	}()

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, ticker := range tickers {
					ticker.UpdatePrice()
				}
			}
		}
	}()

	<-quit
	ticker.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	return
}
