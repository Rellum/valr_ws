package valr_ws

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const host = "api.valr.com"

func NewMarketSummaryUpdatesStream(ctx context.Context, apiKey, apiSecret string, pairs []string, fn func(MarketSummaryUpdate)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c := &client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	err := c.connect(ctx, "/ws/trade", func(r responseType, bl []byte) {
		if r != responseTypeMarketSummaryUpdate {
			return
		}

		var res marketSummaryUpdateResponse
		err := json.Unmarshal(bl, &res)
		if err != nil {
			log.Println("json unmarshal error:", err)
			return
		}

		fn(res.Data)
	})
	if err != nil {
		return err
	}

	err = subscribe(c.conn, eventTypeMarketSummaryUpdate, pairs)
	if err != nil {
		return err
	}

	return pingForever(ctx, c.conn, c.done)
}

func NewAggregatedOrderbookUpdatesStream(ctx context.Context, apiKey, apiSecret string, pairs []string, fn func(AggregatedOrderbookUpdate)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c := &client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	err := c.connect(ctx, "/ws/trade", func(r responseType, bl []byte) {
		if r != responseTypeAggregatedOrderbookUpdate {
			return
		}

		var res aggregatedOrderbookUpdateResponse
		err := json.Unmarshal(bl, &res)
		if err != nil {
			log.Println("json unmarshal error:", err)
			return
		}

		fn(res.Data)
	})
	if err != nil {
		return err
	}

	err = subscribe(c.conn, eventTypeAggregatedOrderbookUpdate, pairs)
	if err != nil {
		return err
	}

	return pingForever(ctx, c.conn, c.done)
}

func NewTradeBucketStream(ctx context.Context, apiKey, apiSecret string, pairs []string, fn func(NewTradeBucket)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c := &client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	err := c.connect(ctx, "/ws/trade", func(r responseType, bl []byte) {
		if r != responseTypeNewTradeBucket {
			return
		}

		var res newTradeBucketResponse
		err := json.Unmarshal(bl, &res)
		if err != nil {
			log.Println("json unmarshal error:", err)
			return
		}

		fn(res.Data)
	})
	if err != nil {
		return err
	}

	err = subscribe(c.conn, eventTypeNewTradeBucket, pairs)
	if err != nil {
		return err
	}

	return pingForever(ctx, c.conn, c.done)
}

func NewTradeStream(ctx context.Context, apiKey, apiSecret string, pairs []string, fn func(NewTrade)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c := &client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	err := c.connect(ctx, "/ws/trade", func(r responseType, bl []byte) {
		if r != responseTypeNewTrade {
			return
		}

		var res newTradeResponse
		err := json.Unmarshal(bl, &res)
		if err != nil {
			log.Println("json unmarshal error:", err)
			return
		}

		fn(res.Data)
	})
	if err != nil {
		return err
	}

	err = subscribe(c.conn, eventTypeNewTrade, pairs)
	if err != nil {
		return err
	}

	return pingForever(ctx, c.conn, c.done)
}

type AccountStreamOpt func(*accountStreamOpts)

type accountStreamOpts struct {
	onAny               func(responseType string, bl []byte)
	onBalanceUpdate     func(BalanceUpdate)
	onOpenOrdersUpdate  func([]OpenOrderUpdate)
	onOrderStatusUpdate func(OrderStatusUpdate)
}

func OnAny(fn func(responseType string, bl []byte)) AccountStreamOpt {
	return func(o *accountStreamOpts) {
		o.onAny = fn
	}
}

func OnBalanceUpdate(fn func(BalanceUpdate)) AccountStreamOpt {
	return func(o *accountStreamOpts) {
		o.onBalanceUpdate = fn
	}
}

func OnOpenOrdersUpdate(fn func([]OpenOrderUpdate)) AccountStreamOpt {
	return func(o *accountStreamOpts) {
		o.onOpenOrdersUpdate = fn
	}
}

