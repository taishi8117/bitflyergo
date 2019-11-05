package bitflyergo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

const (
	baseUrl              = "https://api.bitflyer.com"
	HealthNormal         = "NORMAL"
	HealthBusy           = "BUSY"
	HealthVeryBusy       = "VERY_BUSY"
	HealthSuperBusy      = "SUPER_BUSY"
	HealthNoOrder        = "NO_ORDER"
	HealthStop           = "STOP"
	StateRunning         = "RUNNING"
	StateClosed          = "CLOSED"
	StateStarting        = "STARTING"
	StatePreopen         = "PREOPEN"
	StateCircuitBreak    = "CIRCUIT_BREAK"
	StateAwatingSq       = "AWAITING_SQ"
	StateMatured         = "MATURED"
	ChildOrderTypeLimit  = "LIMIT"
	ChildOrderTypeMarket = "MARKET"
	SideBuy              = "BUY"
	SideSell             = "SELL"
	MinimumOrderbleSize  = 0.01
)

func NewBitflyer(
	apiKey string,
	apiSecret string,
	retryStatus []int,
	retryLimit int,
	retryInterval time.Duration) *Bitflyer {
	return &Bitflyer{
		BaseUrl:       baseUrl,
		ApiVersion:    "1",
		apiKey:        apiKey,
		apiSecret:     apiSecret,
		client:        &http.Client{},
		Debug:         false,
		RetryStatus:   retryStatus,
		RetryLimit:    retryLimit,
		RetryInterval: retryInterval,
	}
}

func (bf *Bitflyer) getUrl(path string) string {
	return bf.BaseUrl + "/v" + bf.ApiVersion + path
}

// APIを実行します。指定されたAPIエラーが発生した際はリトライします。
func (bf *Bitflyer) callApiWithRetry(method string, path string, params map[string]string) ([]byte, error) {
	var res []byte
	var err error

	i := 0
	for true {

		// 認証ヘッダを生成
		headers := bf.getAuthHeaders(method, path, params)

		// 指定されたメソッドでAPIを実行する
		if strings.ToLower(method) == "post" {
			res, err = bf.post(bf.BaseUrl+path, params, headers)
		} else if strings.ToLower(method) == "get" {
			res, err = bf.get(bf.BaseUrl+path, params, headers)
		}

		// エラーが発生していないならループ終了
		if err == nil {
			break
		}

		// エラーが発生して、ステータスがリトライ対象かつ試行回数が上限に達していないのであれば、再度送信を行う
		canRetry := false
		switch e := err.(type) {
		case *ApiError:
			if i < bf.RetryLimit {
				for _, status := range bf.RetryStatus {
					if e.Status == status {
						i += 1
						canRetry = true
						log.Println(e)
						log.Printf("Retry [%v/%v] %v\n", i, bf.RetryLimit, path)
						break
					}
				}
			}
		}

		// 発生したエラーがリトライ対象のエラーでない場合
		if !canRetry {
			return nil, err
		}

		// 再度エラーが発生する可能性が高いため、一定間隔を空けてからリトライを実施する
		time.Sleep(bf.RetryInterval)
	}
	return res, nil
}

func (bf *Bitflyer) get(url string, params map[string]string, headers map[string]string) ([]byte, error) {
	if params != nil {
		url += makeQueryString(params)
	}
	return bf.request("GET", url, headers, nil)
}

func (bf *Bitflyer) post(url string, params map[string]string, headers map[string]string) ([]byte, error) {
	var reader io.Reader
	if params != nil {
		paramsJson, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		reader = strings.NewReader(string(paramsJson))
	}
	return bf.request("POST", url, headers, reader)
}

func (bf *Bitflyer) request(method string, url string, headers map[string]string, reader io.Reader) ([]byte, error) {

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	// add header
	if headers != nil {
		for name, value := range headers {
			req.Header.Set(name, value)
		}
	}

	if bf.Debug {
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Printf("%s", dump)
	}

	// send request
	st := time.Now()
	resp, err := bf.client.Do(req)
	if bf.Debug {
		log.Println(method, url, time.Now().Sub(st))
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// return error if response status is not 200
	if resp.StatusCode != http.StatusOK {
		apiErr := &ApiError{}
		err = json.Unmarshal(body, apiErr)
		if err != nil {
			return nil, err
		}
		return nil, apiErr
	}
	return body, nil
}

func sign(message string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func makeQueryString(params map[string]string) string {
	qs := ""
	if params != nil {
		qs += "?"
		for k, v := range params {
			qs += k + "=" + v + "&"
		}
	}
	return qs[0 : len(qs)-1]
}

func (bf *Bitflyer) getDefaultHeaders() map[string]string {
	headers := map[string]string{}
	headers["Content-Type"] = "application/json"
	//headers["User-Agent"] = ""
	return headers
}

func (bf *Bitflyer) getAuthHeaders(method string, path string, params map[string]string) map[string]string {

	url := path
	body := ""
	if params != nil {
		if strings.ToUpper(method) == "GET" {
			url += makeQueryString(params)
		} else {
			data, err := json.Marshal(params)
			if err != nil {
				log.Fatal(err)
			}
			body = string(data)
		}
	}

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	message := ts + strings.ToUpper(method) + url + body
	sign := sign(message, bf.apiSecret)

	headers := bf.getDefaultHeaders()
	headers["ACCESS-KEY"] = bf.apiKey
	headers["ACCESS-TIMESTAMP"] = ts
	headers["ACCESS-SIGN"] = sign
	return headers
}
