package entity

import "time"

type Event struct {
	ID             int64
	PersonID       string
	Episode        string
	Type           PersonaType
	Sign           Sign
	Inevitability  Inevitability
	LogOddsApplied float64
	Timestamp      time.Time
}
