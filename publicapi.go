package bitflyergo

import (
	"encoding/json"
	"time"
)

var (
	PathGetMarkets    = "/getmarkets"
	PathGetBoard      = "/getboard"
	PathGetTicker     = "/getticker"
	PathGetExecutions = "/getexecutions"
	PathGetBoardState = "/getboardstate"
	PathGetHealth     = "/gethealth"
)

func (bf *Bitflyer) GetMarkets() ([]Market, error) {
	res, err := bf.get(bf.getUrl(PathGetMarkets), nil, bf.getDefaultHeaders())
	if err != nil {
		return nil, err
	}
	var markets []Market
	err = json.Unmarshal(res, &markets)
	if err != nil {
		return nil, err
	}
	return markets, nil
}

func (bf *Bitflyer) GetTicker(productCode string) (*Ticker, error) {
	params := map[string]string{"product_code": productCode}
	res, err := bf.get(bf.getUrl(PathGetTicker), params, bf.getDefaultHeaders())
	if err != nil {
		return nil, err
	}
	var ticker Ticker
	err = json.Unmarshal(res, &ticker)
	if err != nil {
		return nil, err
	}
	return &ticker, nil
}

func (bf *Bitflyer) GetExecutions(params map[string]string) ([]Execution, error) {
	res, err := bf.get(bf.getUrl(PathGetExecutions), params, bf.getDefaultHeaders())
	if err != nil {
		return nil, err
	}
	var executions []Execution
	err = json.Unmarshal(res, &executions)
	if err != nil {
		return nil, err
	}
	return executions, nil
}

func (bf *Bitflyer) GetBoard(productCode string) (*Board, error) {
	params := map[string]string{"product_code": productCode}
	res, err := bf.get(bf.getUrl(PathGetBoard), params, bf.getDefaultHeaders())
	if err != nil {
		return nil, err
	}
	var b map[string]interface{}
	err = json.Unmarshal(res, &b)
	if err != nil {
		return nil, err
	}

	board := Board{
		Time: time.Now(),
		Asks: map[float64]float64{},
		Bids: map[float64]float64{},
	}

	// 板情報は、{"price":xxx, "size":xxx}というjson形式から、価格をキー、サイズを値とするmapに変換する
	asks := b["asks"].([]interface{})
	for _, ask := range asks {
		m := ask.(map[string]interface{})
		board.Asks[m["price"].(float64)] = m["size"].(float64)
	}
	bids := b["bids"].([]interface{})
	for _, bid := range bids {
		m := bid.(map[string]interface{})
		board.Bids[m["price"].(float64)] = m["size"].(float64)
	}
	return &board, nil
}

func (bf *Bitflyer) GetBoardState(productCode string) (*BoardState, error) {
	params := map[string]string{"product_code": productCode}
	res, err := bf.get(bf.getUrl(PathGetBoardState), params, bf.getDefaultHeaders())
	if err != nil {
		return nil, err
	}
	var boardState BoardState
	err = json.Unmarshal(res, &boardState)
	if err != nil {
		return nil, err
	}
	return &boardState, nil
}

func (bf *Bitflyer) GetHealth() (*Health, error) {
	res, err := bf.get(bf.getUrl(PathGetHealth), nil, bf.getDefaultHeaders())
	if err != nil {
		return nil, err
	}
	var health Health
	err = json.Unmarshal(res, &health)
	if err != nil {
		return nil, err
	}
	return &health, nil
}
