package rms

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func TestSearchOrder_初期化なし(t *testing.T) {
	a := RMSApi{}
	_, err := a.SearchOrder(3, time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, 1), nil)
	if err == nil {
		t.Error("このテストはエラーを発生させるテストですが、エラーは出ませんでした。")
		t.FailNow()
	}
	if err != nil && err.Error() != "Uninitialized" {
		t.Errorf("expected: Uninitialized, actual: %s", err.Error())
		t.FailNow()
	}
}

func TestSearchOrder_認証失敗(t *testing.T) {
	a := RMSApi{}
	a.Initialize("hoge", "fuga")
	_, err := a.SearchOrder(3, time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, 1), nil)
	if err == nil {
		t.Error("このテストはエラーを発生させるテストですが、エラーは出ませんでした。")
		t.FailNow()
	}
	if err != nil && err.Error() != "Uninitialized" {
		t.Errorf("expected: Uninitialized, actual: %s", err.Error())
		t.FailNow()
	}
}

func TestSearchOrder_引数なし(t *testing.T) {
	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))
	r, err := a.SearchOrder(3, time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, 1), nil)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}
	if r.CommonMessageModelResponseList[0].MessageType != "INFO" {
		t.Errorf("Happend error expected: INFO, acctual: %s", r.CommonMessageModelResponseList[0].MessageType)
	}
}

func TestSearchOrder_ステータス指定の検索(t *testing.T) {
	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))
	cond := SearchOrderCondition{}
	cond.OrderProgressList = append(cond.OrderProgressList, 100)
	r, err := a.SearchOrder(3, time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, 1), &cond)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}
	if r.CommonMessageModelResponseList[0].MessageType != "INFO" {
		t.Errorf("Happend error expected: INFO, acctual: %s", r.CommonMessageModelResponseList[0].MessageType)
	}
}

func TestSearchOrder_データ数を指定して検索(t *testing.T) {
	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))
	cond := SearchOrderCondition{}
	cond.RequestRecordsAmount = 2
	cond.RequestPage = 1

	r, err := a.SearchOrder(3, time.Now().AddDate(0, 0, -30), time.Now().AddDate(0, 0, 1), &cond)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}

	if len(r.OrderNumberList) != 2 {
		t.Errorf("Error. epected 2, actual %d", len(r.OrderNumberList))
	}
}

func TestSearchOrder_containsArrayの含まれるテスト(t *testing.T) {
	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))
	cond := SearchOrderCondition{}
	cond.SettlementMethod = 1

	r, err := a.SearchOrder(3, time.Now().AddDate(0, 0, -30), time.Now().AddDate(0, 0, 1), &cond)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}

	rd, err := a.GetOrder(r.OrderNumberList, 3)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}

	for _, o := range rd.OrderModelList {
		if o.GetOrderSettlementModel.SettlementMethod != "代金引換" {
			t.Errorf("Error. epected 代金引換, actual %s", o.GetOrderSettlementModel.SettlementMethod)
		}
	}
}

func TestSearchOrder_containsArrayの含まれないテスト(t *testing.T) {
	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))
	cond := SearchOrderCondition{}
	cond.SettlementMethod = -1

	r, err := a.SearchOrder(3, time.Now().AddDate(0, 0, -30), time.Now().AddDate(0, 0, 1), &cond)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}

	if len(r.OrderNumberList) == 0 {
		t.Errorf("Error. epected >0, actual %d", len(r.OrderNumberList))
	}
}

func TestGetOrder_初期化なし(t *testing.T) {
	a := RMSApi{}
	_, err := a.GetOrder([]string{}, 3)
	if err == nil {
		t.Error("このテストはエラーを発生させるテストですが、エラーは出ませんでした。")
		t.FailNow()
	}
	if err != nil && err.Error() != "Uninitialized" {
		t.Errorf("expected: Uninitialized, actual: %s", err.Error())
		t.FailNow()
	}
}

func TestGetOrder_認証失敗(t *testing.T) {
	a := RMSApi{}
	a.Initialize("hoge", "fuga")
	_, err := a.GetOrder([]string{}, 3)
	if err == nil {
		t.Error("このテストはエラーを発生させるテストですが、エラーは出ませんでした。")
		t.FailNow()
	}
	if err != nil && err.Error() != "Uninitialized" {
		t.Errorf("expected: Uninitialized, actual: %s", err.Error())
		t.FailNow()
	}
}

func TestGetOrder_引数なし(t *testing.T) {
	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))
	r, err := a.GetOrder([]string{}, 3)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}
	if r.GetOrderMessageModelList[0].MessageType != "ERROR" {
		t.Errorf("Happend error expected: ERROR, acctual: %s", r.GetOrderMessageModelList[0].MessageType)
	}
}

