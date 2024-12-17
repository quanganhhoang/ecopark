package service

import "backend/repository"

type Service struct {
	Reservation ReservationService
}

func New(repository repository.Repository) Service {
	return Service{
		Reservation: *NewReservationService(
			repository.Reservation,
		),
	}
}