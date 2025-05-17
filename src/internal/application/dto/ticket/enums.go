package ticket

// TicketStatusEnum defines the status of a ticket
type TicketStatusEnum int

const (
	New TicketStatusEnum = iota
	InProgress
	Closed
)

// TicketCategoryEnum defines the product_category of a ticket
type TicketCategoryEnum int

const (
	Bug TicketCategoryEnum = iota
	Enhancement
	FeatureRequest
	Question
	Documentation
	Financial
)

// TicketPriorityEnum defines the priority of a ticket
type TicketPriorityEnum int

const (
	Low TicketPriorityEnum = iota
	Medium
	High
	Critical
)
