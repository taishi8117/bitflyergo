package bitflyergo

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (

	// PathGetMyExecutions is path of api to get own executions
	PathGetMyExecutions = "/me/getexecutions"

	// PathGetChildOrders is path of api to get own child orders
	PathGetChildOrders = "/me/getchildorders"

	// PathGetPositions is path of api to get positions
	PathGetPositions = "/me/getpositions"

	// PathGetCollateral is path of api to get collateral
	PathGetCollateral = "/me/getcollateral"

	// PathGetBalance is path of api to get balance
	PathGetBalance = "/me/getbalance"

	// PathSendChildOrder is path of api to send child order
	PathSendChildOrder = "/me/sendchildorder"

	// PathCancelChildOrder is path of api to cancel child order
	PathCancelChildOrder = "/me/cancelchildorder"

	// PathCancelAllChildOrders is path of api to cancel all child orders
	PathCancelAllChildOrders = "/me/cancelallchildorders"
)

// GetMyExecutions gets own executions.
func (bf *Bitflyer) GetMyExecutions(params map[string]string) ([]MyExecution, error) {
	res, err := bf.callApiWithRetry("GET", "/v"+bf.ApiVersion+PathGetMyExecutions, params)
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

// GetChildOrders gets own child orders.
//
// Required parameters
// - product_code
func (bf *Bitflyer) GetChildOrders(params map[string]string) ([]ChildOrder, error) {
	res, err := bf.callApiWithRetry("GET", "/v"+bf.ApiVersion+PathGetChildOrders, params)
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

// GetPositions gets positions.
func (bf *Bitflyer) GetPositions(productCode string) ([]Position, error) {
	params := map[string]string{"product_code": productCode}
	res, err := bf.callApiWithRetry("GET", "/v"+bf.ApiVersion+PathGetPositions, params)
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

// GetCollateral gets collateral.
func (bf *Bitflyer) GetCollateral() (*Collateral, error) {
	res, err := bf.callApiWithRetry("GET", "/v"+bf.ApiVersion+PathGetCollateral, nil)
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

// GetBalance gets balance.
func (bf *Bitflyer) GetBalance() (*[]Balance, error) {
	res, err := bf.callApiWithRetry("GET", "/v"+bf.ApiVersion+PathGetBalance, nil)
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

// SendChildOrder send child order.
func (bf *Bitflyer) SendChildOrder(productCode string, childOrderType string,
	side string, size float64, params map[string]string) (map[string]string, error) {

	if size < MinimumOrderbleSize {
		return nil, fmt.Errorf(
			"Sizes less than %v can not be ordered. [%v]\n", MinimumOrderbleSize, size)
	}

	if params == nil {
		params = map[string]string{}
	}

	params["product_code"] = productCode
	params["child_order_type"] = childOrderType
	params["side"] = side
	params["size"] = strconv.FormatFloat(size, 'g', 8, 64)

	res, err := bf.callApiWithRetry("POST", "/v"+bf.ApiVersion+PathSendChildOrder, params)
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

// CancelAllChildOrders cancels all child orders.
func (bf *Bitflyer) CancelAllChildOrders(productCode string) error {
	params := map[string]string{
		"product_code": productCode,
	}
	_, err := bf.callApiWithRetry("POST", "/v"+bf.ApiVersion+PathCancelAllChildOrders, params)
	return err
}

// CancelChildOrder cancels child orders.
func (bf *Bitflyer) CancelChildOrder(productCode string, childOrderAcceptanceId string) error {
	params := map[string]string{
		"product_code":              productCode,
		"child_order_acceptance_id": childOrderAcceptanceId,
	}
	_, err := bf.callApiWithRetry("POST", "/v"+bf.ApiVersion+PathCancelChildOrder, params)
	return err
}
