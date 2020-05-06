package rms

import (
	"os"
	"testing"
)

func TestGetShopCalendar_通常取得(t *testing.T) {
	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))
	r, err := a.GetShopCalendar("", -1)
	if err != nil {
		t.Errorf("このテストは正常にデータが取得できることを期待するテストですが、エラーが発生しました。%v", err)
		t.FailNow()
	}
	t.Errorf("%v", r)
	if r.Result == nil {
		t.Error("営業日カレンダーの取得に失敗しました。")
		t.FailNow()
	}

}
