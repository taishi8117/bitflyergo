package bitflyergo_test

import (
	"testing"
	"time"

	"github.com/mitsutoshi/bitflyergo"
)

var bf = bitflyergo.NewBitflyer("", "", nil, 0, 1*time.Second)

func TestGetMarket(t *testing.T) {
	markets, err := bf.GetMarkets()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range markets {
		if m.ProductCode == "" {
			t.Fatal("ProductCode is blank.")
		}
		if m.MarketType == "" {
			t.Fatal("MarketType  is blank.")
		}
	}
}

func TestGetTicker(t *testing.T) {
}

func TestGetExecutions(t *testing.T) {
}

func TestGetBoard(t *testing.T) {
}

func TestGetBoardState(t *testing.T) {
	bs, err := bf.GetBoardState(bitflyergo.ProductCodeBtcJpy)
	if err != nil {
		t.Fatal(err)
	}
	if bs.Health == "" {
		t.Fatal("Health is blank.")
	}
	if bs.State == "" {
		t.Fatal("State is blank.")
	}
}

func TestGetHealth(t *testing.T) {
}
