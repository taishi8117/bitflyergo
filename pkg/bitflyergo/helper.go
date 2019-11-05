package bitflyergo

import (
	"time"
)

// 約定履歴を指定された時間単位のOHLCに変換します。
//
// ローソク足の時刻は、該当の時刻からtimeFrameSecの秒の未来までの情報を示します。
//
// e.g.
//
//   timeFrameSec: 3 -> xx:xx:00.000000 - xx:xx:02.999999 の範囲の約定情報を反映
//   timeFrameSec: 5 -> xx:xx:00.000000 - xx:xx:04.999999 の範囲の約定情報を反映
func CreateOHLC(executions []Execution, timeFrameSec int) ([]OHLC, error) {

	if len(executions) == 0 {
		return nil, nil
	}

	// 戻り値用のローソク足
	var candles []OHLC

	// ローソク足の時間単位
	timeFrame := time.Duration(timeFrameSec) * time.Second

	// 先頭の約定履歴の約定日時を取得
	e := executions[0]

	// 先頭のローソク足の時刻を計算
	t, err := getCandleTime(e, timeFrameSec)
	if err != nil {
		return nil, err
	}

	// 最初の１本目のローソク足の初期状態を作成
	ohlc := &OHLC{
		Time:   t,
		Open:   e.Price,
		High:   e.Price,
		Low:    e.Price,
		Close:  e.Price,
		Volume: 0.0,
		Delay:  time.Duration(0),
	}

	// 次のローソク足の時刻を計算
	nextTime := ohlc.Time.Add(time.Duration(timeFrameSec) * time.Second)

	var delaySec []time.Duration // 遅延時間（秒）

	// 各約定履歴の内容からローソク足を作成する
	for _, e := range executions {

		execDateSec, err := parseExecTime(e)
		if err != nil {
			return nil, err
		}

		// 約定履歴の時刻が、次のローソク足の範囲のものかをチェック
		if execDateSec.After(nextTime) || execDateSec.Equal(nextTime) {

			// 次のローソク足の作成に移行
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
			delaySec = append(delaySec, e.Delay)
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
			delaySec = append(delaySec, e.Delay)
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

// 遅延時間の平均値を計算します。
func meanDelay(delays []time.Duration) time.Duration {

	// このローソク足の遅延時間の平均値をとるため、遅延時間を合算
	sumSec := 0.0
	for _, delay := range delays {
		sumSec += delay.Seconds()
	}

	if sumSec > 0.0 {

		// 遅延時間の平均値（秒）を計算
		delaySecMean := sumSec / float64(len(delays))

		// Durationはintを受け取るため、.xの秒表現では精度が落ちてしまう。そのため、一旦ミリ秒に変換してからセットする
		return time.Duration(delaySecMean*1000) * time.Millisecond
	}
	return time.Duration(0)
}
