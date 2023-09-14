package main

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func TestMain(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		t.Skip("Skip test")
	})
}

func TestMigration(t *testing.T) {
	ctx, can := context.WithTimeout(context.Background(), 10*time.Second)
	defer can()

	dbString := "postgresql://postgres:password@localhost:5432/kontrakan?sslmode=disable"

	pool, err := pgxpool.New(ctx, dbString)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database")
	}
	defer pool.Close()

	m := newMigration(pool)

	t.Run("migrate", func(t *testing.T) {
		err = m.migrate(ctx)

		if err != nil {
			t.Errorf("Error when migrate: %v", err)
		}
	})
}
