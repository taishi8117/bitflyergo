bitflyergo
===========

[![CircleCI](https://circleci.com/gh/mitsutoshi/bitflyergo.svg?style=svg)](https://circleci.com/gh/mitsutoshi/bitflyergo)

bitflyergo is golang library to trade cryptocurrency on bitFlyer.

Support following functions.

* REST-API Client
* Realtime API (WebSocket JSON-RPC)
* Utility functions

## How to use

Here are some ways to use the library. If you want to know API specification, See tha public [documents](https://lightning.bitflyer.com/docs).

### Initialize

You need call `bitflyer.NewBitflyer` with arguments. First argument is your api key, second that is your api secret. If you don't use private api, you may specify blank.

```go
apiKey := "<Your API Key>"
apiSecret := "<Your API Secret>"
bf := bitflyer.NewBitflyer(apiKey, apiSecret)
```

### Public API

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

### Private API

#### /v1/me/getcollateral

```go
collateral, err := api.GetCollateral()
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
