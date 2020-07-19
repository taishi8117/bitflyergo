package bitflyergo

import (
	"encoding/hex"
	"fmt"
	"log"
	url2 "net/url"
	"strconv"
	"strings"
	"time"

	"crypto/rand"
	"encoding/json"
	"reflect"

	"github.com/gorilla/websocket"
)

const (
	url                  = "ws.lightstream.bitflyer.com"
	channelBoard         = "lightning_board_"
	channelBoardSnapshot = "lightning_board_snapshot_"
	channelExecutions    = "lightning_executions_"
	channelTicker        = "lightning_ticker_"
	channelChildOrder    = "child_order_events"
	channelParentOrder   = "parent_order_events"
	authJsonRpcId        = 1
)

type subscribeParams struct {
	Channel string `json:"channel"`
}

type authParams struct {
	ApiKey    string `json:"api_key"`
	Timestamp int64  `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}

type jsonRPC2 struct {
	Version string      `json:"version"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
}

type WebSocketClient struct {
	Con   *websocket.Conn
	Debug bool
	Cb    Callback
}

// Callback is the callback functions when receiving data from websocket.
type Callback interface {

	// OnReceiveBoard is the callbck when board is received from websocket.
	OnReceiveBoard(channelName string, board *Board)

	// OnReceiveBoardSnapshot is the callbck when board snapshot is received from websocket.
	OnReceiveBoardSnapshot(channelName string, board *Board)

	// OnReceiveExecutions is the callbck when executions is received from websocket.
	OnReceiveExecutions(channelName string, executions []Execution)

	// OnReceiveTicker is the callbck when ticker is received from websocket.
	OnReceiveTicker(channelName string, ticker *Ticker)

	// OnReceiveChildOrderEvents is the callbck when child order event is received from websocket.
	OnReceiveChildOrderEvents(channelName string, event []ChildOrderEvent)

	// OnReceiveParentOrderEvents is the callbck when board is received from websocket.
	OnReceiveParentOrderEvents(channelName string, event []ParentOrderEvent)

	// OnErrorOccur is the callbck when error is occurred during receiving stream data.
	OnErrorOccur(channelName string, err error)
}

// ChildOrderEvent is type of child order event receiving from websocket.
type ChildOrderEvent struct {
	ProductCode            string         `json:"product_code"`              // product_code
	ChildOrderId           string         `json:"child_order_id"`            // child_order_id
	ChildOrderAcceptanceId string         `json:"child_order_acceptance_id"` // child_order_acceptance_id
	EventDate              EventTime      `json:"event_date"`                // event_date
	EventType              string         `json:"event_type"`                // event_type
	ChildOrderType         string         `json:"child_order_type"`          // child_order_type
	ExpireDate             TimeWithSecond `json:"expire_date"`               // expire_date
	Reason                 string         `json:"reason"`                    // reason
	ExecId                 int            `json:"exec_id"`                   // exec_id
	Side                   string         `json:"side"`                      // side
	Price                  int            `json:"price"`                     // price
	Size                   float64        `json:"size"`                      // size
	Commission             float64        `json:"commission"`                // commission
	Sfd                    float64        `json:"sfd"`                       // sfd
}

type EventTime struct {
	*time.Time
}

// Logger is logger.
var Logger *log.Logger

func logf(format string, v ...interface{}) {
	if Logger == nil {
		log.Printf(format, v...)
		return
	}
	Logger.Printf(format, v...)
}

// UnmarshalJSON converts json data to EventTime.
func (tt *EventTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("\"2006-01-02T15:04:05.9Z\"", string(data))
	*tt = EventTime{&t}
	return err
}

// String converts all fields value to one string.
func (t *ChildOrderEvent) String() string {
	tp := reflect.TypeOf(t)
	return fmt.Sprintf(
		"ChildOrderEvent[%s=%s, %s=%s, %s=%s, %s=%v, %s=%s, %s=%s, %s=%v, %s=%v, %s=%v, %s=%v, %s=%v, %s=%v, %s=%v, %s=%v]",
		tp.Field(0).Tag.Get("json"), t.ProductCode,
		tp.Field(1).Tag.Get("json"), t.ChildOrderId,
		tp.Field(2).Tag.Get("json"), t.ChildOrderAcceptanceId,
		tp.Field(3).Tag.Get("json"), t.EventDate,
		tp.Field(4).Tag.Get("json"), t.EventType,
		tp.Field(5).Tag.Get("json"), t.ChildOrderType,
		tp.Field(6).Tag.Get("json"), t.ExpireDate,
		tp.Field(7).Tag.Get("json"), t.Reason,
		tp.Field(8).Tag.Get("json"), t.ExecId,
		tp.Field(9).Tag.Get("json"), t.Side,
		tp.Field(10).Tag.Get("json"), t.Price,
		tp.Field(11).Tag.Get("json"), t.Size,
		tp.Field(12).Tag.Get("json"), t.Commission,
		tp.Field(13).Tag.Get("json"), t.Sfd)
}

