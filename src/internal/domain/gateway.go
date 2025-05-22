package domain

import (
	"time"
)

// Gateway represents Payment.Gateways table
type Gateway struct {
	ID                                  int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	SiteID                              int64     `json:"site_id" gorm:"column:site_id;type:bigint;not null"`
	SamanMerchantId                     string    `json:"saman_merchant_id" gorm:"column:saman_merchant_id;type:longtext;null"`
	SamanPassword                       string    `json:"saman_password" gorm:"column:saman_password;type:longtext;null"`
	IsActiveSaman                       string    `json:"is_active_saman" gorm:"column:is_active_saman;type:longtext;not null"`
	MellatTerminalId                    *int64    `json:"mellat_terminal_id" gorm:"column:mellat_terminal_id;type:bigint;null"`
	MellatUserName                      string    `json:"mellat_user_name" gorm:"column:mellat_user_name;type:longtext;null"`
	MellatUserPassword                  string    `json:"mellat_user_password" gorm:"column:mellat_user_password;type:longtext;null"`
	IsActiveMellat                      string    `json:"is_active_mellat" gorm:"column:is_active_mellat;type:longtext;not null"`
	ParsianLoginAccount                 string    `json:"parsian_login_account" gorm:"column:parsian_login_account;type:longtext;null"`
	IsActiveParsian                     string    `json:"is_active_parsian" gorm:"column:is_active_parsian;type:longtext;not null"`
	PasargadMerchantCode                string    `json:"pasargad_merchant_code" gorm:"column:pasargad_merchant_code;type:longtext;null"`
	PasargadTerminalCode                string    `json:"pasargad_terminal_code" gorm:"column:pasargad_terminal_code;type:longtext;null"`
	PasargadPrivateKey                  string    `json:"pasargad_private_key" gorm:"column:pasargad_private_key;type:longtext;null"`
	IsActivePasargad                    string    `json:"is_active_pasargad" gorm:"column:is_active_pasargad;type:longtext;not null"`
	IranKishTerminalId                  string    `json:"iran_kish_terminal_id" gorm:"column:iran_kish_terminal_id;type:longtext;null"`
	IranKishAcceptorId                  string    `json:"iran_kish_acceptor_id" gorm:"column:iran_kish_acceptor_id;type:longtext;null"`
	IranKishPassPhrase                  string    `json:"iran_kish_pass_phrase" gorm:"column:iran_kish_pass_phrase;type:longtext;null"`
	IranKishPublicKey                   string    `json:"iran_kish_public_key" gorm:"column:iran_kish_public_key;type:longtext;null"`
	IsActiveIranKish                    string    `json:"is_active_iran_kish" gorm:"column:is_active_iran_kish;type:longtext;not null"`
	MelliTerminalId                     string    `json:"melli_terminal_id" gorm:"column:melli_terminal_id;type:longtext;null"`
	MelliMerchantId                     string    `json:"melli_merchant_id" gorm:"column:melli_merchant_id;type:longtext;null"`
	MelliTerminalKey                    string    `json:"melli_terminal_key" gorm:"column:melli_terminal_key;type:longtext;null"`
	IsActiveMelli                       string    `json:"is_active_melli" gorm:"column:is_active_melli;type:longtext;not null"`
	AsanPardakhtMerchantConfigurationId string    `json:"asan_pardakht_merchant_configuration_id" gorm:"column:asan_pardakht_merchant_configuration_id;type:longtext;null"`
	AsanPardakhtUserName                string    `json:"asan_pardakht_user_name" gorm:"column:asan_pardakht_user_name;type:longtext;null"`
	AsanPardakhtPassword                string    `json:"asan_pardakht_password" gorm:"column:asan_pardakht_password;type:longtext;null"`
	AsanPardakhtKey                     string    `json:"asan_pardakht_key" gorm:"column:asan_pardakht_key;type:longtext;null"`
	AsanPardakhtIV                      string    `json:"asan_pardakht_iv" gorm:"column:asan_pardakht_iv;type:longtext;null"`
	IsActiveAsanPardakht                string    `json:"is_active_asan_pardakht" gorm:"column:is_active_asan_pardakht;type:longtext;not null"`
	SepehrTerminalId                    *int64    `json:"sepehr_terminal_id" gorm:"column:sepehr_terminal_id;type:bigint;null"`
	IsActiveSepehr                      string    `json:"is_active_sepehr" gorm:"column:is_active_sepehr;type:longtext;not null"`
	ZarinPalMerchantId                  string    `json:"zarin_pal_merchant_id" gorm:"column:zarin_pal_merchant_id;type:longtext;null"`
	ZarinPalAuthorizationToken          string    `json:"zarin_pal_authorization_token" gorm:"column:zarin_pal_authorization_token;type:longtext;null"`
	ZarinPalIsSandbox                   *bool     `json:"zarin_pal_is_sandbox" gorm:"column:zarin_pal_is_sandbox;type:tinyint(1);null"`
	IsActiveZarinPal                    string    `json:"is_active_zarin_pal" gorm:"column:is_active_zarin_pal;type:longtext;not null"`
	PayIrApi                            string    `json:"pay_ir_api" gorm:"column:pay_ir_api;type:longtext;null"`
	PayIrIsTestAccount                  *bool     `json:"pay_ir_is_test_account" gorm:"column:pay_ir_is_test_account;type:tinyint(1);null"`
	IsActivePayIr                       string    `json:"is_active_pay_ir" gorm:"column:is_active_pay_ir;type:longtext;not null"`
	IdPayApi                            string    `json:"id_pay_api" gorm:"column:id_pay_api;type:longtext;null"`
	IdPayIsTestAccount                  *bool     `json:"id_pay_is_test_account" gorm:"column:id_pay_is_test_account;type:tinyint(1);null"`
	IsActiveIdPay                       string    `json:"is_active_id_pay" gorm:"column:is_active_id_pay;type:longtext;not null"`
	YekPayMerchantId                    string    `json:"yek_pay_merchant_id" gorm:"column:yek_pay_merchant_id;type:longtext;null"`
	IsActiveYekPay                      string    `json:"is_active_yek_pay" gorm:"column:is_active_yek_pay;type:longtext;not null"`
	PayPingAccessToken                  string    `json:"pay_ping_access_token" gorm:"column:pay_ping_access_token;type:longtext;null"`
	IsActivePayPing                     string    `json:"is_active_pay_ping" gorm:"column:is_active_pay_ping;type:longtext;not null"`
	IsActiveParbadVirtual               string    `json:"is_active_parbad_virtual" gorm:"column:is_active_parbad_virtual;type:longtext;not null"`
	UserID                              int64     `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt                           time.Time `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt                           time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`

	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Site *Site `json:"site" gorm:"foreignKey:SiteID"`
	User *User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for Gateway
func (Gateway) TableName() string {
	return "gateways"
}
