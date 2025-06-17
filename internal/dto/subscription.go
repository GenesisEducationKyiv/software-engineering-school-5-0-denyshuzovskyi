package dto

type SubscriptionRequest struct {
	Email     string `validate:"required,email"`
	City      string `validate:"required"`
	Frequency string `validate:"required,oneof=hourly daily"`
}
