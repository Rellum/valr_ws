package valr_ws

import (
	"time"

	"github.com/shopspring/decimal"
)

type requestType string

const requestTypePing requestType = "PING"
const requestTypeSubscribe requestType = "SUBSCRIBE"

type request struct {
	Type          requestType    `json:"type"`
	Subscriptions []subscription `json:"subscriptions"`
}

type eventType string

const eventTypeMarketSummaryUpdate eventType = "MARKET_SUMMARY_UPDATE"
const eventTypeAggregatedOrderbookUpdate eventType = "AGGREGATED_ORDERBOOK_UPDATE"
const eventTypeNewTrade eventType = "NEW_TRADE"
const eventTypeNewTradeBucket eventType = "NEW_TRADE_BUCKET"

type subscription struct {
	Event eventType `json:"event"`
	Pairs []string  `json:"pairs"`
}

type responseType string

const responseTypeAuthenticated responseType = "AUTHENTICATED"
const responseTypePong responseType = "PONG"
const responseTypeSubscribed responseType = "SUBSCRIBED"
const responseTypeMarketSummaryUpdate responseType = "MARKET_SUMMARY_UPDATE"
const responseTypeAggregatedOrderbookUpdate responseType = "AGGREGATED_ORDERBOOK_UPDATE"
const responseTypeNewTrade responseType = "NEW_TRADE"
const responseTypeNewTradeBucket responseType = "NEW_TRADE_BUCKET"

type response struct {
	Type responseType `json:"type"`
}

type pongResponse struct {
	Type    responseType `json:"type"`
	Message string       `json:"message"`
}

type marketSummaryUpdateResponse struct {
	Type responseType        `json:"type"`
	Pair string              `json:"currencyPairSymbol"`
	Data MarketSummaryUpdate `json:"data"`
}

type MarketSummaryUpdate struct {
	Pair        string          `json:"currencyPairSymbol"`
	Ask         decimal.Decimal `json:"askPrice"`
	Bid         decimal.Decimal `json:"bidPrice"`
	Last        decimal.Decimal `json:"lastTradedPrice"`
	Close       decimal.Decimal `json:"previousClosePrice"`
	BaseVolume  decimal.Decimal `json:"baseVolume"`
	QuoteVolume decimal.Decimal `json:"quoteVolume"`
	High        decimal.Decimal `json:"highPrice"`
	Low         decimal.Decimal `json:"lowPrice"`
	Created     time.Time       `json:"created"`
	Change      decimal.Decimal `json:"changeFromPrevious"`
}

type aggregatedOrderbookUpdateResponse struct {
	Type responseType              `json:"type"`
	Pair string                    `json:"currencyPairSymbol"`
	Data AggregatedOrderbookUpdate `json:"data"`
}

type AggregatedOrderbookUpdate struct {
	Pair       string           `json:"currencyPairSymbol"`
	Asks       []OrderbookEntry `json:"Asks"`
	Bids       []OrderbookEntry `json:"Bids"`
	LastChange time.Time        `json:"LastChange"`
}

type Side string

const SideBuy Side = "buy"
const SideSell Side = "sell"

type OrderbookEntry struct {
	Price      decimal.Decimal `json:"price"`
	Quantity   decimal.Decimal `json:"quantity"`
	Side       Side            `json:"side"`
	Pair       string          `json:"currencyPair"`
	OrderCount int             `json:"orderCount"`
}

type newTradeResponse struct {
	Type responseType `json:"type"`
	Pair string       `json:"currencyPairSymbol"`
	Data NewTrade     `json:"data"`
}

type NewTrade struct {
	Pair     string          `json:"currencyPair"`
	TradedAt time.Time       `json:"tradedAt"`
	Side     Side            `json:"takerSide"`
	Price    decimal.Decimal `json:"price"`
	Quantity decimal.Decimal `json:"quantity"`
}

type newTradeBucketResponse struct {
	Type responseType   `json:"type"`
	Pair string         `json:"currencyPairSymbol"`
	Data NewTradeBucket `json:"data"`
}

type NewTradeBucket struct {
	Pair   string          `json:"currencyPair"`
	Period int             `json:"bucketPeriodInSeconds"`
	Start  time.Time       `json:"startTime"`
	Open   decimal.Decimal `json:"open"`
	High   decimal.Decimal `json:"high"`
	Low    decimal.Decimal `json:"low"`
	Close  decimal.Decimal `json:"close"`
	Volume decimal.Decimal `json:"volume"`
}
