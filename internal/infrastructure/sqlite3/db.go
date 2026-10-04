package sqlite3

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"ros/migrations"
)

func NewDB(ctx context.Context, dbPath string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite: %w", err)
	}

	// SQLiteの書き込みロック競合を防ぐ設定
	db.SetMaxOpenConns(1)

	// マイグレーションSQLの適用
	if _, err := db.ExecContext(ctx, migrations.InitSQL); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return db, nil
}
