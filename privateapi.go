package bitflyergo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var (
	PathGetMyExecutions      = "/me/getexecutions"        // path of api to get own executions
	PathGetChildOrders       = "/me/getchildorders"       // path of api to get own child orders
	PathGetPositions         = "/me/getpositions"         // path of api to get positions
	PathGetCollateral        = "/me/getcollateral"        // path of api to get collateral
	PathGetBalance           = "/me/getbalance"           // path of api to get balance
	PathSendChildOrder       = "/me/sendchildorder"       // path of api to send child order
	PathCancelChildOrder     = "/me/cancelchildorder"     // path of api to cancel child order
	PathCancelAllChildOrders = "/me/cancelallchildorders" // path of api to cancel all child orders
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
		return nil, errors.New(fmt.Sprintf(
			"Sizes less than %v can not be ordered. [%v]\n", MinimumOrderbleSize, size))
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

// CancelChildOrders cancels child orders.
func (bf *Bitflyer) CancelChildOrder(productCode string, childOrderAcceptanceId string) error {
	params := map[string]string{
		"product_code":              productCode,
		"child_order_acceptance_id": childOrderAcceptanceId,
	}
	_, err := bf.callApiWithRetry("POST", "/v"+bf.ApiVersion+PathCancelChildOrder, params)
	return err
}
