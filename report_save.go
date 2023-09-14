package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *ReportConnection) Save(ctx context.Context, in Report) (res Report, err error) {

	conn, err := r.DB.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable, AccessMode: pgx.ReadWrite})
	if err != nil {
		return
	}

	err = tx.QueryRow(ctx, `INSERT INTO reports (title, type, description, image, user_report, price) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, title, type, description, image, user_report, price`, in.Title, in.Type, in.Description, in.Image, in.UserReport, in.Price).Scan(&res.Id, &res.Title, &res.Type, &res.Description, &res.Image, &res.UserReport, &res.Price)
	if err != nil {
		errRoll := tx.Rollback(ctx)
		if errRoll != nil {
			err = errRoll
			return
		}
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		return
	}
	return
}
