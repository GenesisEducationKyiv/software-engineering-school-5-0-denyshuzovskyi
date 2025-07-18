package validator

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/dto"
)

type SubscriptionValidator struct {
	validator validator
}

func NewSubscriptionValidator(validator validator) *SubscriptionValidator {
	return &SubscriptionValidator{
		validator: validator,
	}
}

func (v *SubscriptionValidator) ValidateSubscriptionRequest(subscriptionReq dto.SubscriptionRequest) error {
	return v.validator.Struct(subscriptionReq)
}
