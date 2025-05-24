package domain

// ParbadTransaction represents Payment.ParbadTransactions table
type ParbadTransaction struct {
	ID             int64   `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Amount         float64 `json:"amount" gorm:"column:amount;type:decimal(65,30);not null"`
	Type           uint8   `json:"type" gorm:"column:type;type:tinyint unsigned;not null"`
	IsSucceed      bool    `json:"is_succeed" gorm:"column:is_succeed;type:tinyint(1);not null"`
	Message        string  `json:"message" gorm:"column:message;type:longtext;null"`
	AdditionalData string  `json:"additional_data" gorm:"column:additional_data;type:longtext;null"`
	PaymentID      int64   `json:"payment_id" gorm:"column:payment_id;type:bigint;not null"`

	// Relations
	ParbadPayment *ParbadPayment `json:"parbad_payment" gorm:"foreignKey:PaymentID"`
}

// TableName specifies the table name for ParbadTransaction
func (ParbadTransaction) TableName() string {
	return "parbad_transactions"
}
func (m *ParbadTransaction) GetID() int64 {
	return m.ID
}
func (m *ParbadTransaction) GetUserID() *int64 {
	return nil
}
func (m *ParbadTransaction) GetCustomerID() *int64 {
	return nil
}
func (m *ParbadTransaction) GetSiteID() *int64 {
	return nil
}
