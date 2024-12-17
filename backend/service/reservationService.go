package service

import (
	"backend/repository"
)

type ReservationService struct {
	Repository repository.ReservationRepositoryImpl
}

func NewReservationService(repository repository.ReservationRepositoryImpl) *ReservationService {
	return &ReservationService{
		repository,
	}
}