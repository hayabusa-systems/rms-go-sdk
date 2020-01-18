/*
このパッケージはGoLangから楽天の店舗管理用API RMS WEB SERVICEを操作するためのSDKです。
*/
package rms

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

const (
	// RMS WEB SERVICEの楽天ペイ受注APIの注文検索用のエンドポイントです。
	SEARCH_ORDER_URL = "https://api.rms.rakuten.co.jp/es/2.0/order/searchOrder/"

	// RMS WEB SERVICEの楽天ペイ受注APIの注文情報の取得用のエンドポイントです。
	GET_ORDER_URL = "https://api.rms.rakuten.co.jp/es/2.0/order/getOrder/"

	// RMS WEB SERVICEの楽天ペイ受注APIのひとことメモ、担当者情報等の更新用のエンドポイントです。
	UPDATE_ORDER_MEMO_URL = "https://api.rms.rakuten.co.jp/es/2.0/order/updateOrderMemo/"
)

// SearchOrderDateType は期間検索種別を表します。
type SearchOrderDateType int

const (
	DATE_TYPE_ORDER_DATE                    = iota + 1 // 注文日
	DATE_TYPE_ORDER_CONFIRM_DATE                       // 注文確認日
	DATE_TYPE_ORDER_FIX_DATE                           // 注文確定日
	DATE_TYPE_SHIPPING_DATE                            // 発送日
	DATE_TYPE_SHIPPING_COMPLETE_REPORT_DATE            //発送完了報告日
	DATE_TYPE_PAYMENT_FIX_DATE                         // 決済確定日
)

