package bitflyergo

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

const productCode = "FX_BTC_JPY"

func newBitflyerWithAuth() *Bitflyer {
	apiKey := os.Getenv("APIKEY")
	apiSecret := os.Getenv("APISECRET")
	if apiKey != "" && apiSecret != "" {
		return NewBitflyer(apiKey, apiSecret, []int{-1}, 1, 1)
	} else {
		fmt.Println("[warn] APIKEY and APISECRET must be defined if you want to test private APIs.")
	}
	return nil
}

func TestGetPositions(t *testing.T) {
	api := newBitflyerWithAuth()
	if api != nil {
		_, err := api.GetPositions(productCode)
		if err != nil {
			fmt.Printf("%v\n", reflect.TypeOf(err))
			t.Fatal(err)
		}
	}
}

func TestGetCollateral(t *testing.T) {
	api := newBitflyerWithAuth()
	if api != nil {
		_, err := api.GetCollateral()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetMeChildOrders(t *testing.T) {
	api := newBitflyerWithAuth()
	if api != nil {
		params := map[string]string{
			"product_code": productCode,
			"count":        "1",
		}
		_, err := api.GetChildOrders(params)
		if err != nil {
			switch e := err.(type) {
			case *ApiError:
				if e.Status != -500 {
					t.Fatal(err)
				}
			default:
				t.Fatal(err)
			}
		}
	}
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
