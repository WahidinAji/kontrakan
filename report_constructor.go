package main

import "github.com/jackc/pgx/v5/pgxpool"

func NewReportConnection(db *pgxpool.Pool) *ReportConnection {
	return &ReportConnection{
		DB: db,
	}
}
