package repository

import (
	"context"
	"database/sql"
)

type Repository struct {
	Reservation ReservationRepositoryImpl
}

func New(dbConn *sql.DB) Repository {
	return Repository{
		Reservation: ReservationRepositoryImpl{dbConn},
	}
}

// Interface used to unify *sql.DB
// and *sql.Tx for when a query needs to be run in a transaction
type Querier interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}