// Event of parent order happened.
type ParentOrderEvent struct {
	ProductCode             string         `json:"product_code"`
	ParentOrderId           string         `json:"parent_order_id"`
	ParentOrderAcceptanceId string         `json:"parent_order_acceptance_id"`
	EventDate               EventTime      `json:"event_date"`
	EventType               string         `json:"event_type"`
	ParentOrderType         string         `json:"parent_order_type"`
	Reason                  string         `json:"reason"`
	ChildOrderType          string         `json:"child_order_type"`
	ParameterIndex          int            `json:"parameter_index"`
	ChildOrderAcceptanceId  string         `json:"child_order_acceptance_id"`
	Side                    string         `json:"side"`
	Price                   int            `json:"price"`
	Size                    float64        `json:"size"`
	ExpireDate              TimeWithSecond `json:"expire_date"`
}

// Connect connects to bitflyer's realtime api server.
func (bf *WebSocketClient) Connect() error {
	url := url2.URL{Scheme: "wss", Host: url, Path: "/json-rpc"}
	con, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return err
	}
	bf.Con = con
	return nil
}

// Authenticate for subscribing private channel.
func (bf *WebSocketClient) Auth(apiKey string, apiSecret string) error {

	// create message
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	nonce, err := randomHex(16)
	if err != nil {
		return err
	}
	message := strconv.FormatInt(timestamp, 10) + nonce

	// signed message
	jsonRpc := &jsonRPC2{
		Version: "2.0",
		Method:  "auth",
		Params: authParams{
			ApiKey:    apiKey,
			Timestamp: timestamp,
			Nonce:     nonce,
			Signature: sign(message, apiSecret),
		},
		Id: authJsonRpcId,
	}

	// send
	if err := bf.Con.WriteJSON(&jsonRpc); err != nil {
		return err
	}
	return nil
}

// SubscribeTicker subscribes ticker.
func (bf *WebSocketClient) SubscribeTicker(symbol string) {
	bf.subscribe(channelTicker + symbol)
}

// SubscribeExecutions subscribes executions.
func (bf *WebSocketClient) SubscribeExecutions(symbol string) {
	bf.subscribe(channelExecutions + symbol)
}

// SubscribeBoard subscribes board.
func (bf *WebSocketClient) SubscribeBoard(symbol string) {
	bf.subscribe(channelBoard + symbol)
}

// SubscribeBoardSnapshot subscribes board snapshot.
func (bf *WebSocketClient) SubscribeBoardSnapshot(symbol string) {
	bf.subscribe(channelBoardSnapshot + symbol)
}

// SubscribeChildOrder subscribes child orders.
func (bf *WebSocketClient) SubscribeChildOrder() {
	bf.subscribe(channelChildOrder)
}

// SubscribeParentOrder subscribes parent orders.
func (bf *WebSocketClient) SubscribeParentOrder() {
	bf.subscribe(channelParentOrder)
}

func (bf *WebSocketClient) UnsubscribeTicker(symbol string) {
	bf.unsubscribe(channelTicker + symbol)
}

func (bf *WebSocketClient) UnsubscribeExecutions(symbol string) {
	bf.unsubscribe(channelExecutions + symbol)
}

func (bf *WebSocketClient) UnsubscribeBoard(symbol string) {
	bf.unsubscribe(channelBoard + symbol)
}

func (bf *WebSocketClient) UnsubscribeBoardSnapshot(symbol string) {
	bf.unsubscribe(channelBoardSnapshot + symbol)
}

func (bf *WebSocketClient) UnsubscribeChildOrder() {
	bf.unsubscribe(channelChildOrder)
}

func (bf *WebSocketClient) UnsubscribeParentOrder() {
	bf.unsubscribe(channelParentOrder)
}

func (bf WebSocketClient) subscribe(channel string) {
	if bf.Debug {
		log.Println("Subscribe " + channel)
	}
	_ = bf.writeJson(channel, "subscribe")
}

func (bf WebSocketClient) unsubscribe(channel string) {
	if bf.Debug {
		log.Println("Unsubscribe " + channel)
	}
	_ = bf.writeJson(channel, "unsubscribe")
}

func (bf WebSocketClient) writeJson(channel string, method string) error {
	if err := bf.Con.WriteJSON(&jsonRPC2{
		Version: "2.0",
		Method:  method,
		Params:  &subscribeParams{channel}}); err != nil {
		return err
	}
	return nil
}

