package model

import "time"

type Frequency string

const (
	Frequency_Hourly Frequency = "hourly"
	Frequency_Daily  Frequency = "daily"
)

type SubscriptionStatus string

const (
	SubscriptionStatus_Pending   SubscriptionStatus = "pending"
	SubscriptionStatus_Confirmed SubscriptionStatus = "confirmed"
)

type Subscriber struct {
	Id        int32
	Email     string
	CreatedAt time.Time
}

type Subscription struct {
	Id           int32
	SubscriberId int32
	LocationId   int32
	Frequency    Frequency
	Status       SubscriptionStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