func TestGetOrder_注文番号指定(t *testing.T) {
	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))
	r, err := a.GetOrder([]string{os.Getenv("ORDER_NUMBER")}, 3)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}
	t.Logf("%v\n", r)
	if r.GetOrderMessageModelList[0].MessageType != "INFO" {
		t.Errorf("Happend error expected: INFO, acctual: %s", r.GetOrderMessageModelList[0].MessageType)
	}
}

func TestUpdateOrderMemo_データ更新1(t *testing.T) {
	dc := 1
	dd := JsonDate{time.Now()}
	st := 1
	m := "hogefuga"
	o := "hoge"
	mps := "hoge"

	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))

	c := UpdateOrderMemoCondition{}
	c.OrderNumber = os.Getenv("ORDER_NUMBER")
	c.DeliveryClass = &dc
	c.DeliveryDate = &dd
	c.ShippingTerm = &st
	c.Memo = &m
	c.Operator = &o
	c.MailPlugSentence = &mps
	err := a.UpdateOrderMemo(&c)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}

	r, err := a.GetOrder([]string{os.Getenv("ORDER_NUMBER")}, 3)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}
	if r.GetOrderMessageModelList[0].MessageType != "INFO" {
		t.Errorf("Happend error expected: INFO, acctual: %s", r.GetOrderMessageModelList[0].MessageType)
	}
	if dc != *r.OrderModelList[0].GetOrderDeliveryModel.DeliveryClass {
		t.Errorf("Happend error expected: %d, acctual: %d", dc, *r.OrderModelList[0].GetOrderDeliveryModel.DeliveryClass)
	}
	ddStr1, _ := dd.MarshalJSON()
	ddStr2, _ := r.OrderModelList[0].DeliveryDate.MarshalJSON()
	if string(ddStr1) != string(ddStr2) {
		t.Errorf("Happend error expected: %s, acctual: %s", ddStr1, ddStr2)
	}
	if st != *r.OrderModelList[0].ShippingTerm {
		t.Errorf("Happend error expected: %d, acctual: %d", st, *r.OrderModelList[0].ShippingTerm)
	}
	if m != *r.OrderModelList[0].Memo {
		t.Errorf("Happend error expected: %s, acctual: %s", m, *r.OrderModelList[0].Memo)
	}
	if o != *r.OrderModelList[0].Operator {
		t.Errorf("Happend error expected: %s, acctual: %s", o, *r.OrderModelList[0].Operator)
	}
	if mps != *r.OrderModelList[0].MailPlugSentence {
		t.Errorf("Happend error expected: %s, acctual: %s", mps, *r.OrderModelList[0].MailPlugSentence)
	}
}

func TestUpdateOrderMemo_データ更新2(t *testing.T) {
	dc := 0
	dd := JsonDate{time.Now()}
	st := 0
	m := ""
	o := ""
	mps := ""

	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))

	c := UpdateOrderMemoCondition{}
	c.OrderNumber = os.Getenv("ORDER_NUMBER")
	c.DeliveryClass = &dc
	c.DeliveryDate = &dd
	c.ShippingTerm = &st
	c.Memo = &m
	c.Operator = &o
	c.MailPlugSentence = &mps
	err := a.UpdateOrderMemo(&c)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}

	r, err := a.GetOrder([]string{os.Getenv("ORDER_NUMBER")}, 3)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}
	if r.GetOrderMessageModelList[0].MessageType != "INFO" {
		t.Errorf("Happend error expected: INFO, acctual: %s", r.GetOrderMessageModelList[0].MessageType)
	}
	if dc != *r.OrderModelList[0].GetOrderDeliveryModel.DeliveryClass {
		t.Errorf("Happend error expected: %d, acctual: %d", dc, *r.OrderModelList[0].GetOrderDeliveryModel.DeliveryClass)
	}
	ddStr1, _ := dd.MarshalJSON()
	ddStr2, _ := r.OrderModelList[0].DeliveryDate.MarshalJSON()
	if string(ddStr1) != string(ddStr2) {
		t.Errorf("Happend error expected: %s, acctual: %s", ddStr1, ddStr2)
	}
	if st != *r.OrderModelList[0].ShippingTerm {
		t.Errorf("Happend error expected: %d, acctual: %d", st, *r.OrderModelList[0].ShippingTerm)
	}
	if r.OrderModelList[0].Memo != nil {
		t.Errorf("Happend error expected: %s, acctual: %s", m, *r.OrderModelList[0].Memo)
	}
	if r.OrderModelList[0].Operator != nil {
		t.Errorf("Happend error expected: %s, acctual: %s", o, *r.OrderModelList[0].Operator)
	}
	if r.OrderModelList[0].MailPlugSentence != nil {
		t.Errorf("Happend error expected: %s, acctual: %s", mps, *r.OrderModelList[0].MailPlugSentence)
	}
}