func OnOrderStatusUpdate(fn func(OrderStatusUpdate)) AccountStreamOpt {
	return func(o *accountStreamOpts) {
		o.onOrderStatusUpdate = fn
	}
}

func NewAccountStream(ctx context.Context, apiKey, apiSecret string, opts ...AccountStreamOpt) error {
	opt := accountStreamOpts{
		onAny:               func(responseType string, bl []byte) {},
		onBalanceUpdate:     func(BalanceUpdate) {},
		onOpenOrdersUpdate:  func([]OpenOrderUpdate) {},
		onOrderStatusUpdate: func(OrderStatusUpdate) {},
	}

	for _, o := range opts {
		o(&opt)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c := &client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	err := c.connect(ctx, "/ws/account", func(r responseType, bl []byte) {
		opt.onAny(string(r), bl)

		switch r {
		case responseTypeBalanceUpdate:
			var res balanceUpdateResponse
			err := json.Unmarshal(bl, &res)
			if err != nil {
				log.Println("balanceUpdateResponse unmarshal error:", err)
				return
			}

			opt.onBalanceUpdate(res.Data)
		case responseTypeOpenOrdersUpdate:
			var res openOrdersUpdateResponse
			err := json.Unmarshal(bl, &res)
			if err != nil {
				log.Println("openOrdersUpdateResponse unmarshal error:", err)
				return
			}

			opt.onOpenOrdersUpdate(res.Data)
		case responseTypeOrderStatusUpdate:
			var res orderStatusUpdateResponse
			err := json.Unmarshal(bl, &res)
			if err != nil {
				log.Println("orderStatusUpdateResponse unmarshal error:", err)
				return
			}

			opt.onOrderStatusUpdate(res.Data)
		}
	})
	if err != nil {
		return err
	}

	return pingForever(ctx, c.conn, c.done)
}

func subscribe(conn *websocket.Conn, e eventType, pl []string) error {
	req, err := json.Marshal(request{
		Type: requestTypeSubscribe,
		Subscriptions: []subscription{
			{
				Event: e,
				Pairs: pl,
			},
		},
	})
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, req)
	if err != nil {
		return err
	}

	return nil
}

type client struct {
	apiKey, apiSecret string
	conn              *websocket.Conn
	done              chan struct{}
}

func (c *client) connect(ctx context.Context, path string, fn func(responseType, []byte)) error {
	u := url.URL{Scheme: "wss", Host: host, Path: path}

	t0 := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	signature, err := signRequest(c.apiSecret, http.MethodGet, u.Path, t0)
	if err != nil {
		return err
	}

	h := make(http.Header)
	h.Add("X-VALR-API-KEY", c.apiKey)
	h.Add("X-VALR-TIMESTAMP", t0)
	h.Add("X-VALR-SIGNATURE", signature)

	c.conn, _, err = websocket.DefaultDialer.DialContext(ctx, u.String(), h)
	if err != nil {
		return err
	}

	c.conn.PingHandler()

	c.done = make(chan struct{})

	go func() {
		defer close(c.done)
		for {
			if ctx.Err() != nil {
				log.Println("context error:", err)
				return
			}

			_, bl, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read message error:", err)
				return
			}

			var res response
			err = json.Unmarshal(bl, &res)
			if err != nil {
				log.Println("json unmarshal error:", err)
				return
			}

			fn(res.Type, bl)
		}
	}()

	return nil
}

func pingForever(ctx context.Context, conn *websocket.Conn, done <-chan struct{}) error {
	ping, err := json.Marshal(request{Type: requestTypePing})
	if err != nil {
		return err
	}

	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-done:
		case <-ctx.Done():
			ticker.Stop()
			return conn.Close()
		case <-ticker.C:
			err = conn.WriteMessage(websocket.TextMessage, ping)
			if err != nil {
				return err
			}
		}
	}
}

func signRequest(apiSecret, verb, path, timestamp string) (string, error) {
	h := hmac.New(sha512.New, []byte(apiSecret))
	_, err := h.Write([]byte(timestamp + verb + path))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
