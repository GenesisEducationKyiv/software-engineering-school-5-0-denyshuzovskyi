package notification

const (
	Confirmation        = "confirmation"
	ConfirmationSuccess = "confirmation_success"
	UnsubscribeSuccess  = "unsubscribe_success"
	WeatherUpdate       = "weather_update"
)

type NotificationCommand interface {
	Type() string
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

func (c *SendConfirmation) Type() string {
	return Confirmation
}

type SendConfirmationSuccess struct {
	NotificationWithToken
}

func (c *SendConfirmationSuccess) Type() string {
	return ConfirmationSuccess
}

type SendUnsubscribeSuccess struct {
	Notification
}

func (c *SendUnsubscribeSuccess) Type() string {
	return UnsubscribeSuccess
}

type SendWeatherUpdate struct {
	NotificationWithToken
	Weather Weather `json:"weather"`
}

func (c *SendWeatherUpdate) Type() string {
	return WeatherUpdate
}
