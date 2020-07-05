package bitflyergo

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

const productCode = "FX_BTC_JPY"

var api *Bitflyer

func TestMain(m *testing.M) {
	apiKey := os.Getenv("APIKEY")
	apiSecret := os.Getenv("APISECRET")
	code := 1
	if apiKey != "" && apiSecret != "" {
		api = NewBitflyer(apiKey, apiSecret, []int{-1}, 1, 1)
		code = m.Run()
	} else {
		fmt.Println("Environment variables must be defined. [APIKEY=<api key>, APISECRET=<api secret>]")
	}
	os.Exit(code)
}

func TestGetPositions(t *testing.T) {
	_, err := api.GetPositions(productCode)
	if err != nil {
		fmt.Printf("%v\n", reflect.TypeOf(err))
		t.Fatal(err)
	}
}

func TestGetCollateral(t *testing.T) {
	_, err := api.GetCollateral()
	if err != nil {
		t.Fatal(err)
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

func TestGetMeChildOrders(t *testing.T) {
	params := map[string]string{
		"product_code": productCode,
		"count":        "1",
	}
	_, err := api.GetChildOrders(params)
	if err != nil {
		fmt.Printf("%v\n", reflect.TypeOf(err))
		switch e := err.(type) {
		case *ApiError:
			fmt.Printf("APIERROR!!! [%v]", e.Status)
		default:
			fmt.Println("ETC!!!")

		}
		t.Fatal(err)
	}
}
