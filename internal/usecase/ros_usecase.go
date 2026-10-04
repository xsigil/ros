package usecase

import (
	"context"
	"fmt"

	"ros/internal/domain/entity"
	"ros/internal/domain/repository"
)

type ROSUsecase struct {
	repo      repository.ROSRepository
	txManager repository.TxManager
}

func NewROSUsecase(repo repository.ROSRepository, txManager repository.TxManager) *ROSUsecase {
	return &ROSUsecase{
		repo:      repo,
		txManager: txManager,
	}
}

func (u *ROSUsecase) AddPerson(ctx context.Context, id, name, rawType string) error {
	pType, err := entity.ValidatePersonaType(rawType)
	if err != nil {
		return err
	}

	p := &entity.Person{
		ID:          id,
		Name:        name,
		CurrentType: pType,
		Status:      entity.StatusActive,
	}

	return u.txManager.Do(ctx, func(txCtx context.Context) error {
		return u.repo.CreatePerson(txCtx, p)
	})
}

func (u *ROSUsecase) RecordObservation(ctx context.Context, personID, episode, rawSign, rawInevitability string) error {
	sign, err := entity.ValidateSign(rawSign)
	if err != nil {
		return err
	}
	inevitability, err := entity.ValidateInevitability(rawInevitability)
	if err != nil {
		return err
	}

	return u.txManager.Do(ctx, func(txCtx context.Context) error {
		person, err := u.repo.GetPerson(txCtx, personID)
		if err != nil {
			return err
		}
		if person == nil {
			return fmt.Errorf("person '%s' not found", personID)
		}

		// オッズ重みをマスタから取得
		weight, err := u.repo.GetWeight(txCtx, person.CurrentType, sign, inevitability)
		if err != nil {
			return err
		}

		event := &entity.Event{
			PersonID:       personID,
			Episode:        episode,
			Type:           person.CurrentType,
			Sign:           sign,
			Inevitability:  inevitability,
			LogOddsApplied: weight,
		}

		if err := u.repo.RecordEvent(txCtx, event); err != nil {
			return err
		}

		// 更新後のスコアを評価し、Early Exit 閾値（<= -10.0 dB）なら受動的監視へ自動移行
		score, err := u.repo.GetPersonScore(txCtx, personID)
		if err != nil {
			return err
		}
		if score != nil && score.CurrentLogOddsDB <= -10.0 && person.Status == entity.StatusActive {
			if err := u.repo.UpdatePersonStatus(txCtx, personID, entity.StatusMonitoringOnly); err != nil {
				return err
			}
			fmt.Printf("[!] Person %s breached early-exit threshold (%.2f dB) -> Status set to monitoring_only\n", person.Name, score.CurrentLogOddsDB)
		}

		return nil
	})
}

func (u *ROSUsecase) Dashboard(ctx context.Context) ([]entity.PersonScored, error) {
	return u.repo.GetDashboard(ctx)
}
