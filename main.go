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
	GET_ORDER_URL    = "https://api.rms.rakuten.co.jp/es/2.0/order/getOrder/"
)

type (
	/*** RMSとの通信時に使用 ***/
	/*** 共通 ***/
	JsonTime struct {
		time.Time
	}

	JsonDate struct {
		time.Time
	}

	/*** Request ***/

	/*** Response ***/
	CommonMessageModelResponse struct {
		MessageType string `json:"messageType"`
		MessageCode string `json:"messageCode"`
		Message     string `json:"message"`
	}

	/*** Request ***/
	/*** searchOrder ***/
	SearchOrderReuquest struct {
		OrderProgressList []int    `json:"orderProgressList"`
		DateType          int      `json:"dateType"`
		StartDatetime     JsonTime `json:"startDatetime"`
		EndDatetime       JsonTime `json:"endDatetime"`
		SearchKeywordType int      `json:"searchKeywordType"` // 0 なし, 1 商品名, 2 商品番号, 3 ひとことメモ, 4 注文者お名前, 5 注文者お名前フリガナ, 6 送付先お名前
		SearchKeyword     string   `json:"searchKeyword"`
	}

	/*** getOrder ***/
	GetOrderRequest struct {
		OrderNumberList []string `json:"orderNumberList"`
		Version         int      `json:"version"`
	}

	/*** Response ***/
	/*** searchOrder ***/
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

	/*** getOrder ***/
	GetOrderMessageModel struct {
		CommonMessageModelResponse
		OrderNumber string `json:"orderNumber"`
	}

	GetOrderOrdererModel struct {
		ZipCode1       string  `json:"zipCode1"`
		ZipCode2       string  `json:"zipCode2"`
		Prefecture     string  `json:"prefecture"`
		City           string  `json:"city"`
		SubAddress     string  `json:"subAddress"`
		FamilyName     string  `json:"familyName"`
		FirstName      string  `json:"firstName"`
		FamilyNameKana *string `json:"familyNameKana"`
		FirstNameKana  *string `json:"firstNameKana"`
		PhoneNumber1   string  `json:"phoneNumber1"`
		PhoneNumber2   string  `json:"phoneNumber2"`
		PhoneNumber3   string  `json:"phoneNumber3"`
		EmailAddress   string  `json:"emailAddress"`
		Sex            string  `json:"sex"`
		BirthYear      int     `json:"birthYear"`
		BirthMonth     int     `json:"birthMonth"`
		BirthDay       int     `json:"birthDay"`
	}

	GetOrderSettlementModel struct {
		SettlementMethod    string  `json:"settlementMethod"`
		RpaySettlementFlag  int     `json:"rpaySettlementFlag"`
		CardName            *string `json:"cardName"`
		CardNumber          *string `json:"cardNumber"`
		CardOwner           *string `json:"cardOwner"`
		CardYm              *string `json:"cardYm"`
		CardPayType         int     `json:"cardPayType"`
		CardInstallmentDesc *int    `json:"cardInstallmentDesc"`
	}

	GetOrderDeliveryModel struct {
		DeliveryName  string `json:"deliveryName"`
		DeliveryClass *int   `json:"deliveryClass"`
	}

	GetOrderPointModel struct {
		UsedPoint int `json:"usedPoint"`
	}

	GetOrderWrappingModel struct {
		Title              int    `json:"title"`
		Name               string `json:"name"`
		Price              *int   `json:"price"`
		IncludeTaxFlag     int    `json:"includeTaxFlag"`
		DeleteWrappingFlag int    `json:"deleteWrappingFlag"`
		TaxRate            int    `json:"taxRate"`
		TaxPrice           int    `json:"taxPrice"`
	}

	GetOrderSenderModel struct {
		ZipCode1           string  `json:"zipCode1"`
		ZipCode2           string  `json:"zipCode2"`
		Prefecture         string  `json:"prefecture"`
		City               string  `json:"city"`
		SubAddress         string  `json:"subAddress"`
		FamilyName         string  `json:"familyName"`
		FirstName          string  `json:"firstName"`
		FamilyNameKana     *string `json:"familyNameKana"`
		FirstNameKana      *string `json:"firstNameKana"`
		PhoneNumber1       string  `json:"phoneNumber1"`
		PhoneNumber2       string  `json:"phoneNumber2"`
		PhoneNumber3       string  `json:"phoneNumber3"`
		IsolatedIslandFlag int     `json:"isolatedIslandFlag"`
	}

	GetOrderItemModel struct {
		ItemDetailID                     int     `json:"itemDetailId"`
		ItemName                         string  `json:"itemName"`
		ItemID                           int     `json:"itemId"`
		ItemNumber                       *string `json:"itemNumber"`
		ManageNumber                     string  `json:"manageNumber"`
		Price                            int     `json:"price"`
		Units                            int     `json:"units"`
		IncludePostageFlag               int     `json:"includePostageFlag"`
		IncludeTaxFlag                   int     `json:"includeTaxFlag"`
		IncludeCashOnDeliveryPostageFlag int     `json:"includeCashOnDeliveryPostageFlag"`
		SelectedChoice                   *string `json:"selectedChoice"`
		PointRate                        int     `json:"pointRate"`
		PointType                        int     `json:"pointType"`
		InventoryType                    int     `json:"inventoryType"`
		DelvdateInfo                     *string `json:"delvdateInfo"`
		RestoreInventoryFlag             int     `json:"restoreInventoryFlag"`
		DealFlag                         int     `json:"dealFlag"`
		DrugFlag                         int     `json:"drugFlag"`
		DeleteItemFlag                   int     `json:"deleteItemFlag"`
		TaxRate                          float64 `json:"TaxRate"`
		PriceTaxIncl                     int     `json:"priceTaxIncl"`
		IsSingleItemShipping             int     `json:"isSingleItemShipping"`
	}

	GetOrderDeliveryCvsModel struct {
		CvsCode         *int    `json:"cvsCode"`
		StoreGenreCode  *string `json:"storeGenreCode"`
		StoreCode       *string `json:"storeCode"`
		StoreName       *string `json:"storeName"`
		StoreZip        *string `json:"storeZip"`
		StorePrefecture *string `json:"storePrefecture"`
		StoreAddress    *string `json:"storeAddress"`
		AreaCode        *string `json:"areaCode"`
		Depo            *string `json:"depo"`
		OpenTime        *string `json:"openTime"`
		CloseTime       *string `json:"closeTime"`
		CvsRemarks      *string `json:"cvsRemrks"`
	}

	GetOrderCouponModel struct {
		CouponCode        string   `json:"couponCode"`
		ItemID            int      `json:"itemId"`
		CouponName        string   `json:"couponName"`
		CouponSummary     string   `json:"couponSummary"`
		CouponCapital     string   `json:"couponCapital"`
		CouponCapitalCode int      `json:"couponCapitalCode"`
		ExpiryDate        JsonDate `json:"expiryDate"`
		CouponPrice       int      `json:"couponPrice"`
		CouponUnit        int      `json:"couponUnit"`
		CouponTotalPrice  int      `json:"couponTotalPrice"`
	}

	GetOrderChangeReasonModel struct {
		ChangeID               int       `json:"changeId"`
		ChangeType             *int      `json:"changeType"`
		ChangeTypeDetail       int       `json:"changeTypeDetail"`
		ChangeReason           *int      `json:"changeReason"`
		ChangeReasonDetail     *int      `json:"changeReasonDetail"`
		ChangeApplyDatetime    *JsonTime `json:"changeApplyDatetime"`
		ChangeFixDatetime      *JsonTime `json:"changeFixDatetime"`
		ChangeCompleteDatetime *JsonTime `json:"changeCmplDatetime"`
	}

	GetOrderTaxSummaryModel struct {
		TaxRate       float64 `json:"taxRate"`
		ReqPrice      int     `json:"reqPrice"`
		ReaPriceTax   int     `json:"reqPriceTax"`
		TotalPrice    int     `json:"totalPrice"`
		PaymentCharge int     `json:"paymentCharge"`
		CouponPrice   int     `json:"couponPrice"`
		Point         int     `json:"point"`
	}

	GetOrderShippingModel struct {
		ShippingDetailID   int       `json:"shippingDetailId"`
		ShippingNumber     *string   `json:"shippingNumber"`
		DeliveryCompany    *string   `json:"deliveryCompany"`
		DeliveyCompanyName *string   `json:"deliveryCompanyName"`
		ShippingDate       *JsonDate `json:"shippingDate"`
	}

	GetOrderPackageModel struct {
		BascketID                  int     `json:"bascketId"`
		PostagePrice               int     `json:"postagePrice"`
		PostageTaxRate             float64 `json:"postageTaxRate"`
		DeliveryPrice              int     `json:"deliveryPrice"`
		DeliveryTaxRate            float64 `json:"deliveryTaxRate"`
		GoodsTax                   int     `json:"goodsTax"`
		GoodsPrice                 int     `json:"goodsPrice"`
		TotalPrice                 int     `json:"totalPrice"`
		Noshi                      *string `json:"noshi"`
		PackageDeleteFlag          int     `json:"packageDeleteFlag"`
		GetOrderSenderModel        `json:"senderModel"`
		ItemModelList              []GetOrderItemModel     `json:"ItemModelList"`
		ShippingModelList          []GetOrderShippingModel `json:"ShippingModelList"`
		GetOrderDeliveryCvsModel   `json:"DeliveryCvsModel"`
		DefaultDeliveryCompanyCode string `json:"defaultDeliveryCompanyCode"`
	}

	GetOrderOrderModel struct {
		OrderNumber                    string    `json:"orderNumber"`
		OrderProgress                  int       `json:"orderProgress"` // 100: 注文確認待ち, 200: 楽天処理中, 300: 発送待ち, 400: 変更確定待ち, 500: 発送済み, 600: 支払い手続き中, 700: 支払い手続き済み, 800: キャンセル確定待ち, 900: キャンセル確定
		SubStatusID                    *int      `json:"subStatusId"`
		SubStatusName                  *string   `json:"subStatusName"`
		OrderDatetime                  JsonTime  `json:"orderDatetime"`
		ShopOrderConfirmDatetime       *JsonTime `json:"shopOrderCfmDatetime"`
		OrderFixDatetime               *JsonTime `json:"orderFixDatetime"`
		ShippingInstDatetime           *JsonTime `json:"shippingInstDatetime"`
		ShippingCompleteReportDatetime *JsonTime `json:"shippingCmplRptDatetime"`
		CancelDueDate                  *JsonDate `json:"cancelDueDate"`
		DeliveryDate                   *JsonDate `json"deliveryDate"`
		ShippingTerm                   *int      `json:"shippingTerm"`
		Remarks                        *string   `json:"remarks"`
		GiftCheckFlag                  int       `json:"giftCheckFlag"`
		SeveralSenderFlag              int       `json:"severalSenderFlag"`
		EqualSenderFlag                int       `json:"equalSenderFlag"`
		IsolatedIslandFlag             int       `json:"isolatedIslandFlag"`
		RakutenMemberFlag              int       `json:"rakutenMemberFlag"`
		CarrieCode                     int       `json:"carrierCode"`
		EmailCarrierCode               int       `json:"emailCarrierCode"`
		OrderType                      int       `json:"orderType"`
		ReserveNumber                  *string   `json:reserveNumber"`
		ReserveDeliveryCount           *int      `json:"reserveDeliveryCount"`
		CautionDisplayType             int       `json:"cautionDispalyType"`
		RakutenConfirmFlag             int       `json:"rakutenConfirmFlag"`
		GoodsPrice                     int       `json:"goodsPrice"`
		GoodsTax                       int       `json:"goodsTax"`
		PostagePrice                   int       `json:"postagePrice"`
		DeliveryPrice                  int       `json:"deliveryPrice"`
		PaymentCharge                  int       `json:"paymentCharge"`
		PaymentChargeTaxRate           float64   `json:"paymentChargeTaxRate"`
		TotalPrice                     int       `json:"totalPrice"`
		RequestPrice                   int       `json:"requestPrice"`
		CouponAllTotalPrice            int       `json:"requestPrice"`
		CouponShopPrice                int       `json:"couponShopPrice"`
		CouponOtherPrice               int       `json:"couponOtherPrice"`
		AdditionalFeeOccurAmountToUser int       `json:"additionalFeeOccurAmountToUser"`
		AdditionalFeeOccurAmountToShop int       `json:"additionalFeeOccurAmountToShop"`
		AsurakuFlag                    int       `json:"asurakuFlag"`
		DrugFlag                       int       `json:"drugFlag"`
		DealFlag                       int       `json:"dealFlag"`
		MembershipType                 int       `json:"membershipType"`
		Memo                           *string   `json:"Memo"`
		Operator                       *string   `json:"operator"`
		MailPlugSentence               *string   `json:"mailPlugSentence"`
		ModifyFlag                     int       `json:"modifyFlag"`
		IsTaxRecalc                    int       `json:"isTaxRecalc"`
		GetOrderOrdererModel           `json:"OrdererModel"`
		GetOrderSettlementModel        `json:"SettlementModel"`
		GetOrderDeliveryModel          `json:"DeliveryModel"`
		GetOrderPointModel             `json:"DeliveryModel"`
		WrappingModel1                 GetOrderWrappingModel       `json:"WrappingModel"`
		WrappingModel2                 GetOrderWrappingModel       `json:"WrappingModel"`
		PackageModelList               []GetOrderPackageModel      `json:"PackageModelList"`
		CouponModelList                []GetOrderCouponModel       `json:"CouponModelList"`
		ChangeReasonModelList          []GetOrderChangeReasonModel `json:"ChangeReasonModelList"`
		TaxSummaryModelList            []GetOrderTaxSummaryModel   `json:"TaxSummaryModelList"`
	}

	GetOrderResponse struct {
		GetOrderMessageModelList []GetOrderMessageModel `json:"MessageModelList"`
		OrderModelList           []GetOrderOrderModel   `json:"OrderModelList"`
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

func (j *JsonTime) UnmarshalJSON(d []byte) error {
	t, err := time.Parse("\"2006-01-02T03:04:05+0900\"", string(d))
	*j = JsonTime{t}
	return err
}

func (j JsonDate) format() string {
	return j.Time.Format("2006-01-02")
}

func (j JsonDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + j.format() + `"`), nil
}

func (j *JsonDate) UnmarshalJSON(d []byte) error {
	t, err := time.Parse("\"2006-01-02\"", string(d))
	*j = JsonDate{t}
	return err
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

// ToDo: Required以外の検索用パラメータの指定、struct使ってやる。
func (a *RMSApi) GetOrder(oList []string, v int) (*GetOrderResponse, error) {
	if a.authorization == "" {
		return nil, errors.New("Uninitialized")
	}
	reqBody := GetOrderRequest{}
	// For Required
	reqBody.OrderNumberList = oList
	reqBody.Version = v

	jsonStr, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", GET_ORDER_URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "ESA "+a.authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := new(http.Client)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)
	result := GetOrderResponse{}
	err = json.Unmarshal(byteArray, &result)
	if err != nil {
		return nil, err
	}
	if len(result.GetOrderMessageModelList) == 0 {
		return nil, errors.New("Uninitialized")
	}
	return &result, nil
}
