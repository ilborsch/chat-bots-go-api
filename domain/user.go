package domain

import (
	"time"
)

const (
	FreePlan       = "free_plan"
	BusinessPlan   = "business_plan"
	EnterprisePlan = "enterprise_plan"

	FreeMessagesAmount       = 20
	BusinessMessagesAmount   = 2500
	EnterpriseMessagesAmount = 10000

	FreeBytesDataAmount       = 1024 * 1024 * 1   // 1 MB
	BusinessBytesDataAmount   = 1024 * 1024 * 100 // 100 MB
	EnterpriseBytesDataAmount = 1024 * 1024 * 500 // 500 MB

	FreeBotsAmount       = 1
	BusinessBotsAmount   = 3
	EnterpriseBotsAmount = 10
)

type User struct {
	ID             int64
	Email          string
	Plan           string
	PlanBoughtDate time.Time
	MessagesLeft   int
	BytesDataLeft  int
	BotsLeft       int
}
