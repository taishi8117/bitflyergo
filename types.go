package bitflyergo

import (
	"fmt"
	"net/http"
	"sort"
	"time"
)

type TimeWithSecond struct {
	*time.Time
}

func (tt *TimeWithSecond) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("\"2006-01-02T15:04:05\"", string(data))
	*tt = TimeWithSecond{&t}
	return err
}

type Bitflyer struct {
	BaseUrl       string
	ApiVersion    string
	apiKey        string
	apiSecret     string
	client        *http.Client
	Debug         bool
	RetryLimit    int
	RetryStatus   []int
	RetryInterval time.Duration
}

type Execution struct {
	Id                         int64     `json:"id"`
	ExecDate                   time.Time `json:"exec_date"`
	Price                      float64   `json:"price"`
	Size                       float64   `json:"size"`
	Side                       string    `json:"side"`
	BuyChildOrderAcceptanceId  string    `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceId string    `json:"sell_child_order_acceptance_id"`
	ReceivedTime               time.Time `json:receivedTime`
}

// Returns receiving delayed time
func (e *Execution) Delay() time.Duration {
	return e.ReceivedTime.Sub(e.ExecDate)
}

type Board struct {
	Time     time.Time           `json:"time"`
	MidPrice float64             `json:"mid_price"`
	Bids     map[float64]float64 `json:"bids"`
	Asks     map[float64]float64 `json:"asks"`
}

func (b *Board) TotalAskSize() float64 {
	size := 0.0
	for _, s := range b.Asks {
		size += s
	}
	return size
}

func (b *Board) TotalBidSize() float64 {
	size := 0.0
	for _, s := range b.Bids {
		size += s
	}
	return size
}

func (b *Board) SortAsks() []float64 {
	keys := make([]float64, 0, len(b.Asks))
	for k := range b.Asks {
		keys = append(keys, k)
	}
	sort.Sort(sort.Float64Slice(keys))
	return keys
}

func (b *Board) SortBids() []float64 {
	keys := make([]float64, 0, len(b.Bids))
	for k := range b.Bids {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))
	return keys
}

type Market struct {
	ProductCode string `json:"product_code"`
	Alias       string `json:"alias"`
}

type Ticker struct {
	ProductCode     string     `json:"product_code"`
	Timestamp       TickerTime `json:"timestamp"`
	TickId          int64      `json:"tick_id"`
	BestBid         float64    `json:"best_bid"`
	BestAsk         float64    `json:"best_ask"`
	BestBidSize     float64    `json:"best_bid_size"`
	BestAskSize     float64    `json:"best_ask_size"`
	TotalBidDepth   float64    `json:"total_bid_depth"`
	TotalAskDepth   float64    `json:"total_ask_depth"`
	Ltp             float64    `json:"ltp"`
	Volume          float64    `json:"volume"`
	VolumeByProduct float64    `json:"volume_by_product"`
}

type TickerTime struct {
	*time.Time
}

func (tt *TickerTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("2006-01-02T15:04:05.9", string(data))
	*tt = TickerTime{&t}
	return err
}

type BoardState struct {
	Health string             `json:"health"`
	State  string             `json:"state"`
	Data   map[string]float64 `json:"data"`
}

type Health struct {
	Status string `json:"status"`
}

type Collateral struct {
	Collateral        float64 `json:"collateral"`
	OpenPositionPnl   float64 `json:"open_position_pnl"`
	RequireCollateral float64 `json:"require_collateral"`
	KeepRate          float64 `json:"keep_rate"`
}

type Balance struct {
	CurrencyCode string `json:"currency_code"`
	Amount       int64  `json:"amount"`
	Available    int64  `json:"available"`
}

type ChildOrder struct {
	Id                     int64          `json:"id"`
	ChildOrderId           string         `json:"child_order_id"`
	ProductCode            string         `json:"product_code"`
	Side                   string         `json:"side"`
	ChildOrderType         string         `json:"child_order_type"`
	Price                  float64        `json:"price"`
	AveragePrice           float64        `json:"average_price"`
	Size                   float64        `json:"size"`
	ChildOrderState        string         `json:"child_order_state"`
	ExpireDate             TimeWithSecond `json:"expire_date"`
	ChildOrderDate         TimeWithSecond ` json:"child_order_date"`
	ChildOrderAcceptanceId string         `json:"child_order_acceptance_id"`
	OutstandingSize        float64        `json:"outstanding_size"`
	CancelSize             float64        `json:"cancel_size"`
	ExecutedSize           float64        `json:"executed_size"`
	TotalCommission        float64        `json:"total_commission"`
	Executions             []MyExecution
}

type Position struct {
	ProductCode         string         `json:"product_code"`
	Side                string         `json:"side"`
	Price               float64        `json:"price"`
	Size                float64        `json:"size"`
	Commission          float64        `json:"commission"`
	SwapPointAccumulate float64        `json:"swap_point_accumulate"`
	RequireCollateral   float64        `json:"require_collateral"`
	OpenDate            TimeWithSecond `json:"open_date"`
	Leverage            float64        `json:"leverage"`
	Pnl                 float64        `json:"pnl"`
	Std                 float64        `json:"sfd"`
}

type OHLC struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
	Delay  time.Duration
}

type ApiError struct {
	Status       int    `json:"status"`
	ErrorMessage string `json:"error_message"`
	Data         string `json:"data"`
}

func (err *ApiError) Error() string {
	return fmt.Sprintf(
		"Error -> status: %v, error_message: %v, data: %v\n", err.Status, err.ErrorMessage, err.Data)
}

type MyExecution struct {
	Id                     int64      `json:"id"`
	ChildOrderId           string     `json:"child_order_id"`
	Side                   string     `json:"side"`
	Price                  float64    `json:"price"`
	Size                   float64    `json:"size"`
	Commition              float64    `json:"commision"`
	ExecDate               TickerTime `json:"exec_date"`
	ChildOrderAcceptanceId string     `json:"child_order_acceptance_id"`
}
