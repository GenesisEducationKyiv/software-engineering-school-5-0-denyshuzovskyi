package validator

type validator interface {
	Var(interface{}, string) error
	Struct(interface{}) error
}
