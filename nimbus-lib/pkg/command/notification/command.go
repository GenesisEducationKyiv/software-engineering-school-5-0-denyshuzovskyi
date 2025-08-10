package notification

type CommandType string

const (
	Confirmation        CommandType = "confirmation"
	ConfirmationSuccess CommandType = "confirmation_success"
	UnsubscribeSuccess  CommandType = "unsubscribe_success"
	WeatherUpdate       CommandType = "weather_update"
)

type NotificationCommand interface {
	Type() CommandType
}

type Notification struct {
	To string `json:"to"`
}

type NotificationWithToken struct {
	Notification
	Token string `json:"token"`
}

type Weather struct {
	Location    string  `json:"location"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	Description string  `json:"description"`
}

type SendConfirmation struct {
	NotificationWithToken
}

func (c *SendConfirmation) Type() CommandType {
	return Confirmation
}

type SendConfirmationSuccess struct {
	NotificationWithToken
}

func (c *SendConfirmationSuccess) Type() CommandType {
	return ConfirmationSuccess
}

type SendUnsubscribeSuccess struct {
	Notification
}

func (c *SendUnsubscribeSuccess) Type() CommandType {
	return UnsubscribeSuccess
}

type SendWeatherUpdate struct {
	NotificationWithToken
	Weather Weather `json:"weather"`
}

func (c *SendWeatherUpdate) Type() CommandType {
	return WeatherUpdate
}
