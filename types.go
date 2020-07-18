package bitflyergo

import (
	"fmt"
	"net/http"
	"time"
)

// TimeWithSecond
type TimeWithSecond struct {
	*time.Time
}

// UnmarshalJSON unmarchals json data.
func (tt *TimeWithSecond) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("\"2006-01-02T15:04:05\"", string(data))
	*tt = TimeWithSecond{&t}
	return err
}

// Bitflyer is bitFlyer api client.
type Bitflyer struct {
	BaseUrl       string        // base url
	ApiVersion    string        // api version
	apiKey        string        // api key
	apiSecret     string        // api secret
	Debug         bool          // if true, debug mode
	RetryLimit    int           // retry limit
	RetryStatus   []int         // status to retry
	RetryInterval time.Duration // retry interval
	client        *http.Client
}

// Execution is one of the execution history
type Execution struct {
	Id                         int64     `json:"id"`                             // id
	ExecDate                   time.Time `json:"exec_date"`                      // exec_date
	Price                      float64   `json:"price"`                          // price
	Size                       float64   `json:"size"`                           // size
	Side                       string    `json:"side"`                           // side
	BuyChildOrderAcceptanceId  string    `json:"buy_child_order_acceptance_id"`  // buy_child_order_acceptance_id
	SellChildOrderAcceptanceId string    `json:"sell_child_order_acceptance_id"` // sell_child_order_acceptance_id
	ReceivedTime               time.Time `json:receivedTime`                     // receivedTime
}

// Returns receiving delayed time
func (e *Execution) Delay() time.Duration {
	return e.ReceivedTime.Sub(e.ExecDate)
}

// Board is board.
type Board struct {
	Time     time.Time           `json:"time"`      // time
	MidPrice float64             `json:"mid_price"` // mid_price
	Bids     map[float64]float64 `json:"bids"`      // bids
	Asks     map[float64]float64 `json:"asks"`      // asks
}

// Market is the return value of '/getmarkets' API.
type Market struct {
	ProductCode string `json:"product_code"` // product_code
	MarketType  string `json:"market_type"`  // market_type
	Alias       string `json:"alias"`        // alias
}

// Ticker  is the return value of '/getticker' API.
type Ticker struct {
	ProductCode     string     `json:"product_code"`      // product_code
	Timestamp       TickerTime `json:"timestamp"`         // timestamp
	TickId          int64      `json:"tick_id"`           // tick_id
	BestBid         float64    `json:"best_bid"`          // best_bid
	BestAsk         float64    `json:"best_ask"`          // best_ask
	BestBidSize     float64    `json:"best_bid_size"`     // best_bid_size
	BestAskSize     float64    `json:"best_ask_size"`     // best_ask_size
	TotalBidDepth   float64    `json:"total_bid_depth"`   // total_bid_depth
	TotalAskDepth   float64    `json:"total_ask_depth"`   // total_ask_depth
	Ltp             float64    `json:"ltp"`               // ltp
	Volume          float64    `json:"volume"`            // volume
	VolumeByProduct float64    `json:"volume_by_product"` // volume_by_product
}

// TickerTime
type TickerTime struct {
	*time.Time
}

// UnmarshalJSON unmarchals json data.
func (tt *TickerTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("2006-01-02T15:04:05.9", string(data))
	*tt = TickerTime{&t}
	return err
}

// BoardState
type BoardState struct {
	Health string             `json:"health"` // health
	State  string             `json:"state"`  // state
	Data   map[string]float64 `json:"data"`   // data
}

// Health
type Health struct {
	Status string `json:"status"` // status
}

// Collateral
type Collateral struct {
	Collateral        float64 `json:"collateral"`         // collateral
	OpenPositionPnl   float64 `json:"open_position_pnl"`  // open_position_pnl
	RequireCollateral float64 `json:"require_collateral"` // require_collateral
	KeepRate          float64 `json:"keep_rate"`          // keep_rate
}

// Balance
type Balance struct {
	CurrencyCode string `json:"currency_code"` // currency_code
	Amount       int64  `json:"amount"`        // amount
	Available    int64  `json:"available"`     // available
}

// ChildOrder
type ChildOrder struct {
	Id                     int64          `json:"id"`                        // id
	ChildOrderId           string         `json:"child_order_id"`            // child_order_id
	ProductCode            string         `json:"product_code"`              // product_code
	Side                   string         `json:"side"`                      // side
	ChildOrderType         string         `json:"child_order_type"`          // child_order_type
	Price                  float64        `json:"price"`                     // price
	AveragePrice           float64        `json:"average_price"`             // average_price
	Size                   float64        `json:"size"`                      // size
	ChildOrderState        string         `json:"child_order_state"`         // child_order_state
	ExpireDate             TimeWithSecond `json:"expire_date"`               // expire_date
	ChildOrderDate         TimeWithSecond `json:"child_order_date"`          // child_order_date
	ChildOrderAcceptanceId string         `json:"child_order_acceptance_id"` // child_order_acceptance_id
	OutstandingSize        float64        `json:"outstanding_size"`          // outstanding_size
	CancelSize             float64        `json:"cancel_size"`               // cancel_size
	ExecutedSize           float64        `json:"executed_size"`             // executed_size
	TotalCommission        float64        `json:"total_commission"`          // total_commission
	Executions             []MyExecution
}

// Position
type Position struct {
	ProductCode         string         `json:"product_code"`          // product_code
	Side                string         `json:"side"`                  // side
	Price               float64        `json:"price"`                 // price
	Size                float64        `json:"size"`                  // size
	Commission          float64        `json:"commission"`            // commission
	SwapPointAccumulate float64        `json:"swap_point_accumulate"` // swap_point_accumulate
	RequireCollateral   float64        `json:"require_collateral"`    // require_collateral
	OpenDate            TimeWithSecond `json:"open_date"`             // open_date
	Leverage            float64        `json:"leverage"`              // leverage
	Pnl                 float64        `json:"pnl"`                   // pnl
	Std                 float64        `json:"sfd"`                   // sfd
}

// ApiError is lightning api error.
type ApiError struct {
	Status       int    `json:"status"`        // status
	ErrorMessage string `json:"error_message"` // error_message
	Data         string `json:"data"`          // data
}

// Error returns error string
func (err *ApiError) Error() string {
	return fmt.Sprintf(
		"Error -> status: %v, error_message: %v, data: %v\n", err.Status, err.ErrorMessage, err.Data)
}

// MyExecution
type MyExecution struct {
	Id                     int64      `json:"id"`                        // id
	ChildOrderId           string     `json:"child_order_id"`            // child_order_id
	Side                   string     `json:"side"`                      // side
	Price                  float64    `json:"price"`                     // price
	Size                   float64    `json:"size"`                      // size
	Commission             float64    `json:"commission"`                // commission
	ExecDate               TickerTime `json:"exec_date"`                 // exec_date
	ChildOrderAcceptanceId string     `json:"child_order_acceptance_id"` // child_order_acceptance_id
}
