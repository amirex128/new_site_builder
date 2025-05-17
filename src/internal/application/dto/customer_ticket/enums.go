package customer_ticket

// CustomerTicketStatusEnum defines the status of a customer ticket
type CustomerTicketStatusEnum int

const (
	New CustomerTicketStatusEnum = iota
	InProgress
	Closed
)

// CustomerTicketCategoryEnum defines the product_category of a customer ticket
type CustomerTicketCategoryEnum int

const (
	Bug CustomerTicketCategoryEnum = iota
	Enhancement
	FeatureRequest
	Question
	Documentation
	Financial
)

// CustomerTicketPriorityEnum defines the priority of a customer ticket
type CustomerTicketPriorityEnum int

const (
	Low CustomerTicketPriorityEnum = iota
	Medium
	High
	Critical
)