type (
	/*** RMSとの通信時に使用 ***/
	/*** 共通 ***/

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

	/*** Request ***/

	/*** Response ***/

	// CommonMessageModelResponse はRMS WEB SERVICEのレスポンスの共通部分の構造体です。
	CommonMessageModelResponse struct {
		// MessgeType はメッセージ種別です。以下のいずれかが入力されます。
		// INFO
		// ERROR
		// WARNING
		MessageType string `json:"messageType"`

		// MessageCode はメッセージコードです。
		// メッセージコードは https://webservice.rms.rakuten.co.jp/merchant-portal/view?contents=/ja/common/1-1_service_index/rakutenpayorderapi/rakutenpaymsgcodereference を参照してください(要ログイン)
		MessageCode string `json:"messageCode"`

		// Message はメッセージです。
		Message string `json:"message"`
	}

	/*** Request ***/
	/*** searchOrder ***/

	// SearchOrderSortModel は楽天ペイ受注APIで注文検索の検索条件のうち、ソートに関する条件です。
	SearchOrderSortModel struct {
		// SortColumn は並び替え項目です。以下のいずれかを指定することができます。
		// 1: 注文日時
		SortColumn int `json:"sortColumn"`

		// SortDirection は並び替え方法です。以下のいずれかを指定することができます。
		// 1: 昇順
		// 2: 降順
		SortDirection int `json:"sortDirection"`
	}

	// SearchOrderPaginationRequestModel は楽天ペイ受注APIで注文検索の検索条件のうち、ページングに関する条件です。
	SearchOrderPaginationRequestModel struct {
		// RequestRecordsAmount は1ページあたりの取得結果数です。最大1,000件まで取得可能です。
		RequestRecordsAmount int `json:"requestRecordsAmount"`

		// RequestPage はリクエストページ番号です。
		RequestPage int `json:"requestPage"`

		// SortModelList はソート条件です。
		SortModelList []SearchOrderSortModel `json:"SortModelList,omitempty"`
	}

	// SearchOrderRequest は楽天ペイ受注APIで注文検索の検索条件です。
	SearchOrderReuquest struct {
		// OrderProgressList はステータスリストです。
		// ステータスは以下のいずれかが入ります。
		// 100: 注文確認待ち
		// 200: 楽天処理中
		// 300: 発送待ち
		// 400: 変更確定待ち
		// 500: 発送済
		// 600: 支払い手続き中
		// 700: 支払い手続き済
		// 800: キャンセル確定待ち
		// 900: キャンセル確定
		OrderProgressList []int `json:"orderProgressList,omitempty"`

		// SubStatusIdList はサブステータスIDリストです。
		// ユーザが作成したサブステータスを指定することができます。
		SubStatusIDList []int `json:"subStatusIdList,omitempty"`

		// DateType は期間検索種別です。
		// この項目は必須です。
		// 種別は以下のいずれかが入ります。
		// 1: 注文日
		// 2: 注文確認日
		// 3: 注文確定日
		// 4: 発送日
		// 5: 発送完了報告日
		// 6: 決済確定日
		DateType int `json:"dateType"`

		// StartDatetime は期間検索開始日時です。過去2年以内の注文を指定することが可能です。
		// この項目は必須です。
		StartDatetime JsonTime `json:"startDatetime"`

		// EndDatetime は期間検索終了日時です。開始日から63日以内を指定することができます。
		// この項目は必須です。
		EndDatetime JsonTime `json:"endDatetime"`

		// OrderTypeList は販売種別リストです。以下のいずれかを指定することができます。
		// 1: 通常購入
		// 4: 定期購入
		// 5: 頒布会
		// 6: 予約商品
		OrderTypeList []int `json:"orderTypeList,omitempty"`

		// SettlementMethod は支払い方法名です。以下のいずれかを指定することができます。
		// 1: クレジットカード
		// 2: 代金引換
		// 3: 後払い
		// 4: ショッピングクレジット/ローン
		// 5: オートローン
		// 6: リース
		// 7: 請求書払い
		// 9: 銀行振込
		// 12: Apple Pay
		// 13: セブンイレブン(前払)
		// 14: ローソン、郵便局ATM等(前払)
		// 16: Alipay
		// 17: PayPal
		// 21: 後払い決済(楽天市場の共通決済)
		SettlementMethod *int `json:"settlementMethod,omitempty"`

		// DeliveryName は配送方法です。
		DeliveryName *string `json:"deliveryName,omitempty"`

		// ShippingDateBlankFlag は発送日未指定有無フラグです。以下のいずれかを指定することができます。
		// 0: 発送日の指定の有無によらず取得
		// 1: 発送日が未指定のものだけを取得
		ShippingDateBlankFlag *int `json:"shippingDateBlankFlag,omitempty"`

		// ShippingNumberBlankFlag はお荷物伝票番号未指定有無フラグです。以下のいずれかを指定することができます。
		// 0: お荷物伝票番号の指定の有無によらず取得
		// 1: お荷物伝票番号が未指定のものだけを取得
		ShippingNumberBlankFlag *int `json:"shippingNumberBlankFlag,omitempty"`

		// SearchKeywordType は検索キーワード種別です。次のいずれかを指定することができます。
		// 0: なし
		// 1: 商品名
		// 2: 商品番号
		// 3: ひとことメモ
		// 4: 注文者お名前
		// 5: 注文者お名前フリガナ
		// 6: 送付先お名前
		SearchKeywordType *int `json:"searchKeywordType,omitempty"`

		// SearchKeyword は検索キーワードです。32文字以下の入力を受け付けます。
		SearchKeyword *string `json:"searchKeyword,omitempty"`

		// MailSendType は注文メールアドレス種別です。以下のいずれかを指定することができます。
		// 0: PC/モバイル
		// 1: PC
		// 2: モバイル
		MailSendType *int `json:"mailSendType,omitempty"`

		// OrdererMailAddress は注文者メールアドレスです。完全一致である必要があります。
		OrdererMailAddress *string `json:"ordererMailAddress,omitempty"`

		// PhoneNumberType は電話番号種別です。以下のいずれかを指定することができます。
		// 0: 注文者
		// 1: 送付先
		PhoneNumberType *int `json:"phoneNumberType,omitempty"`

		// PhoneNumber は電話番号です。完全一致である必要があります。
		PhoneNumber *string `json:"phoneNumber,omitempty"`

		// ReserveNumber は申込番号です。完全一致である必要があります。
		ReserveNumber *string `json:"reserveNumber,omitempty"`

		// PurchaseSiteType は購入サイトリストです。以下のいずれかを指定することがあります。
		// 0: すべて
		// 1: PCで注文
		// 2: モバイルで注文
		// 3: スマートフォンで注文
		// 4: タブレットで注文
		PurchaseSiteType *int `json:"purchaseSiteType,omitempty"`

		// AsurakuFlag はあす楽希望フラグです。以下のいずれかを指定することができます。
		// 0: あす楽希望の有無にかかわらず取得
		// 1: あす楽希望のものだけを取得
		AsurakuFlag *int `json:"asurakuFlag,omitempty"`

		// CouponUseFlag はクーポン利用有無フラグです。以下のいずれかを指定することができます。
		// 0: クーポン利用の有無にかかわらず取得
		// 1: クーポン利用のものだけを取得
		CouponUseFlag *int `json:"couponUseFlag,omitempty"`

		// DrugFlag は医薬品受注フラグです。以下のいずれかを指定することができます。
		// 0: 医薬品の有無にかかわらず取得
		// 1: 医薬品を含む注文だけを取得
		DrugFlag *int `json:"drugFlag,omitempty"`

		// OverseasFlag は海外カゴ注文フラグです。以下のいずれかを指定することができます。
		// 0: 海外カゴ注文の有無にかかわらず取得
		// 1: 海外カゴ注文のみ取得
		OverseasFlag *int `json:"overseasFlag,omitempty"`

		// SearchOrderPaginationRequestModel はページングに関する情報です。
		SearchOrderPaginationRequestModel `json:"PaginationRequestModel,omitempty"`
	}

	/*** getOrder ***/

	// GetOrderRequest は楽天ペイ受注APIの注文情報取得時の検索条件です。
	GetOrderRequest struct {
		// OrderNumberList は注文番号リストです。最大100件まで指定可能で、過去2年分の注文が対象です。
		// この項目は必須です。
		OrderNumberList []string `json:"orderNumberList"`

		// Version はバージョン番号です。現在は必須ではありませんが、今後は必須化される予定です。
		// 次のいずれかが入力されます。
		// 1: 2018/02版
		// 2: 2018/11版
		// 3: 2019/04版
		// 4: 2019/12版
		Version int `json:"version"`
	}

	/*** Response ***/
	/*** searchOrder ***/

	// SearchOrderPaginationResponseModel は楽天ペイ受注APIの注文検索時に取得したデータ数、ページ数が入力されます。
	SearchOrderPaginationResponseModel struct {
		// TotalRecordsAmount は総結果数です。
		TotalRecordsAmount int `json:"totalRecordsAmount"`

		// TotalPages は総ページ数です。
		TotalPages int `json:"totalPages"`

		// RequestPage はリクエストページ番号です。
		RequestPage int `json:"requestPage"`
	}

	// SearchOrderResponse は楽天ペイ受注APIの注文検索時に取得できるデータです。
	SearchOrderResponse struct {
		// CommonMessageModelResponseList はメッセージモデルリストです。ここにはエラー情報が含まれます。
		CommonMessageModelResponseList []CommonMessageModelResponse `json:"MessageModelList"`

		// OrderNumberList は注文番号リストです。該当する注文番号の一覧が取得できます。
		OrderNumberList []string `json:"orderNumberList"`

		// SearchOrderPaginationResponseModel はページングレスポンスモデルです。ページングに関する情報を取得することができます。
		SearchOrderPaginationResponseModel `json:"PaginationResponseModel"`
	}

	/*** getOrder ***/

	// GetOrderMessageModel は楽天ペイ受注APIの注文情報の取得時のエラー情報等が含まれます。共通メッセージ以外に注文番号が含まれるため、別のクラスとして定義されます。
	GetOrderMessageModel struct {
		// CommonMessageModelResponse はエラー情報が含まれます。
		CommonMessageModelResponse

		// OrderNumber は注文番号です。
		OrderNumber string `json:"orderNumber"`
	}

	// GetOrderOrdererModel は楽天ペイ受注APIの注文情報の取得で得られる注文者の情報です。
	GetOrderOrdererModel struct {
		// ZipCode1 は郵便番号1です。3桁の数値です。(0はじまりがあるため文字列です。)
		ZipCode1 string `json:"zipCode1"`

		// ZipCode2 は郵便番号2です。4桁の数値です。(0はじまりがあるため文字列です。)
		ZipCode2 string `json:"zipCode2"`

		// Prefecture は都道府県です。
		Prefecture string `json:"prefecture"`

		// City は郡市区です。
		City string `json:"city"`

		// SubAddress はCity以降の住所です。
		SubAddress string `json:"subAddress"`

		// FamilyName は姓です。
		FamilyName string `json:"familyName"`

		// FirstName は名です。
		FirstName string `json:"firstName"`

		// FamilyNameKana は姓カナです。
		FamilyNameKana *string `json:"familyNameKana"`

		// FirstNameKana は名カナです。
		FirstNameKana *string `json:"firstNameKana"`

		// PhoneNumber1 は電話暗号1です。電話番号3までのうち、1つまでがnullの可能性があります。
		PhoneNumber1 string `json:"phoneNumber1"`

		// PhoneNumber2 は電話暗号2です。電話番号3までのうち、1つまでがnullの可能性があります。
		PhoneNumber2 string `json:"phoneNumber2"`

		// PhoneNumber3 は電話暗号3です。電話番号3までのうち、1つまでがnullの可能性があります。
		PhoneNumber3 string `json:"phoneNumber3"`

		// EmailAddress はメールアドレスです。マスキングされたものが取得されます。
		EmailAddress string `json:"emailAddress"`

		// Sex は性別です。
		Sex string `json:"sex"`

		// BirthYear は誕生日(年)です。
		BirthYear int `json:"birthYear"`

		// BirthMonth は誕生日(月)です。
		BirthMonth int `json:"birthMonth"`

		// BirthDay は誕生日(日)です。
		BirthDay int `json:"birthDay"`
	}

	// GetOrderSettlementModel は楽天ペイ受注APIの注文情報の取得で得られる支払情報です。
	GetOrderSettlementModel struct {
		// SettlementMethod は支払い方法名です。
		SettlementMethod string `json:"settlementMethod"`

		// RpaySettlementFlag は楽天市場の共通決済手段フラグです。APIのバージョンが2以降の場合取得することができます。以下のいずれかが含まれます。
		// 0: 選択制決済
		// 1: 楽天市場の共通決済手段
		RpaySettlementFlag int `json:"rpaySettlementFlag"`

		// CardName はクレジットカード種類です。SettlementMethod がクレジットカードの場合のみ値が入力されます。
		CardName *string `json:"cardName"`

		// CardNumber はクレジットカード番号です。SettlementMethod がクレジットカードの場合のみ値が入力されます。
		CardNumber *string `json:"cardNumber"`

		// CardOwner はクレジットカード名義人です。SettlementMethod がクレジットカードの場合のみ値が入力されます。
		CardOwner *string `json:"cardOwner"`

		// CardYm はクレジットカード有効期限です。SettlementMethod がクレジットカードの場合のみ値が入力されます。
		CardYm *string `json:"cardYm"`

		// CardPayType はクレジットカード支払い方法です。SettlementMethod がクレジットカードの場合のみ値が入力されます。以下のいずれかが入力されます。
		// 0: 一括払い
		// 1: リボ払い
		// 2: 分割払い
		// 3: その他払い
		// 4: ボーナス一括払い
		CardPayType int `json:"cardPayType"`

		// CardInstallmentDesc はクレジットカードの支払回数です。SettlementMethod がクレジットカードの場合かつ、CardPayType が分割払いの場合のみ値が入力されます。以下のいずれかが入力されます。
		// 103: 3回払い
		// 105: 5回払い
		// 106: 6回払い
		// 110: 10回払い
		// 112: 12回払い
		// 115: 15回払い
		// 118: 18回払い
		// 120: 20回払い
		// 124: 24回払い
		CardInstallmentDesc *string `json:"cardInstallmentDesc"`
	}

	// GetOrderDeliveryModel は楽天ペイ受注APIの注文情報の取得で得られる配送方法です。
	GetOrderDeliveryModel struct {
		// DeliveryName は配送方法です。店舗設定で設定した配送方法が入力されます。
		DeliveryName string `json:"deliveryName"`

		// DeliveryClass は配送区分です。以下のいずれかが入力されます。
		// 0: 選択なし
		// 1: 普通
		// 2: 冷蔵
		// 3: 冷凍
		// 4: その他1
		// 5: その他2
		// 6: その他3
		// 7: その他4
		// 8: その他5
		DeliveryClass *int `json:"deliveryClass"`
	}

	// GetOrderPointModel は楽天ペイ受注APIの注文情報の取得で得られるポイントの利用額です。
	GetOrderPointModel struct {
		// UsedPoint はポイント利用額です。
		UsedPoint int `json:"usedPoint"`
	}

	// GetOrderWrappingModel は楽天ペイ受注APIの注文情報の取得で得られるラッピングに関する情報です。
	GetOrderWrappingModel struct {
		// Title はラッピングタイトルです。次のいずれかが入力されます。
		// 1: 包装紙
		// 2: リボン
		Title int `json:"title"`

		// Name はラッピング名です。
		Name string `json:"name"`

		// Price は料金です。
		Price *int `json:"price"`

		// IncludeTaxFlag は税込み別です。次のいずれかが入力されます。
		// 0: 税別
		// 1: 税込
		IncludeTaxFlag int `json:"includeTaxFlag"`

		// DeleteWrappingFlag はラッピング削除フラグです。
		DeleteWrappingFlag int `json:"deleteWrappingFlag"`

		// TaxRate はラッピング税率です。APIのバージョンが3以降の場合取得可能です。
		TaxRate int `json:"taxRate"`

		// TaxPriceはラッピング税額です。APIのバージョンが3以降の場合取得可能です。
		TaxPrice int `json:"taxPrice"`
	}

	// GetOrderSenderModel は楽天ペイ受注APIの注文情報の取得で得られる発送先に関する情報です。
	GetOrderSenderModel struct {
		// ZipCode1 は郵便番号1です。3桁の数値が入力されます。(0始まりがあるため文字列です。)
		ZipCode1 string `json:"zipCode1"`

		// ZipCode2 は郵便番号2です。4桁の数値が入力されます。(0始まりがあるため文字列です。)
		ZipCode2 string `json:"zipCode2"`

		// Prefecture は都道府県です。
		Prefecture string `json:"prefecture"`

		// City は郡市区です。
		City string `json:"city"`

		// SubAddress はCity以降の住所です。
		SubAddress string `json:"subAddress"`

		// FamilyName は姓です。
		FamilyName string `json:"familyName"`

		// FirstName は名です。
		FirstName string `json:"firstName"`

		// FamilyNameKana は姓カナです。
		FamilyNameKana *string `json:"familyNameKana"`

		// FirstNameKana は名カナです。
		FirstNameKana *string `json:"firstNameKana"`

		// PhoneNumber1 は電話番号1です。電話番号3までのうち1つだけnullの場合があります
		PhoneNumber1 string `json:"phoneNumber1"`

		// PhoneNumber2 は電話番号2です。電話番号3までのうち1つだけnullの場合があります
		PhoneNumber2 string `json:"phoneNumber2"`

		// PhoneNumber3 は電話番号3です。電話番号3までのうち1つだけnullの場合があります
		PhoneNumber3 string `json:"phoneNumber3"`

		// IsolatedIslandFlag は離島フラグです。APIのバージョンが2以降の場合取得可能です。以下のいずれかが入力されます。
		// 0: 離島ではない
		// 1: 離島である
		IsolatedIslandFlag int `json:"isolatedIslandFlag"`
	}

	// GetOrderItemModel は楽天ペイ受注APIの注文情報の取得で得られる販売した商品の情報です。
	GetOrderItemModel struct {
		// ItemDetailID は商品明細IDです。
		ItemDetailID int `json:"itemDetailId"`

		// ItemName は商品名です。
		ItemName string `json:"itemName"`

		// ItemID は商品IDです。
		ItemID int `json:"itemId"`

		// ItemNumber は商品番号です。項目選択肢別在庫が指定された商品の場合、以下のルールで値が表示されます。
		// 商品番号 + 項目選択肢ID(横軸) + 項目選択肢ID(縦軸)
		ItemNumber *string `json:"itemNumber"`

		// ManageNumber は商品管理番号です。
		ManageNumber string `json:"manageNumber"`

		// Price は単価です。
		Price int `json:"price"`

		// Units は個数です。
		Units int `json:"units"`

		// IncludePostageFlag は送料込別です。以下のいずれかが入力されます。
		// 0: 送料別
		// 1: 送料込みもしくは送料無料
		IncludePostageFlag int `json:"includePostageFlag"`

		// IncludePostageFlag は税込み別です。以下のいずれかが入力されます。
		// 0: 税別
		// 1: 税込み
		IncludeTaxFlag int `json:"includeTaxFlag"`

		// IncludeCashOnDeliveryPostageFlag は代引き手数料込別です。以下のいずれかが入力されます。
		// 0: 代引き手数料別
		// 1: 代引き手数料込み
		IncludeCashOnDeliveryPostageFlag int `json:"includeCashOnDeliveryPostageFlag"`

		// SelectedChoice は項目・選択肢です。注文種別が通常購入、共同購入、定期購入、頒布会、予約商品の場合入力されます。
		// 表示形式は"横軸項目名:横軸選択肢 縦軸項目名:縦軸選択肢" で、HTMLタグは除去されます。
		SelectedChoice *string `json:"selectedChoice"`

		// PointRate はポイント倍率です。
		PointRate int `json:"pointRate"`

		// PointType はポイントタイプです。APIのバージョンが2以降の場合入力されます。以下のいずれかが入力されます。
		// 0: 変倍なし
		// 1: 店舗別変倍
		// 2: 商品別変倍
		// -99: エラー時無効値
		PointType int `json:"pointType"`

		// InventoryType は在庫タイプです。以下のいずれかが入力されます。
		// 0: 在庫設定無し
		// 1: 通常在庫設定
		// 2: 項目選択肢在庫設定
		InventoryType int `json:"inventoryType"`

		// DelvdateInfo は納期情報です。
		DelvdateInfo *string `json:"delvdateInfo"`

		// RestoreInventoryFlag は在庫連動オプションです。以下のいずれかが入力されます。
		// 0: 商品の設定に従う
		// 1: 在庫連動する
		// 2: 在庫連動しない
		RestoreInventoryFlag int `json:"restoreInventoryFlag"`

		// DealFlag hあ楽天スーパーDEAL商品フラグです。APIのバージョンが2以降の場合入力されます。以下のいずれかが入力されます。
		// 0: 楽天スーパーDEAL商品ではない
		// 1: 楽天スーパーDEAL商品である
		DealFlag int `json:"dealFlag"`

		// DrugFlag は医薬品フラグです。APIのバージョンが2以降の場合入力されます。以下のいずれかが入力されます。
		// 0: 医薬品ではない。
		// 1: 医薬品である。
		DrugFlag int `json:"drugFlag"`

		// DeleteItemFlag は商品削除フラグです。以下のいずれかが入力されます。
		// 0: 商品を削除しない
		// 1: 商品を削除する
		DeleteItemFlag int `json:"deleteItemFlag"`

		// TaxRate は商品税率です。APIのバージョンが3以降の場合入力されます。店舗で税率が更新されていない注文の場合NULLとなります。
		TaxRate float64 `json:"TaxRate"`

		// PriceTaxIncl は商品毎税込価格です。APIのバージョンが3以降の場合入力されます。
		// 税込商品の場合、商品単価は商品毎税込価格となります。
		// 税別商品の場合、税込価格(商品単価 * (1 + 税率))です。端数処理は店舗設定に準じます。
		PriceTaxIncl int `json:"priceTaxIncl"`

		// IsSingleItemShipping は単品配送フラグです。APIのバージョンが4以降の場合入力されます。以下のいずれかが入力されます。
		// 0: 単品配送ではない
		// 1: 単品配送である
		IsSingleItemShipping int `json:"isSingleItemShipping"`
	}

	// GetOrderDeliveryCvsModel は楽天ペイ受注APIの注文情報の取得で得られる配送のコンビニ・郵便局受け取りに関する情報です。
	GetOrderDeliveryCvsModel struct {
		// CvsCode はコンビニコードです。以下のいずれかが入力されます。
		// 1: ファミリーマート
		// 20: ミニストップ
		// 40: サークルK
		// 41: サンクス
		// 50: ローソン
		// 60: 郵便局
		// 70: スリーエフ
		// 71: エブリワン
		// 72: ココストア
		// 74: セーブオン
		// 80: デイリーヤマザキ
		// 81: ヤマザキデイリーストア
		// 82: ニューヤマザキデイリーストア
		// 85: ニューデイズ
		// 90: ポプラ
		// 91: くらしハウス
		// 92: スリーエイト
		// 93: 生活彩家
		CvsCode *int `json:"cvsCode"`

		// StoreGenreCode はストア分類コードです。APIのバージョンが2以降の場合入力されます。
		StoreGenreCode *string `json:"storeGenreCode"`

		// StoreCode はストアコードです。APIのバージョンが2以降の場合入力されます。
		StoreCode *string `json:"storeCode"`

		// StoreName はストア名称です。APIのバージョンが2以降の場合入力されます。
		StoreName *string `json:"storeName"`

		// StoreZip は郵便番号です。APIのバージョンが2以降の場合入力されます。
		StoreZip *string `json:"storeZip"`

		// StorePrefecture は都道府県です。APIのバージョンが2以降の場合のみ入力されます。
		StorePrefecture *string `json:"storePrefecture"`

		// StoreAddress はその他住所です。APIのバージョンが2以降の場合入力されます。
		StoreAddress *string `json:"storeAddress"`

		// AreaCode は発注エリアコードです。APIのバージョンが2以降の場合入力されます。
		AreaCode *string `json:"areaCode"`

		// Depo はセンターデポコードです。APIのバージョンが2以降の場合入力されます。
		Depo *string `json:"depo"`

		// OpenTime は開店時間です。APIのバージョンが2以降の場合入力されます。
		OpenTime *string `json:"openTime"`

		// CloseTime は閉店時間です。APIのバージョンが2以降の場合のみ入力されます。
		CloseTime *string `json:"closeTime"`

		// CvsRemarks は特記事項です。APIのバージョンが2以降の場合のみ入力されます。
		CvsRemarks *string `json:"cvsRemrks"`
	}

	// GetOrderCouponModel は楽天ペイ受注APIの注文情報の取得で得られる利用されたクーポンの情報です。
	GetOrderCouponModel struct {
		// CouponCode はクーポンコードです。
		CouponCode string `json:"couponCode"`

		// ItemID はクーポン対象の商品IDです。該当する商品がない場合は0が指定されます。
		ItemID int `json:"itemId"`

		// CouponName はクーポン名です。
		CouponName string `json:"couponName"`

		// CouponSummary はクーポン効果(サマリー)です。
		CouponSummary string `json:"couponSummary"`

		// CouponCapital はクーポン原資です。以下のいずれかが入力されます。
		// ・ショップ
		// ・メーカー
		// ・サービス
		CouponCapital string `json:"couponCapital"`

		// CouponCapitalCode はクーポン原資コードです。以下のいずれかが入力されます。
		// 1: ショップ
		// 2: メーカー
		// 3: サービス
		CouponCapitalCode int `json:"couponCapitalCode"`

		// ExpiryDate は有効期限です。日付のみ取得可能です。
		ExpiryDate JsonDate `json:"expiryDate"`

		// CouponPrice はクーポン割引単価です。
		CouponPrice int `json:"couponPrice"`

		// CouponUnit はクーポン利用数です。
		CouponUnit int `json:"couponUnit"`

		// CouponTotalPrice はクーポン利用金額です。クーポン割引単価またはクーポン利用数がnullの場合、-9999になります。
		CouponTotalPrice int `json:"couponTotalPrice"`
	}

	// GetOrderChangeReasonModel は楽天ペイ受注APIの注文情報の取得で得られる注文の変更理由です。
	GetOrderChangeReasonModel struct {
		// ChangeID は変更IDです。
		ChangeID int `json:"changeId"`

		// ChangeType は変更種別です。以下のいずれかが入力されます。
		// 0: キャンセル申請
		// 1: キャンセル確定
		// 2: キャンセル完了
		// 3: キャンセル取消
		// 4: 変更申請
		// 5: 変更確定
		// 6: 変更完了
		// 7: 変更取消
		// 8: 注文確認
		// 9: 再決済手続き
		ChangeType *int `json:"changeType"`

		// ChangeTypeDetail は変更種別(詳細)です。以下のいずれかが入力されます。
		// 0: 減額
		// 1: 増額
		// 2: その他
		// 10: 支払い方法変更
		// ※その他は後払い決済選択注文で金額以外の変更が行われた場合のみ
		ChangeTypeDetail int `json:"changeTypeDetail"`

		// ChangeReason は変更理由です。以下のいずれかが入力されます。
		// 0: 店舗様都合
		// 1: お客様都合
		ChangeReason *int `json:"changeReason"`

		// ChangeReasonDetail は変更理由(小分類)です。以下のいずれかが入力されます。
		// 1: キャンセル
		// 2: 受取後の返品
		// 3: 長期不在による受取拒否
		// 4: 未入金
		// 5: 代引決済の受取拒否
		// 6: お客様都合 - その他
		// 8: 欠品
		// 10: 店舗様都合 - その他
		// 13: 発送遅延
		// 14: 顧客・配送対応先注意表示
		// 15: 返品(破損・品間違い)
		ChangeReasonDetail *int `json:"changeReasonDetail"`

		// ChangeApplyDatetime は変更申請日です。
		ChangeApplyDatetime *JsonTime `json:"changeApplyDatetime"`

		// ChangeFixDatetime は変更確定日です。
		ChangeFixDatetime *JsonTime `json:"changeFixDatetime"`

		// ChangeCompleteDatetime は変更完了日です。
		ChangeCompleteDatetime *JsonTime `json:"changeCmplDatetime"`
	}

	// GetOrderTaxSummaryModel は楽天ペイ受注APIの注文情報の取得で得られる税金を含む請求額に関する情報です。
	GetOrderTaxSummaryModel struct {
		// TaxRate は税率です。
		TaxRate float64 `json:"taxRate"`

		// ReqPrice は請求金額です。計算式は商品金額 + 送料 + ラッピング料 + 決済手数料 + 注文者負担金 - クーポン割引額 - 利用ポイント数。
		// なお以下の場合、-9999となります。
		// ・送料未確定
		// ・代引手数料未確定
		ReqPrice int `json:"reqPrice"`

		// ReqPriceTax は請求額に対する税額です。請求金額に対する税額です。請求金額が-9999の場合、ReqPriceTaxは-9999となります。
		ReqPriceTax int `json:"reqPriceTax"`

		// TotalPrice は合計金額です。計算式は商品金額 + 送料 + ラッピング料です。送料が未確定の場合、-9999となります。
		TotalPrice int `json:"totalPrice"`

		// PaymentCharge は決済手数料です。代引手数料未確定の場合、-9999となります。
		PaymentCharge int `json:"paymentCharge"`

		// CouponPrice はクーポン割引額です。
		CouponPrice int `json:"couponPrice"`

		// Point は利用ポイント数です。
		Point int `json:"point"`
	}

	// GetOrderShippingModel は楽天ペイ受注APIの注文情報の取得で得られる配送情報です。
	GetOrderShippingModel struct {
		// ShippingDetailID は発送明細IDです。このIDは楽天が発行したもので、更新・削除の場合に使用します。
		ShippingDetailID int `json:"shippingDetailId"`

		// ShippingNumber はお荷物伝票番号です。
		ShippingNumber *string `json:"shippingNumber"`

		// DeliveryCompany は配送会社です。以下のいずれかが入力されます。
		// 1000: その他
		// 1001: ヤマト運輸
		// 1002: 佐川急便
		// 1003: 日本郵便
		// 1004: 西濃運輸
		// 1005: セイノースーパーエクスプレス
		// 1006: 福山通運
		// 1007: 名鉄運輸
		// 1008: トナミ運輸
		// 1009: 第一貨物
		// 1010: 新潟運輸
		// 1011: 中越運送
		// 1012: 岡山県貨物運送
		// 1013: 久留米運送
		// 1014: 山陽自動車運送
		// 1015: 日本トラック
		// 1016: エコ配
		// 1017: EMS
		// 1018: DHL
		// 1019: FedEx
		// 1020: UPS
		// 1021: 日本通運
		// 1022: TNT
		// 1023: OCS
		// 1024: USPS
		// 1025: SFエクスプレス
		// 1026: Aramex
		// 1027: SGHグローバル・ジャパン
		// 1028: Rakuten EXPRESS
		DeliveryCompany *string `json:"deliveryCompany"`

		// DeliveyCompanyName は配送会社名です。
		DeliveyCompanyName *string `json:"deliveryCompanyName"`

		// ShippingDate は発送日です。
		ShippingDate *JsonDate `json:"shippingDate"`
	}

	// GetOrderPackageModel は楽天ペイ受注APIの注文情報の取得で得られる送料に関する情報です。
	GetOrderPackageModel struct {
		// BascketID は送付先IDです。
		BascketID int `json:"bascketId"`

		// PostagePrice は送料です。未設定の場合、-9999になります。
		PostagePrice int `json:"postagePrice"`

		// PostageTaxRate は送料税率です。APIのバージョンが3以降の場合のみ入力されます。
		PostageTaxRate float64 `json:"postageTaxRate"`

		// DeliveryPrice は代引料です。未設定の場合、-9999になります。
		DeliveryPrice int `json:"deliveryPrice"`

		// DeliveryTaxRate は代引料税率です。APIのバージョンが3以降の場合のみ入力されます。
		DeliveryTaxRate float64 `json:"deliveryTaxRate"`

		// GoodsTax は送付先外税合計です。税込商品の場合は0、未設定の場合は-9999が入力されます。APIのバージョンが3以降の場合のみ入力されます。
		GoodsTax int `json:"goodsTax"`

		// GoodsPrice は商品合計金額です。商品金額 + ラッピング料です。
		GoodsPrice int `json:"goodsPrice"`

		// TotalPrice は合計金額です。 商品金額 + 送料 + ラッピング料です。未確定の場合は-9999になります。APIのバージョンが1を指定すると代引手数料が含まれます。
		TotalPrice int `json:"totalPrice"`

		// Noshi はのしです。
		Noshi *string `json:"noshi"`

		// PackageDeleteFlag は送付先モデル削除フラグです。以下のいずれかが入力されます。
		// 0: 送付先モデルを削除しない
		// 1: 送付先モデルを削除する
		PackageDeleteFlag int `json:"packageDeleteFlag"`

		// GetOrderSenderModel は送付者モデルです。
		GetOrderSenderModel `json:"senderModel"`

		// ItemModelList は商品モデルリストです。
		ItemModelList []GetOrderItemModel `json:"ItemModelList"`

		// ShippingModelList は発送モデルリストです。
		ShippingModelList []GetOrderShippingModel `json:"ShippingModelList"`

		// GetOrderDeliveryCvsModel はコンビニ配送モデルです。
		GetOrderDeliveryCvsModel `json:"DeliveryCvsModel"`

		// DefaultDeliveryCompanyCode は購入時配送会社です。APIのバージョンが4以降の場合のみ入力されます。以下のいずれかが入力されます。
		// 1000: その他
		// 1001: ヤマト運輸
		// 1002: 佐川急便
		// 1003: 日本郵便
		// 1004: 西濃運輸
		// 1005: セイノースーパーエクスプレス
		// 1006: 福山通運
		// 1007: 名鉄運輸
		// 1008: トナミ運輸
		// 1009: 第一貨物
		// 1010: 新潟運輸
		// 1011: 中越運送
		// 1012: 岡山県貨物運送
		// 1013: 久留米運送
		// 1014: 山陽自動車運送
		// 1015: 日本トラック
		// 1016: エコ配
		// 1017: EMS
		// 1018: DHL
		// 1019: FedEx
		// 1020: UPS
		// 1021: 日本通運
		// 1022: TNT
		// 1023: OCS
		// 1024: USPS
		// 1025: SFエクスプレス
		// 1026: Aramex
		// 1027: SGHグローバル・ジャパン
		// 1028: Rakuten EXPRESS
		DefaultDeliveryCompanyCode string `json:"defaultDeliveryCompanyCode"`
	}

	// GetOrderOrderModel は楽天ペイ受注APIの注文情報の取得で得られる注文情報です。
	GetOrderOrderModel struct {
		// OrderNumber は注文番号です。
		OrderNumber string `json:"orderNumber"`
		// OrderProgress はステータスです。以下のいずれかが入力されます。
		// 100: 注文確認待ち
		// 200: 楽天処理中
		// 300: 発送待ち
		// 400: 変更確定待ち
		// 500: 発送済み
		// 600: 支払い手続き中
		// 700: 支払い手続き済
		// 800: キャンセル確定待ち
		// 900: キャンセル確定
		OrderProgress int `json:"orderProgress"`

		// SubStatusID はサブステータスIDです。
		SubStatusID *int `json:"subStatusId"`

		// SubStatusName はサブステータスです。
		SubStatusName *string `json:"subStatusName"`

		// OrderDatetime は注文日時です。
		OrderDatetime JsonTime `json:"orderDatetime"`

		// ShopOrderConfirmDatetime は注文確認日時です。
		ShopOrderConfirmDatetime *JsonTime `json:"shopOrderCfmDatetime"`

		// OrderFixDatetime は注文確定日時です。
		OrderFixDatetime *JsonTime `json:"orderFixDatetime"`

		// ShippingInstDatetime は発送指示日時です。
		ShippingInstDatetime *JsonTime `json:"shippingInstDatetime"`

		// ShippingCompleteReportDatetime は発送完了報告日時です。
		ShippingCompleteReportDatetime *JsonTime `json:"shippingCmplRptDatetime"`

		// CancelDueDate はキャンセル期限日です。
		CancelDueDate *JsonDate `json:"cancelDueDate"`

		// DeliveryDate はお届け日指定です。
		DeliveryDate *JsonDate `json"deliveryDate"`

		// ShippingTerm はお届け時間帯です。以下のいずれかが入力されます。
		// 0: なし
		// 1: 午前
		// 2: 午後
		// 9: その他
		// h1h2: h1時-h2時(h1は7~24, h2は07-24まで任意の数値で指定可能)
		ShippingTerm *int `json:"shippingTerm"`

		// Remarks はコメントです。備考欄に入るデータのことです。
		Remarks *string `json:"remarks"`

		// GiftCheckFlag はギフト配送希望フラグです。以下のいずれかが入力されます。
		// 0: ギフト注文ではない
		// 1: ギフト注文である
		GiftCheckFlag int `json:"giftCheckFlag"`

		// SeveralSenderFlag は複数送付先フラグです。以下のいずれかが入力されます。
		// 0: 複数配送先なし
		// 1: 複数配送先あり
		SeveralSenderFlag int `json:"severalSenderFlag"`

		// EqualSenderFlag は送付先一致フラグです。以下のいずれかが入力されます。
		// 0: 注文者と送付者の住所が同じではない。
		// 1: 注文が単数で注文者と送付先の住所が同じ
		EqualSenderFlag int `json:"equalSenderFlag"`

		// IsolatedIslandFlag は離島フラグです。以下のいずれかが入力されます。
		// 0: 送付先に離島が含めれていない。
		// 1: 送付先に離島あ含まれている。
		IsolatedIslandFlag int `json:"isolatedIslandFlag"`

		// RakutenMemberFlag は楽天会員フラグです。以下のいずれかが入力されます。
		// 0: 楽天会員ではない。
		// 1: 楽天会員である。
		RakutenMemberFlag int `json:"rakutenMemberFlag"`

		// CarrieCode は利用端末です。以下のいずれかが入力されます。
		// 0: PC(Windows系のスマートフォン・タブレットを含む)
		// 1: モバイル(docomo フィーチャーフォン)
		// 2: モバイル(KDDI フィーチャーフォン)
		// 3: モバイル(Softbank フィーチャーフォン)
		// 5: モバイル(WILLCOM フィーチャーフォン)
		// 11: スマートフォン(iPhone系)
		// 12: スマートフォン(Android系)
		// 19: スマートフォン(その他)
		// 21: タブレット(iPad系)
		// 22: タブレット(Android系)
		// 29: タブレット(その他)
		// 99: その他(不明な場合を含む)
		CarrieCode int `json:"carrierCode"`

		// EmailCarrierCode はメールキャリアコードです。以下のいずれかが入力されます。
		// 0: PC("@i.softbank.jp"を含む)
		// 1: DoCoMo
		// 2: au
		// 3: SoftBank
		// 5: WILLCOM
		// 99: その他
		EmailCarrierCode int `json:"emailCarrierCode"`

		// OrderType は注文種別です。以下のいずれかが入力されます。
		// 1: 通常購入
		// 4: 定期購入
		// 5: 頒布会
		// 6: 予約商品
		OrderType int `json:"orderType"`

		// ReserveNumber は申込番号です。定期購入、頒布会、予約商品に付与されます。
		ReserveNumber *string `json:reserveNumber"`

		// ReserveDeliveryCount は申込お届け回数です。予約商品は常に1、定期購入、頒布会は確定した回数が入力されます。
		ReserveDeliveryCount *int `json:"reserveDeliveryCount"`

		// CautionDisplayType は警告表示タイプです。以下のいずれかが入力されます。
		// 0: 表示なし
		// 1: 表示あり(注意喚起)
		// 2: 表示あり(キャンセル確定)
		CautionDisplayType int `json:"cautionDispalyType"`

		// RakutenConfirmFlag は楽天確認中フラグです。以下のいずれかが入力されます。
		// 0: 楽天確認中ではない
		// 1: 楽天確認中
		RakutenConfirmFlag int `json:"rakutenConfirmFlag"`

		// GoodsPrice は商品合計金額です。商品金額 + ラッピング料です。
		GoodsPrice int `json:"goodsPrice"`

		// GoodsTax は外税合計です。税込商品の場合は0、未確定の場合は-9999です。APIのバージョンが3以降の場合、この値ではなく請求金額に対する税額(reqPriceTax)をご使用ください。
		GoodsTax int `json:"goodsTax"`

		// PostagePrice は送料合計です。送付先が複数ある場合、その合計です。未確定の場合、-9999です。
		PostagePrice int `json:"postagePrice"`

		// DeliveryPrice は代引料合計です。代引手数料がかからない決済手段の場合は0、未確定の場合は-9999です。
		DeliveryPrice int `json:"deliveryPrice"`

		// PaymentCharge は決済手数料合計です。決済手数料がかからない決済手段の場合は0、未確定の場合は-9999です。APIのバージョンが2以降の場合のみ入力されます。
		PaymentCharge int `json:"paymentCharge"`

		// PaymentChargeTaxRate は決済手数料税率です。APIのバージョンが3以降の場合のみ入力されます。
		PaymentChargeTaxRate float64 `json:"paymentChargeTaxRate"`

		// TotalPrice は合計金額です。商品金額 + 送料 + ラッピング料です。未確定の場合は-9999です。
		TotalPrice int `json:"totalPrice"`

		// RequestPrice は請求金額です。商品金額 + 送料 + ラッピング料 + 決済手数料 + 注文者負担金 - クーポン利用総額 - ポイント利用額です。未確定の場合は-9999です。
		RequestPrice int `json:"requestPrice"`

		// CouponAllTotalPrice はクーポン利用総額です。
		CouponAllTotalPrice int `json:"requestPrice"`

		// CouponShopPrice は店舗発行クーポン利用額です。クーポン原資コードが1のクーポンが対象です。未確定の場合は-9999です。
		CouponShopPrice int `json:"couponShopPrice"`

		// CouponOtherPrice は楽天発行クーポン利用額です。クーポン原資コードが1以外のクーポンが対象です。未確定の場合は-9999です。
		CouponOtherPrice int `json:"couponOtherPrice"`

		// AdditionalFeeOccurAmountToUser は注文者負担金合計です。負担金がない場合は0、未確定の場合は-9999です。APIのバージョンが2以降の場合のみ入力されます。
		AdditionalFeeOccurAmountToUser int `json:"additionalFeeOccurAmountToUser"`

		// AdditionalFeeOccurAmountToShop は店舗負担金合計です。負担金がない場合は0、未確定の場合は-9999です。APIのバージョンが2以降の場合のみ入力されます。
		AdditionalFeeOccurAmountToShop int `json:"additionalFeeOccurAmountToShop"`

		// AsurakuFlag はあす楽希望フラグです。以下のいずれかが入力されます。
		// 0: 希望なし
		// 1: 希望あり
		AsurakuFlag int `json:"asurakuFlag"`

		// DrugFlag は医薬品受注フラグです。以下のいずれかが入力されます。
		// 0: 医薬品を含む注文ではない。
		// 1: 医薬品を含む注文である。
		DrugFlag int `json:"drugFlag"`

		// DealFlag は楽天スーパーDEAL商品受注フラグです。以下のいずれかが入力されます。
		// 0: 楽天スーパーDEAL商品を含まない。
		// 1: 楽天スーパーDEAL商品を含む。
		DealFlag int `json:"dealFlag"`

		// MembershipType はメンバーシッププログラム受注タイプです。以下のいずれかが入力されます。
		// 0: 対象外
		// 1: 楽天プレミアム対象受注
		// 2: 楽天学割対象受注
		MembershipType int `json:"membershipType"`

		// Memo はひとことメモです。
		Memo *string `json:"Memo"`

		// Operator は担当者です。
		Operator *string `json:"operator"`

		// MailPlugSentence はメール差し込み文(お客様へのメッセージ)です
		MailPlugSentence *string `json:"mailPlugSentence"`

		// ModifyFlag は購入履歴修正有無フラグです。以下のいずれかが入力されます。
		// 0: 修正なし
		// 1: 修正あり
		ModifyFlag int `json:"modifyFlag"`

		// IsTaxRecalc は消費税再計算フラグです。APIのバージョンが3以降の場合のみ入力されます。以下のいずれかが入力されます。
		// 0: 再計算しない
		// 1: 再計算する
		IsTaxRecalc int `json:"isTaxRecalc"`

		// GetOrderOrdererModel は注文者モデルです。
		GetOrderOrdererModel `json:"OrdererModel"`

		// GetOrderSettlementModel は支払い方法モデルです。
		GetOrderSettlementModel `json:"SettlementModel"`

		// GetOrderDeliveryModel は配送方法モデルです。
		GetOrderDeliveryModel `json:"DeliveryModel"`

		// GetOrderPointModel はポイントモデルです。
		GetOrderPointModel `json:"PointModel"`

		// WrappingModel1 はラッピングモデル1です。
		WrappingModel1 GetOrderWrappingModel `json:"WrappingModel"`

		// WrappingModel2 はラッピングモデル2です。
		WrappingModel2 GetOrderWrappingModel `json:"WrappingModel"`

		// PackageModelList は送付先モデルリストです。
		PackageModelList []GetOrderPackageModel `json:"PackageModelList"`

		// CouponModelList はクーポンモデルリストです。
		CouponModelList []GetOrderCouponModel `json:"CouponModelList"`

		// ChangeReasonModelList は変更・キャンセルモデルリストです。
		ChangeReasonModelList []GetOrderChangeReasonModel `json:"ChangeReasonModelList"`

		// TaxSummaryModelList は税情報モデルリストです。APIのバージョンが3以降の場合のみ入力されます。2019/7/30以前の注文の場合、空のモデルが返却されます。
		TaxSummaryModelList []GetOrderTaxSummaryModel `json:"TaxSummaryModelList"`
	}

	// GetOrderResponse は楽天ペイ受注APIの注文情報の取得で得られるレスポンスです。
	GetOrderResponse struct {
		// GetOrderMessageModelList はメッセージモデルリストです。
		GetOrderMessageModelList []GetOrderMessageModel `json:"MessageModelList"`

		// OrderModelList は受注モデルリストです。
		OrderModelList []GetOrderOrderModel `json:"OrderModelList"`
	}

	/*** updateOrderMemo ***/

	// UpdateOrderMemoMessageModel は楽天ペイ受注APIのひとことメモの更新で得られるエラー情報です。
	UpdateOrderMemoMessageModel struct {
		// CommonMessageModelResponse はエラー情報が含まれます。
		CommonMessageModelResponse

		// OrderNumber は注文番号です。
		OrderNumber string `json:"orderNumber"`
	}

	// UpdateOrderMemoResponse は楽天ペイ受注APIのひとことメモの更新で得られるレスポンスです。
	UpdateOrderMemoResponse struct {
		// MessageModelList はメッセージモデルリストです。
		MessageModelList []UpdateOrderMemoMessageModel `json:"MessageModelList"`
	}

	/*** 内部メソッド ***/

	RMSApi struct {
		serviceSecret string
		licenseKey    string
		authorization string
	}

	// SearchOrderCondition は楽天ペイ受注APIの注文検索の必須以外の検索条件です。
	SearchOrderCondition struct {
		// SortDirection は並び替え方法です。以下のいずれかを指定することができます。
		// 1: 昇順
		// 2: 降順
		SortDirection int

		// RequestRecordsAmount は1ページあたりの取得結果数です。最大1,000件まで取得可能です。
		RequestRecordsAmount int

		// RequestPage はリクエストページ番号です。
		RequestPage int

		// OrderProgressList は注文ステータスです。以下のいずれかを指定することができます。
		// 100: 注文確認待ち
		// 200: 楽天処理中
		// 300: 発送待ち
		// 400: 変更確定待ち
		// 500: 発送済み
		// 600: 支払手続き中
		// 700: 支払い手続き済み
		// 800: キャンセル確定待ち
		// 900: キャンセル確定
		OrderProgressList []int

		// SubStatusIdList はサブステータスIDリストです。
		// ユーザが指定したサブステータスを指定することができます。
		SubStatusIDList []int

		// OrderTypeList は販売種別リストです。以下のいずれかを指定することができます。
		// 1: 通常購入
		// 4: 定期購入
		// 5: 頒布会
		// 6: 予約商品
		OrderTypeList []int

		// SettlementMethod は支払い方法名です。以下のいずれかを指定することができます。
		// 1: クレジットカード
		// 2: 代金引換
		// 3: 後払い
		// 4: ショッピングクレジット/ローン
		// 5: オートローン
		// 6: リース
		// 7: 請求書払い
		// 9: 銀行振込
		// 12: Apple Pay
		// 13: セブンイレブン(前払)
		// 14: ローソン、郵便局ATM等(前払)
		// 16: Alipay
		// 17: PayPal
		// 21: 後払い決済(楽天市場の共通決済)
		SettlementMethod int

		// DeliveryName は配送方法です。
		DeliveryName string

		// ShippingDateBlankFlag は発送日未指定有無フラグです。
		ShippingDateBlankFlag bool

		// ShippingNumberBlankFlag はお荷物伝票番号未指定有無フラグです。
		ShippingNumberBlankFlag bool

		// SearchKeywordType は検索キーワード種別です。次のいずれかを指定することができます。
		// 0: なし
		// 1: 商品名
		// 2: 商品番号
		// 3: ひとことメモ
		// 4: 注文者お名前
		// 5: 注文者お名前フリガナ
		// 6: 送付先お名前
		SearchKeywordType int

		// SearchKeyword は検索キーワードです。32文字以下の入力を受け付けます。
		SearchKeyword string

		// MailSendType は注文メールアドレス種別です。以下のいずれかを指定することができます。
		// 0: PC/モバイル
		// 1: PC
		// 2: モバイル
		MailSendType int

		// OrdererMailAddress は注文者メールアドレスです。完全一致である必要があります。
		OrdererMailAddress string

		// PhoneNumberType は電話番号種別です。以下のいずれかを指定することができます。
		// 0: 注文者
		// 1: 送付先
		PhoneNumberType int

		// PhoneNumber は電話番号です。完全一致である必要があります。
		PhoneNumber string

		// ReserveNumber は申込番号です。完全一致である必要があります。
		ReserveNumber string

		// PurchaseSiteType は購入サイトリストです。以下のいずれかを指定することがあります。
		// 0: すべて
		// 1: PCで注文
		// 2: モバイルで注文
		// 3: スマートフォンで注文
		// 4: タブレットで注文
		PurchaseSiteType int

		// AsurakuFlag はあす楽希望フラグです。
		AsurakuFlag bool

		// CouponUseFlag はクーポン利用有無フラグです。
		CouponUseFlag bool

		// DrugFlag は医薬品受注フラグです。
		DrugFlag bool

		// OverseasFlag は海外カゴ注文フラグです。
		OverseasFlag bool
	}

	// UpdateOrderMemoCondition は楽天ペイ受注APIのひとことメモの更新内容です。
	UpdateOrderMemoCondition struct {
		// OrderNumber は注文番号です。
		OrderNumber string `json:"orderNumber"`

		// SubStatusID はサブステータスIDです。
		SubStatusID *int `json:"subStatusId,omitempty"`

		// DeliveryClass は配送区分です。以下のいずれかを指定することができます。
		// 0: 選択なし
		// 1: 普通
		// 2: 冷蔵
		// 3: 冷凍
		// 4: その他1
		// 5: その他2
		// 6: その他3
		// 7: その他4
		// 8: その他5
		DeliveryClass *int `json:"deliveryClass,omitempty"`

		// DeliveryDate はお届け日指定です。
		DeliveryDate *JsonDate `json:"deliveryDate,omitempty"`

		// ShippingTerm はお届け時間帯です。以下のいずれかを指定することができます。
		// 0: なし
		// 1: 午前
		// 2: 午後
		// 9: その他
		// h1h2: h1時-h2時(h1は7~24, h2は07-24まで任意の数値で指定可能)
		ShippingTerm *int `json:"shippingTerm,omitempty"`

		// Memo はひとことメモです。全角半角にかかわらず32文字以下でなければいけません。
		Memo *string `json:"memo,omitempty"`

		// Operator は担当者です。全角半角にかかわらず6文字以下でなければいけません。
		Operator *string `json:"operator,omitempty"`

		// MailPlugSentence はメール差し込み文(お客様へのメッセージ)です。全角半角にかかわらず1024文字以下でなければいけません。
		MailPlugSentence *string `json:"mailPlugSentence,omitempty"`
	}
)

