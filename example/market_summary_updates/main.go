package main

import (
	"flag"
	"log"
	"time"
	valr_ws "valr_go"
)

var keyID = flag.String("k", "", "Valr API Key ID")
var keySecret = flag.String("s", "", "Valr API Key secret")

func main() {
	flag.Parse()
	c := valr_ws.New(*keyID, *keySecret)

	err := c.Connect()
	if err != nil {
		log.Fatal(err)
	}

	c.OnMarketSummaryUpdates([]string{"BTCZAR"}, func(update valr_ws.MarketSummaryUpdate) {
		log.Printf("%+v", update)
	})

	time.Sleep(30 * time.Second)
	c.Close()
}
