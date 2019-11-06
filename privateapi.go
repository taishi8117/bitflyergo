package bitflyergo

import (
	"encoding/json"
	"strconv"
)

var (
	pathGetMyExecutions      = "/me/getexecutions"
	pathGetChildOrders       = "/me/getchildorders"
	pathGetPositions         = "/me/getpositions"
	pathGetCollateral        = "/me/getcollateral"
	pathGetBalance           = "/me/getbalance"
	pathSendChildOrder       = "/me/sendchildorder"
	pathCancelChildOrder     = "/me/cancelchildorder"
	pathCancelAllChildOrders = "/me/cancelallchildorders"
)

func (bf *Bitflyer) GetMyExecutions(params map[string]string) ([]MyExecution, error) {
	res, err := bf.callApiWithRetry("GET", "/v"+bf.ApiVersion+pathGetMyExecutions, params)
	if err != nil {
		return nil, err
	}
	var executions []MyExecution
	err = json.Unmarshal(res, &executions)
	if err != nil {
		return nil, err
	}
	return executions, nil
}

func (bf *Bitflyer) GetChildOrders(params map[string]string) ([]ChildOrder, error) {
	res, err := bf.callApiWithRetry("GET", "/v"+bf.ApiVersion+pathGetChildOrders, params)
	if err != nil {
		return nil, err
	}
	var childOrders []ChildOrder
	err = json.Unmarshal(res, &childOrders)
	if err != nil {
		return nil, err
	}
	return childOrders, nil
}

func (bf *Bitflyer) GetPositions(productCode string) ([]Position, error) {
	params := map[string]string{"product_code": productCode}
	res, err := bf.callApiWithRetry("GET", "/v"+bf.ApiVersion+pathGetPositions, params)
	if err != nil {
		return nil, err
	}
	var positions []Position
	err = json.Unmarshal(res, &positions)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

func (bf *Bitflyer) GetCollateral() (*Collateral, error) {
	res, err := bf.callApiWithRetry("GET", "/v"+bf.ApiVersion+pathGetCollateral, nil)
	if err != nil {
		return nil, err
	}
	var collateral Collateral
	err = json.Unmarshal(res, &collateral)
	if err != nil {
		return nil, err
	}
	return &collateral, nil
}

func (bf *Bitflyer) GetBalance() (*[]Balance, error) {
	res, err := bf.callApiWithRetry("GET", "/v"+bf.ApiVersion+pathGetBalance, nil)
	if err != nil {
		return nil, err
	}
	var balances []Balance
	err = json.Unmarshal(res, &balances)
	if err != nil {
		return nil, err
	}
	return &balances, nil
}

func (bf *Bitflyer) SendChildOrder(productCode string, childOrderType string,
	side string, size float64, params map[string]string) (map[string]string, error) {

	if params == nil {
		params = map[string]string{}
	}

	params["product_code"] = productCode
	params["child_order_type"] = childOrderType
	params["side"] = side
	params["size"] = strconv.FormatFloat(size, 'g', 8, 64)

	res, err := bf.callApiWithRetry("POST", "/v"+bf.ApiVersion+pathSendChildOrder, params)
	if err != nil {
		return nil, err
	}

	var orderResult map[string]string
	err = json.Unmarshal(res, &orderResult)
	if err != nil {
		return nil, err
	}
	return orderResult, nil
}

func (bf *Bitflyer) CancelAllChildOrders(productCode string) error {
	params := map[string]string{
		"product_code": productCode,
	}
	_, err := bf.callApiWithRetry("POST", "/v"+bf.ApiVersion+pathCancelAllChildOrders, params)
	return err
}

func (bf *Bitflyer) CancelChildOrder(productCode string, childOrderAcceptanceId string) error {
	params := map[string]string{
		"product_code":              productCode,
		"child_order_acceptance_id": childOrderAcceptanceId,
	}
	_, err := bf.callApiWithRetry("POST", "/v"+bf.ApiVersion+pathCancelChildOrder, params)
	return err
}
