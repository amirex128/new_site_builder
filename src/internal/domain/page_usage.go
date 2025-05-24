package domain

// PageArticleUsage represents Site.PageArticleUsages table - a join table
type PageArticleUsage struct {
	ID        int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	PageID    int64 `json:"page_id" gorm:"column:page_id;type:bigint;not null"`
	ArticleID int64 `json:"article_id" gorm:"column:article_id;type:bigint;not null;index"`
	SiteID    int64 `json:"site_id" gorm:"column:site_id;type:bigint;not null;index"`
	UserID    int64 `json:"user_id" gorm:"column:user_id;type:bigint;not null"`

	// Relations
	Page    *Page    `json:"page" gorm:"foreignKey:PageID"`
	Article *Article `json:"article" gorm:"foreignKey:ArticleID"`
	Site    *Site    `json:"site" gorm:"foreignKey:SiteID"`
	User    *User    `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for PageArticleUsage
func (PageArticleUsage) TableName() string {
	return "page_article_usages"
}
func (m *PageArticleUsage) GetID() int64 {
	return m.ID
}
func (m *PageArticleUsage) GetUserID() *int64 {
	return &m.UserID
}
func (m *PageArticleUsage) GetCustomerID() *int64 {
	return nil
}
func (m *PageArticleUsage) GetSiteID() *int64 {
	return &m.SiteID
}

// PageProductUsage represents Site.PageProductUsages table - a join table
type PageProductUsage struct {
	ID        int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	PageID    int64 `json:"page_id" gorm:"column:page_id;type:bigint;not null"`
	ProductID int64 `json:"product_id" gorm:"column:product_id;type:bigint;not null;index"`
	SiteID    int64 `json:"site_id" gorm:"column:site_id;type:bigint;not null;index"`
	UserID    int64 `json:"user_id" gorm:"column:user_id;type:bigint;not null"`

	// Relations
	Page    *Page    `json:"page" gorm:"foreignKey:PageID"`
	Product *Product `json:"article" gorm:"foreignKey:ProductID"`
	Site    *Site    `json:"site" gorm:"foreignKey:SiteID"`
	User    *User    `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for PageProductUsage
func (PageProductUsage) TableName() string {
	return "page_product_usages"
}
func (m *PageProductUsage) GetID() int64 {
	return m.ID
}
func (m *PageProductUsage) GetUserID() *int64 {
	return &m.UserID
}
func (m *PageProductUsage) GetCustomerID() *int64 {
	return nil
}
func (m *PageProductUsage) GetSiteID() *int64 {
	return &m.SiteID
}

// PageHeaderFooterUsage represents Site.PageHeaderFooterUsages table - a join table
type PageHeaderFooterUsage struct {
	ID             int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	PageID         int64 `json:"page_id" gorm:"column:page_id;type:bigint;not null"`
	HeaderFooterID int64 `json:"header_footer_id" gorm:"column:header_footer_id;type:bigint;not null;index"`
	SiteID         int64 `json:"site_id" gorm:"column:site_id;type:bigint;not null;index"`
	UserID         int64 `json:"user_id" gorm:"column:user_id;type:bigint;not null"`

	// Relations
	Page         *Page         `json:"page" gorm:"foreignKey:PageID"`
	HeaderFooter *HeaderFooter `json:"header_footer" gorm:"foreignKey:HeaderFooterID"`
	Site         *Site         `json:"site" gorm:"foreignKey:SiteID"`
	User         *User         `json:"user" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for PageHeaderFooterUsage
func (PageHeaderFooterUsage) TableName() string {
	return "page_header_footer_usages"
}
func (m *PageHeaderFooterUsage) GetID() int64 {
	return m.ID
}
func (m *PageHeaderFooterUsage) GetUserID() *int64 {
	return &m.UserID
}
func (m *PageHeaderFooterUsage) GetCustomerID() *int64 {
	return nil
}
func (m *PageHeaderFooterUsage) GetSiteID() *int64 {
	return &m.SiteID
}
