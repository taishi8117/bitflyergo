package bitflyergo

import (
	"encoding/json"
	"time"
)

const (
	// PathGetMarkets is path of '/getmarkets'
	PathGetMarkets = "/getmarkets"

	// PathGetBoard is path of '/getboard'
	PathGetBoard = "/getboard"

	// PathGetTicker is path of '/getticker'
	PathGetTicker = "/getticker"

	// PathGetExecutions is path of '/getexecutions'
	PathGetExecutions = "/getexecutions"

	// PathGetBoardState is path of '/getboardstate'
	PathGetBoardState = "/getboardstate"

	// PathGetHealth is path of '/gethealth'
	PathGetHealth = "/gethealth"
)

// GetMarkets gets market information.
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

// GetTicker gets ticker of specified product_code.
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

// GetExecutions gets executions.
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

// GetBoard gets board of specified product_code.
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

// GetBoardState gets board state of spefied product_code.
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

// GetHealth gets heath of market.
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
