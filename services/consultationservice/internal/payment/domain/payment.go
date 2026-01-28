package domain

import "time"

type PaymentStatus string
type PaymentMethod string

const (
	PaymentPending  PaymentStatus = "pending"
	PaymentPaid     PaymentStatus = "paid"
	PaymentFailed   PaymentStatus = "failed"
	PaymentRefunded PaymentStatus = "refunded"
)

type Payment struct {
	ID             string
	SessionID      string
	Amount         float64
	Currency       string
	Provider       string
	Status         PaymentStatus
	TransactionRef string
	CreatedAt      time.Time
	PaidAt         *time.Time
}
