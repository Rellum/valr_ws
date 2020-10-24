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
const responseTypeOrderStatusUpdate responseType = "ORDER_STATUS_UPDATE"
const responseTypeOpenOrdersUpdate responseType = "OPEN_ORDERS_UPDATE"
const responseTypeBalanceUpdate responseType = "BALANCE_UPDATE"

type OrderStatusType string

const OrderStatusTypePlaced OrderStatusType = "Placed"
const OrderStatusTypeFailed OrderStatusType = "Failed"
const OrderStatusTypeCancelled OrderStatusType = "Cancelled"
const OrderStatusTypeFilled OrderStatusType = "Filled"
const OrderStatusTypePartiallyFilled OrderStatusType = "Partially Filled"
const OrderStatusTypeInstantOrderReserveFailed OrderStatusType = "Instant Order Balance Reserve Failed"
const OrderStatusTypeInstantOrderReserved OrderStatusType = "Instant Order Balance Reserved"
const OrderStatusTypeInstantOrderCompleted OrderStatusType = "Instant Order Completed"

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

type OrderType string

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

type orderStatusUpdateResponse struct {
	Type responseType      `json:"type"`
	Data OrderStatusUpdate `json:"data"`
}

type OrderStatusUpdate struct {
	OrderID           string          `json:"orderId"`
	OrderStatusType   OrderStatusType `json:"orderStatusType"`
	CurrencyPair      string          `json:"currencyPair"`
	OriginalPrice     decimal.Decimal `json:"originalPrice"`
	RemainingQuantity decimal.Decimal `json:"remainingQuantity"`
	OriginalQuantity  decimal.Decimal `json:"originalQuantity"`
	OrderSide         Side            `json:"orderSide"`
	OrderType         OrderType       `json:"orderType"`
	FailedReason      string          `json:"failedReason"`
	OrderUpdatedAt    time.Time       `json:"orderUpdatedAt"`
	OrderCreatedAt    time.Time       `json:"orderCreatedAt"`
	CustomerOrderID   string          `json:"customerOrderId"`
}

type openOrdersUpdateResponse struct {
	Type responseType      `json:"type"`
	Data []OpenOrderUpdate `json:"data"`
}

type OpenOrderUpdate struct {
	OrderID          string          `json:"orderId"`
	Side             Side            `json:"side"`
	Quantity         decimal.Decimal `json:"quantity"`
	Price            decimal.Decimal `json:"price"`
	CurrencyPair     string          `json:"currencyPair"`
	CreatedAt        time.Time       `json:"createdAt"`
	OriginalQuantity decimal.Decimal `json:"originalQuantity"`
	FilledPercentage decimal.Decimal `json:"filledPercentage"`
	CustomerOrderID  string          `json:"customerOrderId"`
}

type balanceUpdateResponse struct {
	Type responseType  `json:"type"`
	Data BalanceUpdate `json:"data"`
}

type BalanceUpdate struct {
	Currency  Currency        `json:"currency"`
	Available decimal.Decimal `json:"available"`
	Reserved  decimal.Decimal `json:"reserved"`
	Total     decimal.Decimal `json:"total"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type Currency struct {
	ID                             int    `json:"id"`
	Symbol                         string `json:"symbol"`
	DecimalPlaces                  int    `json:"decimalPlaces"`
	IsActive                       bool   `json:"isActive"`
	ShortName                      string `json:"shortName"`
	LongName                       string `json:"longName"`
	CurrencyDecimalPlaces          int    `json:"currencyDecimalPlaces"`
	SupportedWithdrawDecimalPlaces int    `json:"supportedWithdrawDecimalPlaces"`
}

type CurrencyPair struct {
	ID             int             `json:"id"`
	Symbol         string          `json:"symbol"`
	BaseCurrency   Currency        `json:"baseCurrency"`
	QuoteCurrency  Currency        `json:"quoteCurrency"`
	ShortName      string          `json:"shortName"`
	Exchange       string          `json:"exchange"`
	Active         bool            `json:"active"`
	MinBaseAmount  decimal.Decimal `json:"minBaseAmount"`
	MaxBaseAmount  decimal.Decimal `json:"maxBaseAmount"`
	MinQuoteAmount decimal.Decimal `json:"minQuoteAmount"`
	MaxQuoteAmount decimal.Decimal `json:"maxQuoteAmount"`
}
