package service

import (
	"backend/repository"
)

type ReservationService struct {
	Repository repository.ReservationRepository
}

func NewReservationService(repository repository.ReservationRepository) *ReservationService {
	return &ReservationService{
		repository,
	}
}