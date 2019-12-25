package rms

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	SEARCH_ORDER_URL = "https://api.rms.rakuten.co.jp/es/2.0/order/searchOrder/"
	GET_ORDER_URL    = ""
)

type (
	/*** RMSとの通信時に使用 ***/
	JsonTime struct {
		time.Time
	}

	CommonMessageModelResponse struct {
		MessageType string `json:"messageType"`
		MessageCode string `json:"messageCode"`
		Message     string `json:"message"`
	}

	SearchOrderReuquest struct {
		OrderProgressList []int    `json:"orderProgressList"`
		DateType          int      `json:"dateType"`
		StartDatetime     JsonTime `json:"startDatetime"`
		EndDatetime       JsonTime `json:"endDatetime"`
		SearchKeywordType int      `json:"searchKeywordType"` // 0 なし, 1 商品名, 2 商品番号, 3 ひとことメモ, 4 注文者お名前, 5 注文者お名前フリガナ, 6 送付先お名前
		SearchKeyword     string   `json:"searchKeyword"`
	}

	SearchOrderPaginationResponseModel struct {
		TotalRecordsAmount int `json:"totalRecordsAmount"`
		TotalPages         int `json:"totalPages"`
		RequestPage        int `json:"requestPage"`
	}

	SearchOrderResponse struct {
		CommonMessageModelResponseList     []CommonMessageModelResponse `json:"MessageModelList"`
		OrderNumberList                    []string                     `json:"orderNumberList"`
		SearchOrderPaginationResponseModel `json:"PaginationResponseModel"`
	}

	GetOrderRequest struct {
	}

	/*** 内部メソッド ***/
	RMSApi struct {
		serviceSecret string
		licenseKey    string
		authorization string
	}

	SearchOrderCondition struct {
		OrderProgressList []int // 100: 注文確認待ち, 200: 楽天処理中, 300: 発送待ち, 400: 変更確定待ち, 500: 発送済み, 600: 支払手続き中, 700: 支払い手続き済み, 800: キャンセル確定待ち, 900: キャンセル確定
	}
)

func (j JsonTime) format() string {
	return j.Time.Format("2006-01-02T03:04:05") + "+0900"
}

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + j.format() + `"`), nil
}

func (a *RMSApi) Initialize(ss, lk string) {
	a.serviceSecret = ss
	a.licenseKey = lk
	a.authorization = base64.StdEncoding.EncodeToString([]byte(a.serviceSecret + ":" + a.licenseKey))
}

// ToDo: Required以外の検索用パラメータの指定、struct使ってやる。
func (a *RMSApi) SearchOrder(dateType int, startDatetime, endDatetime time.Time, cond *SearchOrderCondition) (*SearchOrderResponse, error) {
	if a.authorization == "" {
		return nil, errors.New("Uninitialized")
	}
	reqBody := SearchOrderReuquest{}
	// For Required
	reqBody.DateType = dateType
	reqBody.StartDatetime = JsonTime{startDatetime}
	reqBody.EndDatetime = JsonTime{endDatetime}

	// For Optional
	if cond != nil {
		if len(cond.OrderProgressList) > 0 {
			reqBody.OrderProgressList = cond.OrderProgressList
		}
	}

	jsonStr, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", SEARCH_ORDER_URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "ESA "+a.authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := new(http.Client)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)
	result := SearchOrderResponse{}
	err = json.Unmarshal(byteArray, &result)
	if err != nil {
		return nil, err
	}
	if len(result.CommonMessageModelResponseList) == 0 {
		return nil, errors.New("Uninitialized")
	}
	return &result, nil
}
