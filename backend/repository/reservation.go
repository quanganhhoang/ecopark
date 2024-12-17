package repository

import (
	"backend/models"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type ReservationRepository interface {
	AddReservation(reservation models.Reservation) error
	FindAll() error
}

type ReservationRepositoryImpl struct {
	DB *sql.DB
}

func (reservationRepo ReservationRepositoryImpl) AddReservation(reservation models.Reservation) error {
	query := `
		INSERT INTO reservations (
			email,
			first_name,
			last_name,
			national_id,
			start_date,
			end_date,
			num_guests
		)
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`

	_, err := reservationRepo.DB.Exec(
		query,
		reservation.Email,
		reservation.FirstName,
		reservation.LastName,
		reservation.NationalId,
		reservation.StartDate,
		reservation.EndDate,
		reservation.NumGuests,
	)

	slog.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"adding reservation",
		slog.String("query", query),
		slog.String("email", reservation.Email),
	)

	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	return nil
}

func (reservationRepo ReservationRepositoryImpl) FindAll() ([]models.Reservation, error) {
	query := "select * from reservations.reservations;"

	var reservations []models.Reservation
	rows, err := reservationRepo.DB.Query(query)

	for rows.Next() {
		var reservation models.Reservation
		var startDateBytes []byte
		var endDateBytes []byte
		if err := rows.Scan(
				&reservation.ID,
				&reservation.Email,
				&reservation.FirstName,
				&reservation.LastName,
				&reservation.NationalId,
				&startDateBytes,
				&endDateBytes,
				&reservation.NumGuests,
			); err != nil {
			slog.Info("findAll", "error", err)
		}
		reservation.StartDate, _ = time.Parse("YYYY-MM-dd", string(startDateBytes))
		reservation.EndDate, _ = time.Parse("YYYY-MM-dd", string(endDateBytes))
		reservations = append(reservations, reservation)
	}

	slog.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"fetching all reservations",
		slog.String("query", query),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	return reservations, nil
}

func (reservationRepo ReservationRepositoryImpl) FindById(id string) (models.Reservation, error) {
	query := "SELECT * FROM reservations WHERE id = ?"
	var reservation models.Reservation
	var startDateBytes []byte
	var endDateBytes []byte
	err := reservationRepo.DB.QueryRow(query, id).Scan(
		&reservation.ID,
		&reservation.Email,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.NationalId,
		&startDateBytes,
		&endDateBytes,
		&reservation.NumGuests,
	)
	reservation.StartDate, _ = time.Parse("YYYY-MM-dd", string(startDateBytes))
	reservation.EndDate, _ = time.Parse("YYYY-MM-dd", string(endDateBytes))

	return reservation, err
}