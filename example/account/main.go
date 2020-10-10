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

	anyFn := func(typ string, bl []byte) {
		log.Printf("anyFn - %s: %s", typ, bl)
	}

	balanceUpdateFn := func(update valr_ws.BalanceUpdate) {
		log.Printf("balanceUpdateFn - %+v", update)
	}

	openOrdersUpdateFn := func(update []valr_ws.OpenOrderUpdate) {
		log.Printf("openOrdersUpdateFn - %+v", update)
	}

	onOrderStatusUpdateFn := func(update valr_ws.OrderStatusUpdate) {
		log.Printf("onOrderStatusUpdateFn - %+v", update)
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	for {
		if ctx.Err() != nil {
			break
		}

		err := valr_ws.NewAccountStream(ctx, *keyID, *keySecret,
			valr_ws.OnAny(anyFn),
			valr_ws.OnBalanceUpdate(balanceUpdateFn),
			valr_ws.OnOpenOrdersUpdate(openOrdersUpdateFn),
			valr_ws.OnOrderStatusUpdate(onOrderStatusUpdateFn),
		)
		if err != nil {
			log.Println(err)
		}
	}
}
