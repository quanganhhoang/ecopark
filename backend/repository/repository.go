package repository

import "database/sql"

type Repository struct {
	Reservation ReservationRepositoryImpl
}

func New(dbConn *sql.DB) Repository {
	return Repository{
		Reservation: ReservationRepositoryImpl{dbConn},
	}
}