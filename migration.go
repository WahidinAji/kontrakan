package main

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func newMigration(db *pgxpool.Pool) *Migration {
	return &Migration{
		DB: db,
	}
}

func (m *Migration) migrate(ctx context.Context) error {
	conn, err := m.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable, AccessMode: pgx.ReadWrite})
	if err != nil {
		return err
	}
	tx.Begin(ctx)

	queries := v1dot1()

	for _, q := range queries {
		_, err := tx.Exec(ctx, q)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func v1() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS reports (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			type VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			image VARCHAR(255) NOT NULL,
			user_report VARCHAR(255) NOT NULL,
			price NUMERIC NOT NULL,
			user_id BIGINT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
	}
}

func v1dot1() []string {
	return []string{
		`
		ALTER TABLE reports ALTER COLUMN image TYPE TEXT
		`,
		`
		ALTER TABLE reports ALTER COLUMN id TYPE BIGINT
		`,
	}
}
