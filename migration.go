package main

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func newMigration(db *pgxpool.Pool) *Migration {
	return &Migration{
		DB: db,
	}
}

func (m Migration) migrate(ctx context.Context) error {
	conn, err := m.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable, AccessMode: pgx.ReadWrite})
	if err != nil {
		return err
	}

	//TODO: if the v1dot2 is already exist, then make thhe logic to check the last migration

	// #region check if the migration table is exist
	// var checkMigrationTable bool
	// err = tx.QueryRow(ctx, `SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'migtations')`).Scan(&checkMigrationTable)
	// if err != nil {
	// 	tx.Rollback(ctx)
	// 	return err
	// }
	// var migrateVersion string
	// if checkMigrationTable {
	// 	migrateVersion = "SELECT version FROM migtations ORDER BY id DESC LIMIT 1"
	// }
	// if migrateVersion == "" {
	// 	migrateVersion = "v1"
	// }
	// var version string
	// err = tx.QueryRow(ctx, migrateVersion).Scan(&version)
	// if err != nil {
	// 	if err == p
	// 	if err == pgx.ErrNoRows {
	// 		version = "v1"
	// 	} else {
	// 		tx.Rollback(ctx)
	// 		return err
	// 	}
	// }
	// log.Info().Msgf("Current version: %s", version)
	// var queries []string
	// switch version {
	// case "v1":
	// 	queries = v1()
	// case "v1dot1":
	// 	queries = v1dot1()
	// case "v1dot2":
	//#endregion

	queries := v1dot2()
	// default:
	// 	queries = v1()

	for i, q := range queries {
		log.Info().Msgf("Migration %s, %d", q, i)
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

func v1dot2() []string {
	return []string{
		`
		ALTER TABLE reports ALTER COLUMN image DROP NOT NULL
		`,
	}
}

// func seedVersion() []string {
// 	return []string{
// 		`v1dot2`,
// 		`CREATE TABLE IF NOT EXISTS migtations (
// 			id SERIAL PRIMARY KEY,
// 			version VARCHAR(255) NOT NULL,
// 			description TEXT NOT NULL,
// 			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
// 			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
// 		)`,
// 		`INSERT INTO migtations (version,description) VALUES ($1,$2)`, "v1", "make image nullable and make a new table called migratins to track the migrations",
// 		`INSERT INTO migtations (version,description) VALUES ($1,$2)`, "v1dot1", "update column image to text and make the id bigint",
// 		`INSERT INTO migtations (version,description) VALUES ($1,$2)`, "v1dot2", "make image nullable and make a new table called migratins to track the migrations",
// 	}
// }
