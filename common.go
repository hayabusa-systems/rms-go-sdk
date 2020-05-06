package rms

import "time"

type (
	// JsonTime はRMS WEB SERVICEで使用する日時をGolangで取り扱えるようにするためのラッパークラスです。
	// 表示形式はYYYY-MM-DDThh:mm:ss+0900です。
	JsonTime struct {
		time.Time
	}

	// JsonDate はRMS WEB SERVICEで使用する日付をGolangで取り扱えるようにするためのラッパークラスです。
	// 表示形式はYYYY-MM-DDです。
	JsonDate struct {
		time.Time
	}
)

func (j JsonTime) format() string {
	return j.Time.Format("2006-01-02T15:04:05") + "+0900"
}

// MarshalJSON は値をJSONに変換する際のフォーマット方法を指定します。
func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + j.format() + `"`), nil
}

// UnmarshalJSON はJSONの値をGoの型に変換する際の変換方法を指定します。
func (j *JsonTime) UnmarshalJSON(d []byte) error {
	t, err := time.Parse("\"2006-01-02T15:04:05Z0700\"", string(d))
	*j = JsonTime{t}
	return err
}

// Value は格納されているtime.Timeの値を返却します。
func (j *JsonTime) Value() time.Time {
	return j.Time
}

func (j JsonDate) format() string {
	return j.Time.Format("2006-01-02")
}

// MarshalJSON は値をJSONに変換する際のフォーマット方法を指定します。
func (j JsonDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + j.format() + `"`), nil
}

// UnmarshalJSON はJSONの値をGoの型に変換する際の変換方法を指定します。
func (j *JsonDate) UnmarshalJSON(d []byte) error {
	t, err := time.Parse("\"2006-01-02\"", string(d))
	*j = JsonDate{t}
	return err
}

// Value は格納されているtime.Timeの値を返却します。
func (j *JsonDate) Value() time.Time {
	return j.Time
}
