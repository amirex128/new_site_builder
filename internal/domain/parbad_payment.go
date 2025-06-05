package domain

// ParbadPayment represents Payment.ParbadPayments table
type ParbadPayment struct {
	ID                 int64   `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	TrackingNumber     int64   `json:"tracking_number" gorm:"column:tracking_number;type:bigint;not null"`
	Amount             float64 `json:"amount" gorm:"column:amount;type:decimal(65,30);not null"`
	Token              string  `json:"token" gorm:"column:token;type:longtext;null"`
	TransactionCode    string  `json:"transaction_code" gorm:"column:transaction_code;type:longtext;null"`
	GatewayName        string  `json:"gateway_name" gorm:"column:gateway_name;type:longtext;null"`
	GatewayAccountName string  `json:"gateway_account_name" gorm:"column:gateway_account_name;type:longtext;null"`
	IsCompleted        bool    `json:"is_completed" gorm:"column:is_completed;type:tinyint(1);not null"`
	IsPaid             bool    `json:"is_paid" gorm:"column:is_paid;type:tinyint(1);not null"`

	// Relations
	Transactions []ParbadTransaction `json:"transactions" gorm:"foreignKey:PaymentID"`
}

// TableName specifies the table name for ParbadPayment
func (ParbadPayment) TableName() string {
	return "parbad_payments"
}
func (m *ParbadPayment) GetID() int64 {
	return m.ID
}
func (m *ParbadPayment) GetUserID() *int64 {
	return nil
}
func (m *ParbadPayment) GetCustomerID() *int64 {
	return nil
}
func (m *ParbadPayment) GetSiteID() *int64 {
	return nil
}
