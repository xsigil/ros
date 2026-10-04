package entity

import "time"

type Person struct {
	ID          string
	Name        string
	CurrentType PersonaType
	Status      PersonStatus
	CreatedAt   time.Time
}

type PersonScored struct {
	Person
	CurrentLogOddsDB  float64
	TotalEvents       int
	PermanentAnchorDB float64
	TransientActiveDB float64
}

func (p *PersonScored) Verdict() string {
	if p.CurrentLogOddsDB <= -10.0 {
		return "EARLY_EXIT"
	}
	if p.CurrentLogOddsDB >= 15.0 {
		return "PRIORITY_INVEST"
	}
	return "OBSERVE"
}
