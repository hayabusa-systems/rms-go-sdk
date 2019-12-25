package rms

import (
	"os"
	"testing"
	"time"
)

func TestSearchOrder_引数なし(t *testing.T) {
	a := RMSApi{}
	a.Initialize(os.Getenv("SERVICE_SECRET"), os.Getenv("LICENSE_KEY"))
	r, err := a.SearchOrder(3, time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, 1))
	if err != nil {
		t.Errorf("Happend undefined error: %v", err)
		t.FailNow()
	}
	if r.CommonMessageModelResponseList[0].MessageType != "INFO" {
		t.Errorf("Happend error expected: INFO, acctual: %s", r.CommonMessageModelResponseList[0].MessageType)
	}
}