func containsArray(a interface{}, v interface{}) bool {
	if reflect.TypeOf(a).Kind() != reflect.Array {
		return false
	}
	if reflect.TypeOf(a).Elem().Kind() != reflect.TypeOf(v).Kind() {
		return false
	}
	i := 0
	arr := reflect.ValueOf(a)
	val := reflect.ValueOf(v)
	for i < arr.Len() {
		if val == arr.Index(i) {
			return true
		}
		i++
	}
	return false
}

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

// Initialize はSDKを初期化します。ssはサービスシークレット、lkはライセンスキーです。サービスシークレット、ライセンスキーは https://webservice.rms.rakuten.co.jp/merchant-portal/configurationApi のページから確認してください。
func (a *RMSApi) Initialize(ss, lk string) {
	a.serviceSecret = ss
	a.licenseKey = lk
	a.authorization = base64.StdEncoding.EncodeToString([]byte(a.serviceSecret + ":" + a.licenseKey))
}

// SearchOrder は楽天ペイ受注APIで注文を検索します。注文の検索では日付を指定して検索しなければいけません。dateType は期間検索種別で、startDatetime は開始日、endDatetime は終了日です。開始日は2年以内、終了日は開始日から63日以内を指定する必要があります。
// それ以外の任意の検索条件は cond を通して指定することができます。
func (a *RMSApi) SearchOrder(dateType SearchOrderDateType, startDatetime, endDatetime time.Time, cond *SearchOrderCondition) (*SearchOrderResponse, error) {
	if a.authorization == "" {
		return nil, errors.New("Uninitialized")
	}
	reqBody := SearchOrderReuquest{}
	// For Required
	reqBody.DateType = int(dateType)
	reqBody.StartDatetime = JsonTime{startDatetime}
	reqBody.EndDatetime = JsonTime{endDatetime}
	reqBody.SearchOrderPaginationRequestModel.RequestRecordsAmount = 30
	reqBody.SearchOrderPaginationRequestModel.RequestPage = 1

	// For Optional
	if cond != nil {
		if cond.SortDirection == 2 {
			sm := SearchOrderSortModel{}
			sm.SortDirection = cond.SortDirection
			reqBody.SearchOrderPaginationRequestModel.SortModelList = append(reqBody.SearchOrderPaginationRequestModel.SortModelList, sm)
		}
		if cond.RequestRecordsAmount > 0 && cond.RequestRecordsAmount <= 1000 {
			reqBody.SearchOrderPaginationRequestModel.RequestRecordsAmount = cond.RequestRecordsAmount
		}
		if cond.RequestPage > 0 {
			reqBody.SearchOrderPaginationRequestModel.RequestPage = cond.RequestPage
		}
		if len(cond.OrderProgressList) > 0 {
			reqBody.OrderProgressList = cond.OrderProgressList
		}
		if len(cond.SubStatusIDList) > 0 {
			reqBody.SubStatusIDList = cond.SubStatusIDList
		}
		if len(cond.OrderTypeList) > 0 {
			for _, v := range cond.OrderTypeList {
				if containsArray([]int{1, 4, 5, 6}, v) {
					reqBody.OrderTypeList = append(reqBody.OrderTypeList, v)
				}
			}
		}
		if containsArray([]int{1, 2, 3, 4, 5, 6, 7, 9, 12, 13, 14, 16, 17, 21}, cond.SettlementMethod) {
			reqBody.SettlementMethod = &cond.SettlementMethod
		}
		if cond.DeliveryName != "" {
			reqBody.DeliveryName = &cond.DeliveryName
		}
		if cond.ShippingDateBlankFlag {
			reqBody.ShippingDateBlankFlag = new(int)
			*reqBody.ShippingDateBlankFlag = 1
		}
		if cond.ShippingNumberBlankFlag {
			reqBody.ShippingNumberBlankFlag = new(int)
			*reqBody.ShippingNumberBlankFlag = 1
		}
		if cond.SearchKeywordType > 0 && cond.SearchKeywordType < 7 {
			reqBody.SearchKeywordType = &cond.SearchKeywordType
			reqBody.SearchKeyword = &cond.SearchKeyword
		}
		if cond.MailSendType > 0 && cond.MailSendType < 3 {
			reqBody.MailSendType = &(cond.MailSendType)
		}
		if cond.OrdererMailAddress != "" {
			reqBody.OrdererMailAddress = &cond.OrdererMailAddress
		}
		if cond.PhoneNumberType > 0 && cond.PhoneNumberType < 2 {
			reqBody.PhoneNumberType = &cond.PhoneNumberType
		}
		if cond.PhoneNumber != "" {
			if reqBody.PhoneNumberType == nil {
				reqBody.PhoneNumberType = new(int)
				*reqBody.PhoneNumberType = 0
			}
			reqBody.PhoneNumber = &cond.PhoneNumber
		}
		if cond.ReserveNumber != "" {
			reqBody.ReserveNumber = &cond.ReserveNumber
		}
		if cond.PurchaseSiteType > 0 && cond.PurchaseSiteType < 5 {
			reqBody.PurchaseSiteType = &cond.PurchaseSiteType
		}
		if cond.AsurakuFlag {
			reqBody.AsurakuFlag = new(int)
			*reqBody.AsurakuFlag = 1
		}
		if cond.CouponUseFlag {
			reqBody.CouponUseFlag = new(int)
			*reqBody.CouponUseFlag = 1
		}
		if cond.DrugFlag {
			reqBody.DrugFlag = new(int)
			*reqBody.DrugFlag = 1
		}
		if cond.OverseasFlag {
			reqBody.OverseasFlag = new(int)
			*reqBody.OverseasFlag = 1
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

// GetOrder は楽天ペイ受注APIで注文情報を取得します。 oList は注文番号、v はバージョン番号です。バージョン番号は現在4まで指定することが可能です。
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

// UpdateOrderMemo は楽天ペイ受注APIでひとことメモを更新します。 cond は変更対象のデータです。
func (a *RMSApi) UpdateOrderMemo(cond *UpdateOrderMemoCondition) error {
	if a.authorization == "" {
		return errors.New("Uninitialized")
	}
	jsonStr, _ := json.Marshal(*cond)

	req, _ := http.NewRequest("POST", UPDATE_ORDER_MEMO_URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "ESA "+a.authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := new(http.Client)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)
	result := UpdateOrderMemoResponse{}
	err = json.Unmarshal(byteArray, &result)
	if err != nil {
		return err
	}
	if len(result.MessageModelList) == 0 {
		return errors.New("Uninitialized")
	}
	if result.MessageModelList[0].MessageType != "INFO" {
		return errors.New(result.MessageModelList[0].Message)
	}
	return nil
}
