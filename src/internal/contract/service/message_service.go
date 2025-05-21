package service

type IMessageService interface {
	SendEmail(msg any) error
	SendSms(msg any) error
}
