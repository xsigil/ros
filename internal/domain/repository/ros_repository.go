package repository

import (
	"context"

	"ros/internal/domain/entity"
)

type ROSRepository interface {
	// Person
	CreatePerson(ctx context.Context, p *entity.Person) error
	GetPerson(ctx context.Context, id string) (*entity.Person, error)
	UpdatePersonStatus(ctx context.Context, id string, status entity.PersonStatus) error
	
	// Weights
	GetWeight(ctx context.Context, pType entity.PersonaType, sign entity.Sign, in entity.Inevitability) (float64, error)
	
	// Event (Append-Only)
	RecordEvent(ctx context.Context, event *entity.Event) error
	
	// View / Analytics
	GetDashboard(ctx context.Context) ([]entity.PersonScored, error)
	GetPersonScore(ctx context.Context, id string) (*entity.PersonScored, error)
}
