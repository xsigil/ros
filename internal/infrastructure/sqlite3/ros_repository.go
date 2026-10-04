package sqlite3

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ros/internal/domain/entity"

	"github.com/jmoiron/sqlx"
)

type ROSRepository struct {
	db *sqlx.DB
}

func NewROSRepository(db *sqlx.DB) *ROSRepository {
	return &ROSRepository{db: db}
}

func (r *ROSRepository) CreatePerson(ctx context.Context, p *entity.Person) error {
	conn := GetExtContext(ctx, r.db)
	query := `INSERT INTO person (id, name, current_type, status) VALUES (?, ?, ?, ?)`
	_, err := conn.ExecContext(ctx, query, p.ID, p.Name, string(p.CurrentType), string(p.Status))
	return err
}

func (r *ROSRepository) GetPerson(ctx context.Context, id string) (*entity.Person, error) {
	conn := GetExtContext(ctx, r.db)
	query := `SELECT id, name, current_type, status, created_at FROM person WHERE id = ?`
	row := conn.QueryRowxContext(ctx, query, id)

	var p entity.Person
	var tStr, sStr string
	if err := row.Scan(&p.ID, &p.Name, &tStr, &sStr, &p.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	p.CurrentType = entity.PersonaType(tStr)
	p.Status = entity.PersonStatus(sStr)
	return &p, nil
}

func (r *ROSRepository) UpdatePersonStatus(ctx context.Context, id string, status entity.PersonStatus) error {
	conn := GetExtContext(ctx, r.db)
	query := `UPDATE person SET status = ? WHERE id = ?`
	_, err := conn.ExecContext(ctx, query, string(status), id)
	return err
}

func (r *ROSRepository) GetWeight(ctx context.Context, pType entity.PersonaType, sign entity.Sign, in entity.Inevitability) (float64, error) {
	conn := GetExtContext(ctx, r.db)
	query := `SELECT log_odds_delta FROM odds_weights WHERE type = ? AND sign = ? AND inevitability = ?`
	var weight float64
	if err := sqlx.GetContext(ctx, conn, &weight, query, string(pType), string(sign), string(in)); err != nil {
		return 0.0, fmt.Errorf("weight not found for (%s, %s, %s): %w", pType, sign, in, err)
	}
	return weight, nil
}

func (r *ROSRepository) RecordEvent(ctx context.Context, e *entity.Event) error {
	conn := GetExtContext(ctx, r.db)
	query := `INSERT INTO history (person_id, episode, type, sign, inevitability, log_odds_applied) 
              VALUES (?, ?, ?, ?, ?, ?)`
	_, err := conn.ExecContext(ctx, query, e.PersonID, e.Episode, string(e.Type), string(e.Sign), string(e.Inevitability), e.LogOddsApplied)
	return err
}

func (r *ROSRepository) GetDashboard(ctx context.Context) ([]entity.PersonScored, error) {
	conn := GetExtContext(ctx, r.db)
	query := `SELECT id, name, current_type, status, current_log_odds_db, total_events, permanent_anchor_db, transient_active_db 
              FROM v_person_effective_odds 
              ORDER BY current_log_odds_db DESC`
	rows, err := conn.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.PersonScored
	for rows.Next() {
		var item entity.PersonScored
		var tStr, sStr string
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&tStr,
			&sStr,
			&item.CurrentLogOddsDB,
			&item.TotalEvents,
			&item.PermanentAnchorDB,
			&item.TransientActiveDB,
		); err != nil {
			return nil, err
		}
		item.CurrentType = entity.PersonaType(tStr)
		item.Status = entity.PersonStatus(sStr)
		result = append(result, item)
	}
	return result, nil
}

func (r *ROSRepository) GetPersonScore(ctx context.Context, id string) (*entity.PersonScored, error) {
	conn := GetExtContext(ctx, r.db)
	query := `SELECT id, name, current_type, status, current_log_odds_db, total_events, permanent_anchor_db, transient_active_db 
              FROM v_person_effective_odds WHERE id = ?`
	row := conn.QueryRowxContext(ctx, query, id)

	var item entity.PersonScored
	var tStr, sStr string
	if err := row.Scan(
		&item.ID,
		&item.Name,
		&tStr,
		&sStr,
		&item.CurrentLogOddsDB,
		&item.TotalEvents,
		&item.PermanentAnchorDB,
		&item.TransientActiveDB,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	item.CurrentType = entity.PersonaType(tStr)
	item.Status = entity.PersonStatus(sStr)
	return &item, nil
}
