package bitflyergo

// Extract close price from ohlc.
func Closes(ohlc []OHLC) []float64 {
	var closes []float64
	for _, o := range ohlc {
		closes = append(closes, o.Close)
	}
	return closes
}
