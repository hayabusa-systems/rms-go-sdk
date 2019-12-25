package rms

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

const (
	SEARCH_ORDER_URL = "https://api.rms.rakuten.co.jp/es/2.0/order/searchOrder/"
	GET_ORDER_URL    = ""
)

type (
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

	RMSApi struct {
		serviceSecret string
		licenseKey    string
		authorization string
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
func (a *RMSApi) SearchOrder(dateType int, startDatetime, endDatetime time.Time) (*SearchOrderResponse, error) {
	if a.authorization == "" {
		return nil, errors.New("Uninitialized")
	}
	reqBody := SearchOrderReuquest{}
	// reqBody.OrderProgressList = append(reqBody.OrderProgressList, 100)
	reqBody.DateType = dateType
	reqBody.StartDatetime = JsonTime{startDatetime}
	reqBody.EndDatetime = JsonTime{endDatetime}

	jsonStr, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", SEARCH_ORDER_URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "ESA "+a.authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("%v\n", string(dump))

	client := new(http.Client)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%v\n", string(byteArray))
	result := SearchOrderResponse{}
	err = json.Unmarshal(byteArray, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
