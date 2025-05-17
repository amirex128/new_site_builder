package domain

import (
	"time"
)

// Ticket represents Support.Tickets table
type Ticket struct {
	ID         int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Title      string     `json:"title" gorm:"column:title;type:longtext;not null"`
	Status     string     `json:"status" gorm:"column:status;type:longtext;not null"`
	Category   string     `json:"product_category" gorm:"column:product_category;type:longtext;not null"`
	AssignedTo *int64     `json:"assigned_to" gorm:"column:assigned_to;type:bigint;null"`
	ClosedBy   *int64     `json:"closed_by" gorm:"column:closed_by;type:bigint;null"`
	ClosedAt   *time.Time `json:"closed_at" gorm:"column:closed_at;type:datetime(6);null"`
	Priority   string     `json:"priority" gorm:"column:priority;type:longtext;not null"`
	UserID     int64      `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version    time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted  bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Comments []Comment `json:"comments" gorm:"foreignKey:TicketID"`
	Media    []Media   `json:"media" gorm:"many2many:ticket_media;"`
	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Assigned *User     `json:"assigned" gorm:"foreignKey:AssignedTo"`
	Closer   *User     `json:"closer" gorm:"foreignKey:ClosedBy"`
}

// TableName specifies the table name for Ticket
func (Ticket) TableName() string {
	return "tickets"
}

// Comment represents Support.Comments table
type Comment struct {
	ID           int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	TicketID     int64      `json:"ticket_id" gorm:"column:ticket_id;type:bigint;not null;index"`
	Content      string     `json:"content" gorm:"column:content;type:longtext;not null"`
	RespondentID int64      `json:"respondent_id" gorm:"column:respondent_id;type:bigint;not null"`
	CreatedAt    time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version      time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted    bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Ticket     *Ticket `json:"ticket" gorm:"foreignKey:TicketID"`
	Respondent *User   `json:"respondent" gorm:"foreignKey:RespondentID"`
}

// TableName specifies the table name for Comment
func (Comment) TableName() string {
	return "comments"
}

// TicketMedia represents Support.TicketMedia table - a join table
type TicketMedia struct {
	ID       int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	TicketID int64 `json:"ticket_id" gorm:"column:ticket_id;type:bigint;not null;index"`
	MediaID  int64 `json:"media_id" gorm:"column:media_id;type:bigint;not null"`
}

// TableName specifies the table name for TicketMedia
func (TicketMedia) TableName() string {
	return "ticket_media"
}

// CustomerTicket represents Support.CustomerTickets table
type CustomerTicket struct {
	ID         int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Title      string     `json:"title" gorm:"column:title;type:longtext;not null"`
	Status     string     `json:"status" gorm:"column:status;type:longtext;not null"`
	Category   string     `json:"product_category" gorm:"column:product_category;type:longtext;not null"`
	AssignedTo *int64     `json:"assigned_to" gorm:"column:assigned_to;type:bigint;null"`
	ClosedBy   *int64     `json:"closed_by" gorm:"column:closed_by;type:bigint;null"`
	ClosedAt   *time.Time `json:"closed_at" gorm:"column:closed_at;type:datetime(6);null"`
	Priority   string     `json:"priority" gorm:"column:priority;type:longtext;not null"`
	UserID     int64      `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	CustomerID int64      `json:"customer_id" gorm:"column:customer_id;type:bigint;not null"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version    time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted  bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	Comments []CustomerComment `json:"comments" gorm:"foreignKey:CustomerTicketID"`
	Media    []Media           `json:"media" gorm:"many2many:customer_ticket_media;"`
	User     *User             `json:"user" gorm:"foreignKey:UserID"`
	Customer *Customer         `json:"customer" gorm:"foreignKey:CustomerID"`
	Assigned *User             `json:"assigned" gorm:"foreignKey:AssignedTo"`
	Closer   *User             `json:"closer" gorm:"foreignKey:ClosedBy"`
}

// TableName specifies the table name for CustomerTicket
func (CustomerTicket) TableName() string {
	return "customer_tickets"
}

// CustomerComment represents Support.CustomerComments table
type CustomerComment struct {
	ID               int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	CustomerTicketID int64      `json:"customer_ticket_id" gorm:"column:customer_ticket_id;type:bigint;not null;index"`
	Content          string     `json:"content" gorm:"column:content;type:longtext;not null"`
	RespondentID     int64      `json:"respondent_id" gorm:"column:respondent_id;type:bigint;not null"`
	CreatedAt        time.Time  `json:"created_at" gorm:"column:created_at;type:datetime(6);not null"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime(6);not null"`
	Version          time.Time  `json:"version" gorm:"column:version;type:timestamp(6);default:current_timestamp(6);not null"`
	IsDeleted        bool       `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);not null"`
	DeletedAt        *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime(6);null"`

	// Relations
	CustomerTicket *CustomerTicket `json:"customer_ticket" gorm:"foreignKey:CustomerTicketID"`
	Respondent     *User           `json:"respondent" gorm:"foreignKey:RespondentID"`
}

// TableName specifies the table name for CustomerComment
func (CustomerComment) TableName() string {
	return "customer_comments"
}

// CustomerTicketMedia represents Support.CustomerTicketMedia table - a join table
type CustomerTicketMedia struct {
	ID               int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	CustomerTicketID int64 `json:"customer_ticket_id" gorm:"column:customer_ticket_id;type:bigint;not null;index"`
	MediaID          int64 `json:"media_id" gorm:"column:media_id;type:bigint;not null"`
}

// TableName specifies the table name for CustomerTicketMedia
func (CustomerTicketMedia) TableName() string {
	return "customer_ticket_media"
}
