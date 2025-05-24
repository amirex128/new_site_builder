package domain

// Media is a reference model for all tables with media relationships
type Media struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`

	// Relations
	Articles          []Article         `json:"articles" gorm:"many2many:article_media;"`
	BlogCategories    []ArticleCategory `json:"blog_categories" gorm:"many2many:category_media;"`
	ProductCategories []ProductCategory `json:"product_categories" gorm:"many2many:category_media;"`
	Products          []Product         `json:"products" gorm:"many2many:product_media;"`
	Pages             []Page            `json:"pages" gorm:"many2many:page_media;"`
	Tickets           []Ticket          `json:"tickets" gorm:"many2many:ticket_media;"`
	CustomerTickets   []CustomerTicket  `json:"customer_tickets" gorm:"many2many:customer_ticket_media;"`
	DefaultThemes     []DefaultTheme    `json:"default_themes" gorm:"foreignKey:MediaID"`
}

// TableName specifies the table name for Media
func (Media) TableName() string {
	return "media"
}
func (m *Media) GetID() int64 {
	return m.ID
}
func (m *Media) GetUserID() *int64 {
	return nil
}
func (m *Media) GetCutomerID() *int64 {
	return nil
}
func (m *Media) GetSiteID() *int64 {
	return nil
}
