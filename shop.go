package rms

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// RMS WEB SERVICEの店舗APIの営業日カレンダー設定・長期休暇の告知情報取得用のエンドポイントです。
	SHOP_CALENDAR_URL = "https://api.rms.rakuten.co.jp/es/1.0/shop/shopCalendar"
)

type (
	// ShopBizApiResponse はshopAPIの営業日カレンダー設定・長期休暇の告知情報を取得するAPIの戻り地が格納される構造体です。
	ShopBizApiResponse struct {
		// ResultCode は結果コードです。最大4バイトのデータが格納されます。
		ResultCode string `xml:"resultCode"`
		// ResultMessageList はメッセージ一覧です。
		ResultMessageList ResultMessageList `xml:"resultMessageList"`
		// Result は取得データの本体です。
		Result *ShopCalendarBizModel `xml:"shopCalendarBizModel,omitempty"`
	}

	// ResultMessageList はメッセージを格納する構造体です。不要なのですが、xmlのタグを付与するために存在しています。
	ResultMessageList struct {
		List []ResultMessage `xml:"resultMessage"`
	}

	// ResultMessage はエラーメッセージが格納される構造体です。
	ResultMessage struct {
		// Code は結果コードです。CodeとMessageは対になります。以下の値が格納されます。
		// Code Message    Description
		// N000 Succeeded. 正常に処理が行われました。
		// S900 Under maintenance. メンテナンス中でサービスを提供できません。
		// S901 The Service is currently in read-only mode. 読み取り専用モードで運用中のため、更新機能は利用できません。
		// C001 Request parameter is invalid.  パラメータのフォーマットバリデーションエラー。 Responce の fieldId にパラメータ名が入ります。
		// C005 Requested MediaType is not supported. Accept や ContentType にAPIが対応していない MediaType を指定された場合のコードです。
		// C006 Update data is invalid. 更新データの UniqueKey に問題がある場合などのコードです。 Note: UniqueKey とは更新したいデータを特定する一意のキーです。例：layoutCommonId 等
		// C007 Request Model is invalid. 更新データ Model が不正な場合のコードです。
		// C008 Requested data is not found. 指定された Resource のデータが存在しない場合のコードです。
		// C012 Number of table elements exceeds limit. レコード数が許可されている最大数を超えている場合のコードです。
		// C013 Validation Error. 更新データに不正な値が含まれてる場合のエラーです。 この場合、バリデーションエラーの種類に応じた検証エラーコードが、レスポンスデータに含まれます。 詳細については、validationErrorCode List をご覧ください。
		// C998 Multiple errors occurred. 複数のエラーが発生し、かつすべてのエラーがクライアント要因(コードC~)の場合のコードです。
		// E101 Failed to insert data. データの登録に失敗した場合のコードです。
		// E102 Failed to update data. データの更新に失敗した場合のコードです。
		// E998 Multiple errors occurred. 複数のエラーが発生した場合のコードです。
		// E999 Unknown Execution error. サーバ内でエラーが起こった場合のコードです。 クライアント側で対応不能な場合このコードが返されます
		Code string `xml:"code"`
		// Message はメッセージです。詳しくはCodeの項を参照して下さい。
		Message string `xml:"message"`
		// FieldId はフィールド名です。ここには問題が発生したフィールド名が入ります。パラメータのフォーマットエラーなど特定の項目に関連してエラーが発生した場合にのみ設定されます。
		FieldId string `xml:"fieldId"`
	}

	// ShopCalendarBizModel は店舗カレンダーを格納する構造体です。不要なのですが、xmlのタグを作るためだけに存在しています。
	ShopCalendarBizModel struct {
		Calendar ShopCalendar `xml:"shopCalendar"`
	}

	// ShopCalendar は休業日を格納する構造体です。
	ShopCalendar struct {
		// BusinessHoliday は休業日が格納されます。
		BusinessHoliday CalendarEvent `xml:"businessHoliday"`
		// ShippingHoliday は受注・お問い合わせ業務のみの営業日が格納されます。
		ShippingHoliday CalendarEvent `xml:"shippingHoliday"`
		// ShippingOnly は発送業務のみの営業日が格納されます。
		ShippingOnly CalendarEvent `xml:"shippingOnly"`
		// ShopHoliday は長期休暇の告知が格納されます。
		ShopHoliday ShopHoliday `xml:"shopHoliday"`
	}

	// CalendarEvent は休業日情報を格納する構造体です。
	CalendarEvent struct {
		// RegularSchedule は定期的な休日(曜日)が格納されます。
		RegularSchedule []string `xml:"regularSchedule,omitempty"`
		// EventDateStrs は日付が格納されます。日付はYYYYMMDDの形式です。
		EventDateStrs []string `xml:"eventDates,omitempty"`
		// EventDate はtime.Time型の日付が格納されます。EventDateStrsをxmlから構造体に戻すときにデータを格納します。
		EventDates []*time.Time
	}

	// ShopHoliday は長期休暇の告知を格納する構造体です。
	ShopHoliday struct {
		// Title はWeb告知用タイトルです。382バイトが最大です。
		Title string `xml:"title,omitempty"`
		// StimestampYmd はWeb告知用表示期間の開始日です。YYYY-MM-DDThh:mm:ss+09:00の形式です。
		StimestampYmd string `xml:"stimestampYmd,omitempty"`
		// Stimestamp はWeb告知用表示期間の開始日です。StimestampYmdをxmlから構造体に戻すときデータを格納します。
		Stimestamp *time.Time
		// EtimestampYmd はWeb告知用表示期間の終了日です。YYYY-MM-DDThh:mm:ss+09:00の形式です。
		EtimestampYmd string `xml:"etimestampYmd,omitempty"`
		// Etimestamp はWeb告知用表示期間の開始日です。EtimestampYmdをxmlから構造体に戻すときデータを格納します。
		Etimestamp *time.Time
		// MailMessageはメール告知用メッセージです。3072バイトが最大です。
		MailMessage string `xml:"mailMessage,omitempty"`
		// StimestampMailYmd はメール告知用表示期間の開始日です。YYYY-MM-DDThh:mm:ss+09:00の形式です。
		StimestampMailYmd string `xml:"stimestampYmd,omitempty"`
		// StimestampMail はWeb告知用表示期間の開始日です。StimestampMailYmdをxmlから構造体に戻すときデータを格納します。
		StimestampMail *time.Time
		// EtimestampMailYmd はメール告知用表示期間の終了日です。YYYY-MM-DDThh:mm:ss+09:00の形式です。
		EtimestampMailYmd string `xml:"etimestampYmd,omitempty"`
		// EtimestampMail はWeb告知用表示期間の開始日です。EtimestampMailYmdをxmlから構造体に戻すときデータを格納します。
		EtimestampMail *time.Time
		// MessageはWeb告知用メッセージです。3072バイトが最大です。
		Message string `xml:"message,omitempty"`
	}
)

