package validator

type LocationValidator struct {
	validator validator
}

func NewLocationValidator(validator validator) *LocationValidator {
	return &LocationValidator{
		validator: validator,
	}
}

func (v *LocationValidator) ValidateLocation(location string) error {
	return v.validator.Var(location, "required,max=60")
}