// Receive receives stream data from websocket.
func (bf *WebSocketClient) Receive() {

	for {

		var res map[string]interface{}
		if err := bf.Con.ReadJSON(&res); err != nil {
			logf("Received error: %v\n", err)
			bf.Cb.OnErrorOccur("", err)
			return
		}

		if method, ok := res["method"]; ok {
			if method == "channelMessage" {
				p := res["params"].(map[string]interface{})
				ch := p["channel"].(string)

				if strings.HasPrefix(ch, channelExecutions) {

					receivedTime := time.Now()
					message := p["message"].([]interface{})
					var executions []Execution
					for _, m := range message {
						e := m.(map[string]interface{})
						execDate, err := time.Parse(time.RFC3339Nano, e["exec_date"].(string))
						if err != nil {
							logf("Failed to parse time received from executions channel: %s", e["exec_date"].(string))
							bf.Cb.OnErrorOccur(ch, err)
						}
						executions = append(executions, Execution{
							Id:                         int64(e["id"].(float64)),
							ExecDate:                   execDate,
							Price:                      e["price"].(float64),
							Size:                       e["size"].(float64),
							Side:                       e["side"].(string),
							BuyChildOrderAcceptanceId:  e["buy_child_order_acceptance_id"].(string),
							SellChildOrderAcceptanceId: e["sell_child_order_acceptance_id"].(string),
							ReceivedTime:               receivedTime,
						})
					}

					bf.Cb.OnReceiveExecutions(ch, executions)

				} else if strings.HasPrefix(ch, channelBoardSnapshot) {
					bf.Cb.OnReceiveBoardSnapshot(ch, newBoard(p["message"].(map[string]interface{})))

				} else if strings.HasPrefix(ch, channelBoard) {
					bf.Cb.OnReceiveBoard(ch, newBoard(p["message"].(map[string]interface{})))

				} else if strings.HasPrefix(ch, channelChildOrder) {

					// TODO Improve speed (Don't use Marchal and Unmarchal)
					var events []ChildOrderEvent
					msg := p["message"].(interface{}).([]interface{})
					msgJson, err := json.Marshal(&msg)
					if err != nil {
						bf.Cb.OnErrorOccur(ch, err)
					}
					err = json.Unmarshal(msgJson, &events)
					if err != nil {
						logf("Failed to parse ChildOrderEvent: %v", string(msgJson))
						bf.Cb.OnErrorOccur(ch, err)
					}
					bf.Cb.OnReceiveChildOrderEvents(ch, events)
					//log.Printf("time: %v\n", time.Now().Sub(receivedTime))

				} else if strings.HasPrefix(ch, channelParentOrder) {

					// TODO: need to implement
					var events []ParentOrderEvent
					msg := p["message"].(interface{}).([]interface{})
					msgJson, err := json.Marshal(&msg)
					if err != nil {
						bf.Cb.OnErrorOccur(ch, err)
					}
					err = json.Unmarshal(msgJson, &events)
					if err != nil {
						logf("Failed to parse ParentOrderEvent: %v", string(msgJson))
						bf.Cb.OnErrorOccur(ch, err)
					}
					bf.Cb.OnReceiveParentOrderEvents(ch, events)

				} else if strings.HasPrefix(ch, channelTicker) {

					t := p["message"].(interface{}).(map[string]interface{})
					timestamp, err := time.Parse(time.RFC3339Nano, t["timestamp"].(string))
					if err != nil {
						logf("Failed to parse time received from ticker channel: %s", t["timestamp"].(string))
						//errCh <- err
						bf.Cb.OnErrorOccur(ch, err)
					}
					bf.Cb.OnReceiveTicker(ch, &Ticker{
						ProductCode:     t["product_code"].(string),
						Timestamp:       TickerTime{&timestamp},
						TickId:          int64(t["tick_id"].(float64)),
						BestBid:         t["best_bid"].(float64),
						BestAsk:         t["best_ask"].(float64),
						BestBidSize:     t["best_bid_size"].(float64),
						BestAskSize:     t["best_ask_size"].(float64),
						TotalBidDepth:   t["total_bid_depth"].(float64),
						TotalAskDepth:   t["total_ask_depth"].(float64),
						Ltp:             t["ltp"].(float64),
						Volume:          t["volume"].(float64),
						VolumeByProduct: t["volume_by_product"].(float64),
					})
				}
			}

		} else if id, ok := res["id"]; ok {

			// if res has id and id equals authJsonRpcId, it's a response of request authentication
			if id.(float64) == authJsonRpcId {
				if result, ok := res["result"]; ok {
					if result.(bool) {
						log.Println("Succeeded to authenticate.")
						bf.SubscribeChildOrder()
						// TODO: subscrive parent child order
					} else {
						log.Println("Failed to authenticate.")
					}
				}
			}
		}
	}
	log.Println("Finished receive websocket.")
}

func newBoard(message map[string]interface{}) *Board {

	bidsMessage := message["bids"].([]interface{})
	var bids = make(map[float64]float64, len(bidsMessage))
	for _, bid := range bidsMessage {
		b := bid.(map[string]interface{})
		bids[b["price"].(float64)] = b["size"].(float64)
	}

	asksMessage := message["asks"].([]interface{})
	var asks = make(map[float64]float64, len(asksMessage))
	for _, ask := range asksMessage {
		a := ask.(map[string]interface{})
		asks[a["price"].(float64)] = a["size"].(float64)
	}

	return &Board{
		Time:     time.Now(),
		MidPrice: message["mid_price"].(float64),
		Bids:     bids,
		Asks:     asks,
	}
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
