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
	c := &client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	err := c.connect(ctx, func(r responseType, bl []byte) {
		if r != responseTypeMarketSummaryUpdate {
			return
		}

		var res marketSummaryUpdateResponse
		err := json.Unmarshal(bl, &res)
		if err != nil {
			log.Println("unmarshal:", err)
			return
		}

		fn(res.Data)
	})
	if err != nil {
		return err
	}

	return subscribe(c.conn, eventTypeMarketSummaryUpdate, pairs)
}

func NewAggregatedOrderbookUpdatesStream(ctx context.Context, apiKey, apiSecret string, pairs []string, fn func(AggregatedOrderbookUpdate)) error {
	c := &client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	err := c.connect(ctx, func(r responseType, bl []byte) {
		if r != responseTypeMarketSummaryUpdate {
			return
		}

		var res aggregatedOrderbookUpdateResponse
		err := json.Unmarshal(bl, &res)
		if err != nil {
			log.Println("unmarshal:", err)
			return
		}

		fn(res.Data)
	})
	if err != nil {
		return err
	}

	return subscribe(c.conn, eventTypeAggregatedOrderbookUpdate, pairs)
}

func NewTradeBucketStream(ctx context.Context, apiKey, apiSecret string, pairs []string, fn func(NewTradeBucket)) error {
	c := &client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	err := c.connect(ctx, func(r responseType, bl []byte) {
		if r != responseTypeMarketSummaryUpdate {
			return
		}

		var res newTradeBucketResponse
		err := json.Unmarshal(bl, &res)
		if err != nil {
			log.Println("unmarshal:", err)
			return
		}

		fn(res.Data)
	})
	if err != nil {
		return err
	}

	return subscribe(c.conn, eventTypeNewTradeBucket, pairs)
}

func NewTradeStream(ctx context.Context, apiKey, apiSecret string, pairs []string, fn func(NewTrade)) error {
	c := &client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	err := c.connect(ctx, func(r responseType, bl []byte) {
		if r != responseTypeNewTrade {
			return
		}

		var res newTradeResponse
		err := json.Unmarshal(bl, &res)
		if err != nil {
			log.Println("unmarshal:", err)
			return
		}

		fn(res.Data)
	})
	if err != nil {
		return err
	}

	return subscribe(c.conn, eventTypeNewTrade, pairs)
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

func (c *client) connect(ctx context.Context, fn func(responseType, []byte)) error {
	u := url.URL{Scheme: "wss", Host: host, Path: "/ws/trade"}

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

	c.done = make(chan struct{})

	go func() {
		defer close(c.done)
		for {
			_, bl, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			var res response
			err = json.Unmarshal(bl, &res)
			if err != nil {
				log.Println("read:", err)
				return
			}

			fn(res.Type, bl)
		}
	}()

	go pingForever(c.conn, c.done)

	return nil
}

func pingForever(conn *websocket.Conn, done <-chan struct{}) {
	ping, err := json.Marshal(request{Type: requestTypePing})
	if err != nil {
		log.Println("marshal:", err)
		return
	}

	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-done:
			ticker.Stop()
			conn.Close()
			return
		case <-ticker.C:
			err = conn.WriteMessage(websocket.TextMessage, ping)
			if err != nil {
				log.Println("write:", err)
				return
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
