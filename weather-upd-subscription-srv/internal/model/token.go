package model

import "time"

type TokenType string

const (
	TokenType_Confirmation TokenType = "confirmation"
	TokenType_Unsubscribe  TokenType = "unsubscribe"
)

type Token struct {
	Token          string
	SubscriptionId int32
	Type           TokenType
	CreatedAt      time.Time
	ExpiresAt      time.Time
	UsedAt         time.Time
}
