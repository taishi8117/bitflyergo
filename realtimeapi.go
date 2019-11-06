package bitflyergo

import (
	//"encoding/json"
	"log"
	url2 "net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	url                  = "ws.lightstream.bitflyer.com"
	channelBoard         = "lightning_board_"
	channelBoardSnapshot = "lightning_board_snapshot_"
	channelExecutions    = "lightning_executions_"
	channelTicker        = "lightning_ticker_"
)

type SubscribeParams struct {
	Channel string `json:"channel"`
}

type JsonRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result"`
}

type WebSocketClient struct {
	Con   *websocket.Conn
	Debug bool
}

func (bf *WebSocketClient) Connect() error {
	url := url2.URL{Scheme: "wss", Host: url, Path: "/json-rpc"}
	con, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return err
	}
	bf.Con = con
	return nil
}

// Tickerチャンネルの購読を開始します。
func (bf *WebSocketClient) SubscribeTicker(symbol string) {
	bf.subscribe(channelTicker + symbol)
}

// 約定履歴チャンネルの購読を開始します。
func (bf *WebSocketClient) SubscribeExecutions(symbol string) {
	bf.subscribe(channelExecutions + symbol)
}

// 板チャンネルの購読を開始します。
func (bf *WebSocketClient) SubscribeBoard(symbol string) {
	bf.subscribe(channelBoard + symbol)
}

// 板（スナップショット）チャンネルの購読を開始します。
func (bf *WebSocketClient) SubscribeBoardSnapshot(symbol string) {
	bf.subscribe(channelBoardSnapshot + symbol)
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
	if err := bf.Con.WriteJSON(&JsonRPC2{
		Version: "2.0",
		Method:  method,
		Params:  &SubscribeParams{channel}}); err != nil {
		return err
	}
	return nil
}

func (bf *WebSocketClient) Receive(
	brdSnpCh chan<- Board,
	brdCh chan<- Board,
	excCh chan<- []Execution,
	tkrCh chan<- Ticker,
	errCh chan<- error) {

	defer close(brdSnpCh)
	defer close(brdCh)
	defer close(excCh)
	defer close(tkrCh)
	defer close(errCh)

	for {

		// websocketでメッセージを受信する
		res := new(JsonRPC2)
		if err := bf.Con.ReadJSON(res); err != nil {
			log.Println("Received error:", err)
			errCh <- err
			return
		}

		// メッセージの種類に応じてチャンネルに送信
		if res.Method == "channelMessage" {

			//start := time.Now()
			p := res.Params.(map[string]interface{})
			ch := p["channel"].(string)

			if strings.HasPrefix(ch, channelExecutions) {

				// 約定履歴
				message := p["message"].([]interface{})
				var executions []Execution
				for _, m := range message {
					e := m.(map[string]interface{})
					execDate, err := time.Parse(time.RFC3339Nano, e["exec_date"].(string))
					if err != nil {
						errCh <- err
					}
					execution := Execution{
						Id:                         int64(e["id"].(float64)),
						ExecDate:                   execDate,
						Price:                      e["price"].(float64),
						Size:                       e["size"].(float64),
						Side:                       e["side"].(string),
						BuyChildOrderAcceptanceId:  e["buy_child_order_acceptance_id"].(string),
						SellChildOrderAcceptanceId: e["sell_child_order_acceptance_id"].(string),
						Delay:                      time.Now().Sub(execDate),
					}
					executions = append(executions, execution)
				}
				excCh <- executions

			} else if strings.HasPrefix(ch, channelBoardSnapshot) {

				// 板（スナップショット）
				brdSnpCh <- newBoard(p["message"].(map[string]interface{}))

			} else if strings.HasPrefix(ch, channelBoard) {

				// 板
				brdCh <- newBoard(p["message"].(map[string]interface{}))

			} else if strings.HasPrefix(ch, channelTicker) {

				// Ticker
				t := p["message"].(interface{}).(map[string]interface{})
				timestamp, err := time.Parse(time.RFC3339Nano, t["timestamp"].(string))
				if err != nil {
					errCh <- err
				}
				ticker := Ticker{
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
				}
				tkrCh <- ticker
			}

			// 処理時間
			//fmt.Println("channel:", ch, "time:", time.Now().Sub(start))
		}
	}
	log.Println("Finished receive websocket.")
}

func newBoard(message map[string]interface{}) Board {

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

	return Board{
		Time:     time.Now(),
		MidPrice: message["mid_price"].(float64),
		Bids:     bids,
		Asks:     asks,
	}
}
