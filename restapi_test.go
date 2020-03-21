package bitflyergo

import (
	"fmt"
	"os"
	"testing"
)

const productCode = "FX_BTC_JPY"

var api *Bitflyer

func TestMain(m *testing.M) {
	println("before all...")

	apiKey := os.Getenv("APIKEY")
	apiSecret := os.Getenv("APISECRET")
	code := 1
	if apiKey != "" && apiSecret != "" {
		fmt.Println("apiKey:", apiKey)
		fmt.Println("apiSecret:", apiSecret)
		api = NewBitflyer(apiKey, apiSecret, []int{-1}, 1, 1)
		code = m.Run()
	} else {
		fmt.Println("Environment variables must be defined. [APIKEY=<api key>, APISECRET=<api secret>]")
	}

	println("after all...")
	os.Exit(code)
}

func TestGetPositions(t *testing.T) {
	positions, err := api.GetPositions(productCode)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("positions:", positions)
}

func TestGetCollateral(t *testing.T) {
	collateral, err := api.GetCollateral()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("collateral:", collateral)
}

//func TestGetMeExecutions(t *testing.T) {
//	count := 10
//	params := map[string]string{
//		"product_code": productCode,
//		"count":        strconv.Itoa(count),
//	}
//	executions, err := api.GetMeExecutions(params)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println("executions:", executions)
//}

func TestGetMeChildOrders(t *testing.T) {
	params := map[string]string{
		"product_code": productCode,
		"count":        "1",
	}
	childOrders, err := api.GetChildOrders(params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("childOrders:", childOrders)
}
