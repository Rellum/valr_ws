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

	fn := func(typ string, bl []byte) {
		log.Printf("%s: %s", typ, bl)
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	for {
		if ctx.Err() != nil {
			break
		}

		err := valr_ws.NewAccountStream(ctx, *keyID, *keySecret, fn)
		if err != nil {
			log.Println(err)
		}
	}
}