func (e *CalendarEvent) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	temp := CalendarEvent{}
	if err := d.DecodeElement(&temp, &start); err != nil {
		return err
	}
	*e = temp
	for _, v := range e.EventDateStrs {
		tmpTime, _ := time.Parse("20060102", v)
		e.EventDates = append(e.EventDates, &tmpTime)
	}
	return nil
}

func (h *ShopHoliday) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	temp := ShopHoliday{}
	if err := d.DecodeElement(temp, &start); err != nil {
		return err
	}
	*h = temp
	if h.StimestampYmd != "" {
		*h.Stimestamp, _ = time.Parse("2006-01-02T15:04:05+0900", h.StimestampYmd)
	}
	if h.EtimestampYmd != "" {
		*h.Etimestamp, _ = time.Parse("2006-01-02T15:04:05+0900", h.EtimestampYmd)
	}
	if h.StimestampMailYmd != "" {
		*h.StimestampMail, _ = time.Parse("2006-01-02T15:04:05+0900", h.StimestampMailYmd)
	}
	if h.EtimestampMailYmd != "" {
		*h.EtimestampMail, _ = time.Parse("2006-01-02T15:04:05+0900", h.EtimestampMailYmd)
	}
	return nil
}

// GetShopCalendar はRMSから営業日カレンダー・長期休暇の告知を取得します。fromDate は開始年月日で、YYYY-MM-DDの形式で渡します。指定されない場合は、現在年月日以降の情報を取得します。period は取得する期間です。1~180まで指定することができます。それ以外の場合は90日分のデータを取得します。
func (a *RMSApi) GetShopCalendar(fromDate string, period int) (*ShopBizApiResponse, error) {
	if a.authorization == "" {
		return nil, errors.New("Uninitialized")
	}

	req, _ := http.NewRequest("GET", SHOP_CALENDAR_URL, nil)
	req.Header.Set("Authorization", "ESA "+a.authorization)
	req.Header.Set("Content-Type", "application/sml; charset=utf-8")

	params := req.URL.Query()
	_, err := time.Parse("2006-01-02", fromDate)
	if fromDate != "" && err == nil {
		params.Add("fromDate", fromDate)
	}
	if period > 0 && period <= 180 {
		params.Add("period", fmt.Sprintf("%d", period))
	}
	req.URL.RawQuery = params.Encode()

	client := new(http.Client)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)
	result := ShopBizApiResponse{}
	err = xml.Unmarshal(byteArray, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
