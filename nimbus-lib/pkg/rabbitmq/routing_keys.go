package rabbitmq

type RoutingKey string

const (
	SendConfirmationKey        RoutingKey = "notification.send.confirmation"
	SendConfirmationSuccessKey RoutingKey = "notification.send.confirmation_success"
	SendUnsubscribeSuccessKey  RoutingKey = "notification.send.unsubscribe_success"
	SendWeatherUpdateKey       RoutingKey = "notification.send.weather_update"
)
