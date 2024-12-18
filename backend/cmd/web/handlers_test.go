package main

import (
	"backend/models"
	"backend/service"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository to simulate FindById behavior
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindById(id string) (models.Reservation, error) {
	args := m.Called(id)
	return args.Get(0).(models.Reservation), args.Error(1)
}

func (m *MockRepository) FindAll() ([]models.Reservation, error) {
	args := m.Called()
	return args.Get(0).([]models.Reservation), args.Error(1)
}

func (m *MockRepository) AddReservation(r models.Reservation) error {
	args := m.Called(r)
	return args.Error(1)
}

// TestHandleGetReservationByID tests the handler for different scenarios
func TestHandleGetReservationByID(t *testing.T) {
	// Mock data
	mockID := "123"
	mockFirstName := "John"
	mockLastName := "Doe"
	mockEmail := "john.doe@example.com"
	mockNumGuests := 2
	mockStartDate := time.Date(2024, time.December, 25, 10, 30, 0, 0, time.UTC)
	mockEndDate := time.Date(2024, time.December, 31, 10, 30, 0, 0, time.UTC)
	mockNationalId := "ABC123456"
	mockReservation := models.Reservation{
		ID:         mockID,
		FirstName:  mockFirstName,
		LastName:   mockLastName,
		Email:      mockEmail,
		NumGuests:  mockNumGuests,
		StartDate:  mockStartDate,
		EndDate:    mockEndDate,
		NationalId: mockNationalId,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)

	t.Run("Reservation Found", func(t *testing.T) {
		// Initialize mocks
		mockRepo := new(MockRepository)
		mockRepo.On("FindById", mockID).Return(mockReservation, nil)

		app := &App{
			Service: service.Service{
				Reservation: service.ReservationService{
					Repository: mockRepo,
				},
			},
		}

		// Create a test context
		router := gin.Default()
		router.GET("/reservations/:id", app.HandleGetReservationByID)

		// Simulate request
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/reservations/%s", mockID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertCalled(t, "FindById", mockID)
		assert.Contains(t, w.Body.String(), mockFirstName)
		assert.Contains(t, w.Body.String(), mockLastName)
	})

	t.Run("Reservation Not Found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("FindById", mockID).Return(models.Reservation{}, sql.ErrNoRows)

		app := &App{
			Service: service.Service{
				Reservation: service.ReservationService{
					Repository: mockRepo,
				},
			},
		}

		router := gin.Default()
		router.GET("/reservations/:id", app.HandleGetReservationByID)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/reservations/%s", mockID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertCalled(t, "FindById", mockID)
		assert.Contains(t, w.Body.String(), "Reservation not found")
	})

	t.Run("Server Error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("FindById", mockID).Return(models.Reservation{}, errors.New("database error"))

		app := &App{
			Service: service.Service{
				Reservation: service.ReservationService{
					Repository: mockRepo,
				},
			},
		}

		router := gin.Default()
		router.GET("/reservations/:id", app.HandleGetReservationByID)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/reservations/%s", mockID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertCalled(t, "FindById", mockID)
		assert.Contains(t, w.Body.String(), "Server error")
	})
}