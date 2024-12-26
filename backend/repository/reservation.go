package repository

import (
	"backend/models"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/pkg/errors"
)

type ReservationRepository interface {
	AddReservation(ctx context.Context, reservation models.Reservation) error
	FindAll(ctx context.Context) ([]models.Reservation, error)
	FindById(ctx context.Context, id string) (models.Reservation, error)
	FindAvailableDates(
		ctx context.Context,
		startDate time.Time,
		endDate time.Time,
	) ([]time.Time, error)
	IsDateRangeAvailable(
		ctx context.Context,
		querier Querier,
		startDate time.Time,
		endDate time.Time,
	) (bool, error)
}

type ReservationRepositoryImpl struct {
	DB *sql.DB
}

// Custom error
type DatesNotAvailableError struct {
	StartDate time.Time
	EndDate time.Time
}

func (e *DatesNotAvailableError) Error() string {
	return fmt.Sprintf("Dates not available between %s and %s", e.StartDate, e.EndDate)
}
// Custom error

func (reservationRepo ReservationRepositoryImpl) AddReservation(
		ctx context.Context,
		reservation models.Reservation,
) error {
	tx, err := reservationRepo.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	startDate, endDate := reservation.StartDate, reservation.EndDate
	isDateRangeAvailable, err := reservationRepo.IsDateRangeAvailable(ctx, tx, startDate, endDate)
	if err != nil || !isDateRangeAvailable {
		return &DatesNotAvailableError{startDate, endDate}
	}

	insertReservationQuery := `
		INSERT INTO reservations.reservations (
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

	_, err = tx.QueryContext(
		ctx,
		insertReservationQuery,
		reservation.Email,
		reservation.FirstName,
		reservation.LastName,
		reservation.NationalId,
		startDate,
		endDate,
		reservation.NumGuests,
	)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert reservation: %w", err)
	}

	slog.LogAttrs(
		ctx,
		slog.LevelInfo,
		"adding reservation",
		slog.String("query", insertReservationQuery),
		slog.String("email", reservation.Email),
	)

	updateCalendarQuery := `
		UPDATE calendar SET is_available = FALSE
		WHERE date BETWEEN ? AND ?
	`

	_, err = tx.ExecContext(ctx, updateCalendarQuery, startDate, endDate)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update calendar: %w", err)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (reservationRepo ReservationRepositoryImpl) FindAll(ctx context.Context) ([]models.Reservation, error) {
	query := "select * from reservations.reservations;"

	var reservations []models.Reservation
	rows, err := reservationRepo.DB.QueryContext(ctx, query)

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
		ctx,
		slog.LevelInfo,
		"fetching all reservations",
		slog.String("query", query),
	)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to execute query: %s", query))
	}
	defer rows.Close()

	return reservations, nil
}

func (reservationRepo ReservationRepositoryImpl) FindById(
	ctx context.Context,
	id string,
) (models.Reservation, error) {
	query := "SELECT * FROM reservations WHERE id = ?"
	var reservation models.Reservation
	var startDateBytes []byte
	var endDateBytes []byte
	err := reservationRepo.DB.QueryRowContext(ctx, query, id).Scan(
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

func (reservationRepo ReservationRepositoryImpl) FindAvailableDates(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]time.Time, error) {
	query := `
		select date from calendar
		WHERE true
		AND date BETWEEN ? AND ?
		AND is_available = true
	`

	rows, err := reservationRepo.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var availableDates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		availableDates = append(availableDates, date)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return availableDates, nil
}

func (reservationRepo ReservationRepositoryImpl) IsDateRangeAvailable(
	ctx context.Context,
	querier Querier,
	startDate time.Time,
	endDate time.Time,
) (bool, error) {
	query := `
		SELECT COUNT(*) AS unavailable_dates
		FROM calendar
		WHERE date BETWEEN ? AND ?
		  AND is_available = FALSE;
	`

	var unavailableCount int
	err := querier.QueryRowContext(ctx, query, startDate, endDate).Scan(&unavailableCount)
	if err != nil {
		return false, fmt.Errorf("failed to query availability: %w", err)
	}

	// If no unavailable dates, the range is fully available
	return unavailableCount == 0, nil
}