func TestUpdateOrderShipping_データ更新1(t *testing.T) {
	sdid, _ := strconv.Atoi(os.Getenv("SHIPPINGDETAILID"))
	dc := "1001"
	sn := "1000"
	sd := JsonDate{time.Now()}
	sdf := 0

	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))

	c := UpdateOrderShippingCondition{}
	c.OrderNumber = os.Getenv("ORDER_NUMBER")
	smCond := UpdateOrderShippingBasketidModelCondition{}
	smCond.BasketID, _ = strconv.Atoi(os.Getenv("BASKETID"))
	ssmCond := UpdateOrderShippingShippingModelCondition{}
	ssmCond.ShippingDetailID = &sdid
	ssmCond.DeliveryCompany = &dc
	ssmCond.ShippingNumber = &sn
	ssmCond.ShippingDate = &sd
	ssmCond.ShippingDeleteFlag = &sdf
	smCond.ShippingModelList = append(smCond.ShippingModelList, ssmCond)
	c.BasketidModelList = append(c.BasketidModelList, smCond)

	err := a.UpdateOrderShipping(&c)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}

	r, err := a.GetOrder([]string{os.Getenv("ORDER_NUMBER")}, 3)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}
	if r.GetOrderMessageModelList[0].MessageType != "INFO" {
		t.Errorf("Happend error expected: INFO, acctual: %s", r.GetOrderMessageModelList[0].MessageType)
	}
	if dc != *r.OrderModelList[0].PackageModelList[0].ShippingModelList[0].DeliveryCompany {
		t.Errorf("Happend error expected: %s, acctual: %s", dc, *r.OrderModelList[0].PackageModelList[0].ShippingModelList[0].DeliveryCompany)
	}
	sdStr1, _ := sd.MarshalJSON()
	sdStr2, _ := r.OrderModelList[0].PackageModelList[0].ShippingModelList[0].ShippingDate.MarshalJSON()
	if string(sdStr1) != string(sdStr2) {
		t.Errorf("Happend error expected: %s, acctual: %s", sdStr1, sdStr2)
	}
	if sn != *r.OrderModelList[0].PackageModelList[0].ShippingModelList[0].ShippingNumber {
		t.Errorf("Happend error expected: %s, acctual: %s", sn, *r.OrderModelList[0].PackageModelList[0].ShippingModelList[0].ShippingNumber)
	}
}

func TestUpdateOrderShipping_データ更新2(t *testing.T) {
	sdid, _ := strconv.Atoi(os.Getenv("SHIPPINGDETAILID"))
	dc := "1000"
	sn := ""
	sd := JsonDate{time.Now()}
	sdf := 0

	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))

	c := UpdateOrderShippingCondition{}
	c.OrderNumber = os.Getenv("ORDER_NUMBER")
	smCond := UpdateOrderShippingBasketidModelCondition{}
	smCond.BasketID, _ = strconv.Atoi(os.Getenv("BASKETID"))
	ssmCond := UpdateOrderShippingShippingModelCondition{}
	ssmCond.ShippingDetailID = &sdid
	ssmCond.DeliveryCompany = &dc
	ssmCond.ShippingNumber = &sn
	ssmCond.ShippingDate = &sd
	ssmCond.ShippingDeleteFlag = &sdf
	smCond.ShippingModelList = append(smCond.ShippingModelList, ssmCond)
	c.BasketidModelList = append(c.BasketidModelList, smCond)

	err := a.UpdateOrderShipping(&c)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}

	r, err := a.GetOrder([]string{os.Getenv("ORDER_NUMBER")}, 3)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}
	if r.GetOrderMessageModelList[0].MessageType != "INFO" {
		t.Errorf("Happend error expected: INFO, acctual: %s", r.GetOrderMessageModelList[0].MessageType)
	}
	if dc != *r.OrderModelList[0].PackageModelList[0].ShippingModelList[0].DeliveryCompany {
		t.Errorf("Happend error expected: %s, acctual: %s", dc, *r.OrderModelList[0].PackageModelList[0].ShippingModelList[0].DeliveryCompany)
	}
	sdStr1, _ := sd.MarshalJSON()
	sdStr2, _ := r.OrderModelList[0].PackageModelList[0].ShippingModelList[0].ShippingDate.MarshalJSON()
	if string(sdStr1) != string(sdStr2) {
		t.Errorf("Happend error expected: %s, acctual: %s", sdStr1, sdStr2)
	}
	if r.OrderModelList[0].PackageModelList[0].ShippingModelList[0].ShippingNumber != nil {
		t.Errorf("Happend error expected: %s, acctual: %s", sn, *r.OrderModelList[0].PackageModelList[0].ShippingModelList[0].ShippingNumber)
	}
}
