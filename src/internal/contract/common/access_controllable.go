package common

type AccessControllable interface {
	GetCustomerID() *int64
	GetUserID() *int64
	GetSiteID() *int64
	GetID() *int64
}
