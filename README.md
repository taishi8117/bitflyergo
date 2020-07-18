bitflyergo
===========

[![CircleCI](https://circleci.com/gh/mitsutoshi/bitflyergo.svg?style=svg)](https://circleci.com/gh/mitsutoshi/bitflyergo) [![Go Report Card](https://goreportcard.com/badge/github.com/mitsutoshi/bitflyergo)](https://goreportcard.com/report/github.com/mitsutoshi/bitflyergo)

bitflyergo is golang library to trade cryptocurrency on bitFlyer.

Support following functions.

* REST-API Client
* Realtime API (WebSocket JSON-RPC)
* Utility functions

## How to use

Here are some ways to use the library. If you want to know API specification, See tha public [documents](https://lightning.bitflyer.com/docs).

### Initialize

You need to call `bitflyergo.NewBitflyer` with arguments. First argument is your api key, second that is your api secret. If you don't use private api, you may specify blank.

```go
apiKey := "<Your API Key>"
apiSecret := "<Your API Secret>"
bf := bitflyergo.NewBitflyer(apiKey, apiSecret)
```

### Call Public API

#### /v1/getexecutions

Get the last five execution histories of `FX_BTC_JPY`.

* product_code: FX_BTC_JPY
* count: 5

```go
params := map[string]string{
    "product_code": "FX_BTC_JPY",
    "count":        "5",
}
executions, err := bf.GetExecutions(params)
```

Return value is `[]bitflyergo.Execution`. 

```go
for _, e := range *executions {
    fmt.Println(e)
}
```

```
{800001572 2019-02-09T10:49:55.58 398759 0.01 BUY JRF20190209-XXXXX-XXXXX1 JRF20190209-YYYYYY-YYYYY1}
{800001572 2019-02-09T10:49:55.57 398758 0.01 BUY JRF20190209-XXXXX-XXXXX2 JRF20190209-YYYYYY-YYYYY2}
...
```

#### /v1/getboard

```go
board, err := bf.GetBoard()
```

### Call Private API

#### /v1/me/getexecutions

#### /v1/me/getchildorders

```go
params := map[string]string{
    "": "",
}
childOrders, err := api.GetChildOrders(params)
```

#### /v1/me/getpositions

```go
productCode := "FX_BTC_JPY"
positions, err := api.GetPositions(productCode)
```

#### /v1/me/getcollateral

```go
collateral, err := api.GetCollateral()
```

#### /v1/me/getbalance

```go
balance, err := api.GetBalance()
```

#### /v1/me/sendchildorder

Place the limit order. `SendChildOrder` returns `childOrderAcceptanceId string`. `childOrderAcceptanceId` is when order is accepted ID.

```go
params := map[string]string{
    "price": "400000",
}
childOrderAcceptanceId, err := api.SendChildOrder("FX_BTC_JPY", "LIMIT", "BUY", 0.01, params)
```

Place the market order. market order does't need to specify price of argument.

```go
childOrderAcceptanceId, err := api.SendChildOrder("FX_BTC_JPY", "MARKET", "BUY", 0.01, nil)
```

### Receive streaming data from websocket

bitflyergo provides the APIs to use bitFlyer Lightning Realtime API.

First, you need to implement `Callback` interface's methods.

```go
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
```

Then, Write code for receiving Realtime API data from websocket.

```go
// Create WebSocketClient with Callback interface implement.
ws := WebSocketClient{
    Debug: false,
	Cb:    &YourCallbackImplement{},
}

// connect Realtime API.
err := ws.Connect()
if err != nil {
	log.Fatal(err)
}

// start receiving data. must to use goroutine.
go ws.Receive()

// subscribe channel
ws.SubscribeExecutions("FX_BTC_JPY")

interrupt := make(chan os.Signal, 1)
signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
LOOP:
	for {
		select {
		case _ = <-interrupt:
			break LOOP
		}
	}
```

## How to test

Tests using private api require following environment variables.

```sh
export APIKEY=<value>
export APISECRET=<value>
```
