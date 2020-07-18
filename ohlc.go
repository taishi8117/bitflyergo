package bitflyergo

import (
	"time"
)

type OHLC struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
	Delay  time.Duration
}

// CreateOHLC converts executions to OHLC.
//
// e.g.
//
//   timeFrameSec: 3 -> range is between xx:xx:00.000000 and xx:xx:02.999999
//   timeFrameSec: 5 -> range is between xx:xx:00.000000 and xx:xx:04.999999
func CreateOHLC(executions []Execution, timeFrameSec int) ([]OHLC, error) {

	if len(executions) == 0 {
		return nil, nil
	}

	var candles []OHLC

	timeFrame := time.Duration(timeFrameSec) * time.Second

	e := executions[0]

	// get time of first candle
	t, err := getCandleTime(e, timeFrameSec)
	if err != nil {
		return nil, err
	}

	// create first candle with initial state
	ohlc := &OHLC{
		Time:   t,
		Open:   e.Price,
		High:   e.Price,
		Low:    e.Price,
		Close:  e.Price,
		Volume: 0.0,
		Delay:  time.Duration(0),
	}

	// get time of next candle
	nextTime := ohlc.Time.Add(time.Duration(timeFrameSec) * time.Second)

	var delaySec []time.Duration

	for _, e := range executions {

		execDateSec, err := parseExecTime(e)
		if err != nil {
			return nil, err
		}

		// 約定履歴の時刻が、次のローソク足の範囲のものかをチェック
		if execDateSec.After(nextTime) || execDateSec.Equal(nextTime) {

			ohlc.Delay = meanDelay(delaySec)
			candles = append(candles, *ohlc)

			// 約定履歴の時刻が、次のローソク足の時刻よりも未来のものであれば（歯抜けがある）、該当の時刻を内包する時刻まで進める
			for !execDateSec.Before(nextTime.Add(timeFrame)) {
				nextTime = nextTime.Add(timeFrame)
			}
			ohlc = &OHLC{
				Time:   nextTime,
				Open:   e.Price,
				High:   e.Price,
				Low:    e.Price,
				Close:  e.Price,
				Volume: e.Size,
				Delay:  time.Duration(0),
			}
			delaySec = []time.Duration{}
			delaySec = append(delaySec, e.Delay())
			nextTime = nextTime.Add(time.Duration(timeFrameSec) * time.Second)

		} else {

			// 約定履歴が、現在のローソク足の時刻の範囲であれば、各種属性の更新を行う

			// 高値、安値の更新
			if e.Price > ohlc.High {
				ohlc.High = e.Price
			} else if e.Price < ohlc.Low {
				ohlc.Low = e.Price
			}

			// 出来高を加算
			ohlc.Volume += e.Size

			// 終値を更新
			ohlc.Close = e.Price

			// 遅延時間を加算
			delaySec = append(delaySec, e.Delay())
		}
	}

	// 最後の１件のローソク足を追加
	ohlc.Delay = meanDelay(delaySec)
	candles = append(candles, *ohlc)
	return candles, nil
}

// 約定履歴の時刻からミリ秒以下の情報を落として、秒までの精度に変換します。
//
// e.g.
//   2019-03-01T00:00:00.999999Z -> 2019-03-01T00:00:00.0Z
//
func parseExecTime(e Execution) (time.Time, error) {
	return e.ExecDate.Truncate(time.Second), nil
}

// ローソク足の時刻をtimeFrameの倍数に合わせる。
//
// e.g.
//
//   2 sec: 0, 2, 4...
//   3 sec: 0, 3, 12...
//   5 sec: 0, 5, 10...
//
func getCandleTime(e Execution, timeFrame int) (time.Time, error) {

	// 約定履歴の時刻を秒までの時刻に変換
	execDateSec, err := parseExecTime(e)
	if err != nil {
		return time.Time{}, err
	}

	// 約定履歴の時刻をtimeFrameで割った余りの秒数を取得
	diff := execDateSec.Second() % timeFrame

	// 余りがあるならその分を引いて、timeFrameの倍数に時刻を合わせる
	if diff > 0 {
		execDateSec = execDateSec.Add(time.Duration(diff*-1) * time.Second)
	}
	return execDateSec, nil
}

func meanDelay(delays []time.Duration) time.Duration {
	sumSec := 0.0
	for _, delay := range delays {
		sumSec += delay.Seconds()
	}
	if sumSec > 0.0 {
		delaySecMean := sumSec / float64(len(delays))
		return time.Duration(delaySecMean*1000) * time.Millisecond
	}
	return time.Duration(0)
}
