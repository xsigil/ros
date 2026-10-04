package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"ros/internal/domain/entity"
	"ros/internal/infrastructure/sqlite3"
	"ros/internal/usecase"
	"ros/migrations"
)

// テスト用のインメモリDBとUsecaseをセットアップ
func setupTestApp(t *testing.T) (*usecase.ROSUsecase, *sqlx.DB) {
	t.Helper()
	ctx := context.Background()

	// インメモリDBを開く
	db, err := sqlx.ConnectContext(ctx, "sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory sqlite: %v", err)
	}
	db.SetMaxOpenConns(1)

	// マイグレーション適用
	if _, err := db.ExecContext(ctx, migrations.InitSQL); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	txManager := sqlite3.NewTxManager(db)
	repo := sqlite3.NewROSRepository(db)
	app := usecase.NewROSUsecase(repo, txManager)

	return app, db
}

func TestROSUsecase_AddAndCalculateScore(t *testing.T) {
	app, db := setupTestApp(t)
	defer db.Close()
	ctx := context.Background()

	// 1. ノード追加
	personID := "node_1"
	err := app.AddPerson(ctx, personID, "テスト太郎", string(entity.TypeAdultCustomer))
	if err != nil {
		t.Fatalf("AddPerson failed: %v", err)
	}

	// 2. 日常の規律（daily / pos）を2回記録 (+2.0 dB x 2 = +4.0 dB)
	if err := app.RecordObservation(ctx, personID, "迅速な返信", "pos", "daily"); err != nil {
		t.Fatalf("RecordObservation 1 failed: %v", err)
	}
	if err := app.RecordObservation(ctx, personID, "納期前納品", "pos", "daily"); err != nil {
		t.Fatalf("RecordObservation 2 failed: %v", err)
	}

	// 3. ダッシュボードのスコア検証
	dash, err := app.Dashboard(ctx)
	if err != nil {
		t.Fatalf("Dashboard failed: %v", err)
	}
	if len(dash) != 1 {
		t.Fatalf("expected 1 person, got %d", len(dash))
	}

	p := dash[0]
	if p.TotalEvents != 2 {
		t.Errorf("expected 2 events, got %d", p.TotalEvents)
	}
	// 直後の記録なので減衰はほぼゼロ、+4.0 dB 前後であること
	if p.CurrentLogOddsDB < 3.9 || p.CurrentLogOddsDB > 4.1 {
		t.Errorf("expected score ~4.0 dB, got %.2f dB", p.CurrentLogOddsDB)
	}
	if p.Status != entity.StatusActive {
		t.Errorf("expected status 'active', got '%s'", p.Status)
	}
}

func TestROSUsecase_EarlyExitTrigger(t *testing.T) {
	app, db := setupTestApp(t)
	defer db.Close()
	ctx := context.Background()

	personID := "node_bad"
	// 大人・非顧客 (daily neg は -4.0 dB)
	err := app.AddPerson(ctx, personID, "時間泥棒氏", string(entity.TypeAdultNonCustomer))
	if err != nil {
		t.Fatalf("AddPerson failed: %v", err)
	}

	// 3回連続で daily neg を記録 (-4.0 x 3 = -12.0 dB)
	for i := 1; i <= 3; i++ {
		err := app.RecordObservation(ctx, personID, "無断遅刻・放置", "neg", "daily")
		if err != nil {
			t.Fatalf("RecordObservation %d failed: %v", i, err)
		}
	}

	// ダッシュボードで受動的監視 (monitoring_only) に降格しているか確認
	dash, err := app.Dashboard(ctx)
	if err != nil {
		t.Fatalf("Dashboard failed: %v", err)
	}

	p := dash[0]
	if p.CurrentLogOddsDB > -10.0 {
		t.Errorf("expected score <= -10.0 dB, got %.2f dB", p.CurrentLogOddsDB)
	}
	if p.Status != entity.StatusMonitoringOnly {
		t.Errorf("expected status 'monitoring_only' (Early Exit), got '%s'", p.Status)
	}
	if p.Verdict() != "EARLY_EXIT" {
		t.Errorf("expected verdict 'EARLY_EXIT', got '%s'", p.Verdict())
	}
}

func TestROSUsecase_CriticalMomentPermanence(t *testing.T) {
	app, db := setupTestApp(t)
	defer db.Close()
	ctx := context.Background()

	personID := "node_hero"
	err := app.AddPerson(ctx, personID, "信義の人", string(entity.TypeEmployee))
	if err != nil {
		t.Fatalf("AddPerson failed: %v", err)
	}

	// 窮地での信義 (employee / critical_moment / pos: +20.0 dB)
	if err := app.RecordObservation(ctx, personID, "大赤字案件で自腹を切って完遂", "pos", "critical_moment"); err != nil {
		t.Fatalf("RecordObservation failed: %v", err)
	}

	// タイムスタンプを過去（365日前）に直接巻き戻して擬似的に1年経過させる
	pastTime := time.Now().AddDate(-1, 0, 0).Format("2006-01-02 15:04:05")
	_, err = db.ExecContext(ctx, "UPDATE history SET timestamp = ? WHERE person_id = ?", pastTime, personID)
	if err != nil {
		t.Fatalf("failed to rewind history timestamp: %v", err)
	}

	// 集計結果を確認
	dash, err := app.Dashboard(ctx)
	if err != nil {
		t.Fatalf("Dashboard failed: %v", err)
	}

	p := dash[0]
	// critical_moment は半減期無限大（NULL）のため、1年経っても 20.0 dB のまま減衰しないこと
	if p.PermanentAnchorDB != 20.0 {
		t.Errorf("expected permanent anchor 20.0 dB, got %.2f dB", p.PermanentAnchorDB)
	}
	if p.CurrentLogOddsDB != 20.0 {
		t.Errorf("expected effective odds to remain 20.0 dB after 1 year, got %.2f dB", p.CurrentLogOddsDB)
	}
	if p.Verdict() != "PRIORITY_INVEST" {
		t.Errorf("expected verdict 'PRIORITY_INVEST', got '%s'", p.Verdict())
	}
}
