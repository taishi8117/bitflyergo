package bitflyergo_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/mitsutoshi/bitflyergo/pkg/bitflyergo"
)

func TestCreateOHLC1Sec(t *testing.T) {
	candles, err := bitflyergo.CreateOHLC(getExecutions(), 1)
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range candles {
		fmt.Printf("%v, %v\n", i, c)
	}
	if len(candles) != 7 {
		t.Fatalf("Expect: 7, Actual: %v, %s", len(candles), "length of candles is not match.")
	}
	c0 := candles[0]
	t0, _ := time.Parse(time.RFC3339Nano, "2019-03-01T00:00:00.0Z")
	if !(c0.Time == t0 && c0.Open == 100 && c0.High == 100 && c0.Low == 100 && c0.Close == 100 && c0.Volume == 0.1) {
		t.Fatalf("%v\n", c0)
	}
}

func TestCreateOHLC2Sec(t *testing.T) {
	candles, err := bitflyergo.CreateOHLC(getExecutions(), 2)
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range candles {
		fmt.Printf("%v, %v\n", i, c)
	}
	if len(candles) != 5 {
		t.Fatalf("Expect: 7, Actual: %v, %s", len(candles), "length of candles is not match.")
	}
}

func TestCreateOHLC5Sec(t *testing.T) {
	candles, err := bitflyergo.CreateOHLC(getExecutions(), 5)
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range candles {
		fmt.Printf("%v, %v\n", i, c)
	}
	if len(candles) != 3 {
		t.Fatalf("Expect: 2, Actual: %v, %s", len(candles), "length of candles is not match.")
	}
	c0 := candles[0]
	t0, _ := time.Parse(time.RFC3339Nano, "2019-03-01T00:00:00.0Z")
	if !(c0.Time == t0 && c0.Open == 100 && c0.High == 102 && c0.Low == 99 && c0.Close == 101 && c0.Volume == 0.4 && c0.Delay == 250*time.Millisecond) {
		t.Fatalf("%v\n", c0)
	}
	c1 := candles[1]
	t1, _ := time.Parse(time.RFC3339Nano, "2019-03-01T00:00:05.0Z")
	if !(c1.Time == t1 && c1.Open == 110 && c1.High == 120 && c1.Low == 110 && c1.Close == 120 && c1.Volume == 0.11 && c1.Delay == 550*time.Millisecond) {
		t.Fatalf("%v\n", c1)
	}
	c2 := candles[2]
	t2, _ := time.Parse(time.RFC3339Nano, "2019-03-01T00:00:10.0Z")
	if !(c2.Time == t2 && c2.Open == 200 && c2.High == 300 && c2.Low == 200 && c2.Close == 300 && c2.Volume == 1.02 && c2.Delay == 1*time.Second) {
		t.Fatalf("%v\n", c2)
	}
}

func TestCreateOHLC10Sec(t *testing.T) {
	candles, err := bitflyergo.CreateOHLC(getExecutions(), 10)
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range candles {
		fmt.Printf("%v, %v\n", i, c)
	}
	if len(candles) != 2 {
		t.Fatalf("Expect: 2, Actual: %v, %s", len(candles), "length of candles is not match.")
	}
}

func getExecutions() []bitflyergo.Execution {

	execDate1, _ := time.Parse(time.RFC3339Nano, "2019-03-01T00:00:00.0Z")
	execDate2, _ := time.Parse(time.RFC3339Nano, "2019-03-01T00:00:01.0Z")
	execDate3, _ := time.Parse(time.RFC3339Nano, "2019-03-01T00:00:02.0Z")
	execDate4, _ := time.Parse(time.RFC3339Nano, "2019-03-01T00:00:04.0Z")
	execDate5, _ := time.Parse(time.RFC3339Nano, "2019-03-01T00:00:05.0Z")
	execDate6, _ := time.Parse(time.RFC3339Nano, "2019-03-01T00:00:09.9Z")
	execDate7, _ := time.Parse(time.RFC3339Nano, "2019-03-01T00:00:10.0Z")
	execDate8, _ := time.Parse(time.RFC3339Nano, "2019-03-01T00:00:10.999Z")

	e1 := bitflyergo.Execution{Id: 1, ExecDate: execDate1, Price: 100, Size: 0.1, Side: "BUY", BuyChildOrderAcceptanceId: "buy-1", SellChildOrderAcceptanceId: "sell-1", Delay: 100 * time.Millisecond}
	e2 := bitflyergo.Execution{Id: 2, ExecDate: execDate2, Price: 102, Size: 0.1, Side: "BUY", BuyChildOrderAcceptanceId: "buy-2", SellChildOrderAcceptanceId: "sell-2", Delay: 200 * time.Millisecond}
	e3 := bitflyergo.Execution{Id: 3, ExecDate: execDate3, Price: 99, Size: 0.1, Side: "BUY", BuyChildOrderAcceptanceId: "buy-2", SellChildOrderAcceptanceId: "sell-2", Delay: 300 * time.Millisecond}
	e4 := bitflyergo.Execution{Id: 4, ExecDate: execDate4, Price: 101, Size: 0.1, Side: "BUY", BuyChildOrderAcceptanceId: "buy-2", SellChildOrderAcceptanceId: "sell-2", Delay: 400 * time.Millisecond}
	e5 := bitflyergo.Execution{Id: 5, ExecDate: execDate5, Price: 110, Size: 0.1, Side: "BUY", BuyChildOrderAcceptanceId: "buy-2", SellChildOrderAcceptanceId: "sell-2", Delay: 500 * time.Millisecond}
	e6 := bitflyergo.Execution{Id: 6, ExecDate: execDate6, Price: 120, Size: 0.01, Side: "BUY", BuyChildOrderAcceptanceId: "buy-2", SellChildOrderAcceptanceId: "sell-2", Delay: 600 * time.Millisecond}
	e7 := bitflyergo.Execution{Id: 7, ExecDate: execDate7, Price: 200, Size: 0.01, Side: "BUY", BuyChildOrderAcceptanceId: "buy-2", SellChildOrderAcceptanceId: "sell-2", Delay: 1500 * time.Millisecond}
	e8 := bitflyergo.Execution{Id: 8, ExecDate: execDate8, Price: 300, Size: 1.01, Side: "BUY", BuyChildOrderAcceptanceId: "buy-2", SellChildOrderAcceptanceId: "sell-2", Delay: 500 * time.Millisecond}
	executions := []bitflyergo.Execution{e1, e2, e3, e4, e5, e6, e7, e8}
	return executions
}
