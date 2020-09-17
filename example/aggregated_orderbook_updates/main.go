package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/Rellum/valr_ws"
)

var keyID = flag.String("k", "", "Valr API Key ID")
var keySecret = flag.String("s", "", "Valr API Key secret")

func main() {
	flag.Parse()

	fn := func(update valr_ws.AggregatedOrderbookUpdate) {
		log.Printf("%+v", update)
	}

	err := valr_ws.NewAggregatedOrderbookUpdatesStream(context.Background(), *keyID, *keySecret, []string{"BTCZAR"}, fn)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(30 * time.Second)
}
