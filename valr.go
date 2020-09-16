package valr_ws

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const host = "api.valr.com"

type Client interface {
	Connect() error
	OnMarketSummaryUpdates(pairs []string, fn func(MarketSummaryUpdate)) error
	OnAggregatedOrderbookUpdates(pairs []string, fn func(update AggregatedOrderbookUpdate)) error
	OnNewTrade(pairs []string, fn func(update NewTrade)) error
	OnNewTradeBucket(pairs []string, fn func(update NewTradeBucket)) error
	io.Closer
}

func New(apiKey, apiSecret string) Client {
	return &client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

type client struct {
	apiKey, apiSecret string
	conn              *websocket.Conn
	msgs              chan msg
}

type msg struct {
	t responseType
	d []byte
}

func (c *client) Connect() error {
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

	c.conn, _, err = websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		return err
	}

	c.msgs = make(chan msg)

	go func() {
		defer close(c.msgs)
		for {
			_, bl, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Printf("recv: %s", bl)

			var res response
			err = json.Unmarshal(bl, &res)
			if err != nil {
				log.Println("read:", err)
				return
			}

			c.msgs <- msg{t: res.Type, d: bl}
		}
	}()

	ping, err := json.Marshal(request{Type: requestTypePing})
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)

		for {
			select {
			case <-c.msgs:
				ticker.Stop()
				return
			case <-ticker.C:
				err = c.conn.WriteMessage(websocket.TextMessage, ping)
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		}
	}()

	return nil
}

func signRequest(apiSecret, verb, path, timestamp string) (string, error) {
	h := hmac.New(sha512.New, []byte(apiSecret))
	_, err := h.Write([]byte(timestamp + verb + path))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (c *client) OnMarketSummaryUpdates(pairs []string, fn func(MarketSummaryUpdate)) error {
	req, err := json.Marshal(request{
		Type: requestTypeSubscribe,
		Subscriptions: []subscription{
			{
				Event: eventTypeMarketSummaryUpdate,
				Pairs: pairs,
			},
		},
	})
	if err != nil {
		return err
	}

	err = c.conn.WriteMessage(websocket.TextMessage, req)
	if err != nil {
		return err
	}

	go func() {
		for msg := range c.msgs {
			if msg.t != responseTypeMarketSummaryUpdate {
				continue
			}

			var res marketSummaryUpdateResponse
			err = json.Unmarshal(msg.d, &res)
			if err != nil {
				log.Println("unmarshal:", err)
				return
			}

			fn(res.Data)
		}
	}()

	return nil
}

func (c *client) OnAggregatedOrderbookUpdates(pairs []string, fn func(update AggregatedOrderbookUpdate)) error {
	req, err := json.Marshal(request{
		Type: requestTypeSubscribe,
		Subscriptions: []subscription{
			{
				Event: eventTypeAggregatedOrderbookUpdate,
				Pairs: pairs,
			},
		},
	})
	if err != nil {
		return err
	}

	err = c.conn.WriteMessage(websocket.TextMessage, req)
	if err != nil {
		return err
	}

	go func() {
		for msg := range c.msgs {
			if msg.t != responseTypeAggregatedOrderbookUpdate {
				continue
			}

			var res aggregatedOrderbookUpdateResponse
			err = json.Unmarshal(msg.d, &res)
			if err != nil {
				log.Println("unmarshal:", err)
				return
			}

			fn(res.Data)
		}
	}()

	return nil
}

func (c *client) OnNewTrade(pairs []string, fn func(update NewTrade)) error {
	req, err := json.Marshal(request{
		Type: requestTypeSubscribe,
		Subscriptions: []subscription{
			{
				Event: eventTypeNewTrade,
				Pairs: pairs,
			},
		},
	})
	if err != nil {
		return err
	}

	err = c.conn.WriteMessage(websocket.TextMessage, req)
	if err != nil {
		return err
	}

	go func() {
		for msg := range c.msgs {
			if msg.t != responseTypeNewTrade {
				continue
			}

			var res newTradeResponse
			err = json.Unmarshal(msg.d, &res)
			if err != nil {
				log.Println("unmarshal:", err)
				return
			}

			fn(res.Data)
		}
	}()

	return nil
}

func (c *client) OnNewTradeBucket(pairs []string, fn func(update NewTradeBucket)) error {
	req, err := json.Marshal(request{
		Type: requestTypeSubscribe,
		Subscriptions: []subscription{
			{
				Event: eventTypeNewTradeBucket,
				Pairs: pairs,
			},
		},
	})
	if err != nil {
		return err
	}

	err = c.conn.WriteMessage(websocket.TextMessage, req)
	if err != nil {
		return err
	}

	go func() {
		for msg := range c.msgs {
			if msg.t != responseTypeNewTradeBucket {
				continue
			}

			var res newTradeBucketResponse
			err = json.Unmarshal(msg.d, &res)
			if err != nil {
				log.Println("unmarshal:", err)
				return
			}

			fn(res.Data)
		}
	}()

	return nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
