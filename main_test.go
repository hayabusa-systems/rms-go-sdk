package rms

import (
	"os"
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

func TestGetOrder_クレジット分割払い(t *testing.T) {
	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))
	r, err := a.GetOrder([]string{os.Getenv("CARD_MULTI_PAY_ORDER_NUMBER")}, 3)
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}
	t.Logf("%v\n", r)
	if r.GetOrderMessageModelList[0].MessageType != "INFO" {
		t.Errorf("Happend error expected: INFO, acctual: %s", r.GetOrderMessageModelList[0].MessageType)
	}
}
