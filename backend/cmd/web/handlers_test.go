package main

import (
	"backend/models"
	"backend/service"
	"database/sql"
	"encoding/json"
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

func generateMockData() []models.Reservation {
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

	return []models.Reservation{mockReservation}
}

// TestHandleGetReservationByID tests the handler for different scenarios
func TestHandleGetReservationByID(t *testing.T) {
	mockReservation := generateMockData()[0]
	mockID := mockReservation.ID

	// Setup Gin
	gin.SetMode(gin.TestMode)

	t.Run("Reservation Found", func(t *testing.T) {
		// Initialize mocks
		mockRepo := new(MockRepository)
		mockRepo.On("FindById", mockReservation).Return(mockReservation, nil)

		app := &App{
			Service: service.Service{
				Reservation: service.ReservationService{
					Repository: mockRepo,
				},
			},
		}

		// Create a test context
		router := gin.Default()
		router.GET("/api/reservations/:id", app.HandleGetReservationByID)

		// Simulate request
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/reservations/%s", mockID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertCalled(t, "FindById", mockID)
		assert.Contains(t, w.Body.String(), mockReservation.FirstName)
		assert.Contains(t, w.Body.String(), mockReservation.LastName)
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
		router.GET("/api/reservations/:id", app.HandleGetReservationByID)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/reservations/%s", mockID), nil)
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
		router.GET("/api/reservations/:id", app.HandleGetReservationByID)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/reservations/%s", mockID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertCalled(t, "FindById", mockID)
		assert.Contains(t, w.Body.String(), "Server error")
	})
}

func TestHandleGetReservations(t *testing.T) {
	mockReservations := generateMockData()
	gin.SetMode(gin.TestMode)

	t.Run("Fetch all reservations", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("FindAll").Return(mockReservations, nil)

		app := &App{
			Service: service.Service{
				Reservation: service.ReservationService{
					Repository: mockRepo,
				},
			},
		}

		// Create a test context
		router := gin.Default()
		router.GET("/api/reservations", app.HandleGetReservations)

		// Simulate request
		req, _ := http.NewRequest(http.MethodGet, "/api/reservations", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertCalled(t, "FindAll")

		var reservations []models.Reservation

		var rawMap map[string]json.RawMessage
		err := json.Unmarshal(w.Body.Bytes(), &rawMap)
		if err != nil {
			fmt.Printf("Failed to unmarshal into map: %v\n", err)
			return
		}
		err = json.Unmarshal(rawMap["reservations"], &reservations)

		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v. Response: %s", err, rawMap["reservations"])
		}

		// Assert that there is exactly one reservation
		if len(reservations) != 1 {
			t.Fatalf("Expected 1 reservation, got %d", len(reservations))
		}

	})

	t.Run("Server Error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("FindAll").Return([]models.Reservation{}, errors.New("database error"))

		app := &App{
			Service: service.Service{
				Reservation: service.ReservationService{
					Repository: mockRepo,
				},
			},
		}

		router := gin.Default()
		router.GET("/api/reservations", app.HandleGetReservations)

		req, _ := http.NewRequest(http.MethodGet, "/api/reservations", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRepo.AssertCalled(t, "FindAll")
	})
}