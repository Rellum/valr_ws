package main

import (
	"flag"
	"log"
	"time"

	"github.com/Rellum/valr_ws"
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

	c.OnAggregatedOrderbookUpdates([]string{"BTCZAR"}, func(update valr_ws.AggregatedOrderbookUpdate) {
		log.Printf("%+v", update)
	})

	time.Sleep(30 * time.Second)
	c.Close()
}
