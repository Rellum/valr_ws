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

	fn := func(update valr_ws.MarketSummaryUpdate) {
		log.Printf("%+v", update)
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	for {
		if ctx.Err() != nil {
			break
		}

		err := valr_ws.NewMarketSummaryUpdatesStream(ctx, *keyID, *keySecret, []string{"BTCZAR"}, fn)
		if err != nil {
			log.Println(err)
		}
	}
}
