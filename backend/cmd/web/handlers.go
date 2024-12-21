package main

import (
	"backend/models"
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

func (app *App) HandleGetReservations(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "Invalid request method",
		})

		return
	}

	reservations, err := app.Service.Reservation.Repository.FindAll()
	if err != nil {
		slog.Error("Query failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"reservations": reservations,
		})
	}
}

// func (app *App) HandlePostReservation(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		log.Fatalf("Failed to read response: %v", err)
// 	}

// 	var reservation models.Reservation
// 	if err := json.Unmarshal(body, &reservation); err != nil {
// 		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
// 		return
// 	}

// 	err = app.Database.DB.AddReservation(reservation)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
// 	w.WriteHeader(http.StatusCreated)
// }

func (app *App) HandleGetReservationByID(c *gin.Context) {
	id := c.Param("id")

	reservation, err := app.Service.Reservation.Repository.FindById(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		}
		return
	}

	// Return the reservation in JSON format
	c.JSON(http.StatusOK, reservation)
}

func (app *App) HandlePostReservation(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "Invalid request method",
		})
		return
	}

	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		// If JSON is invalid, return a 400 response
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON format",
		})
		return
	}

	slog.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"incoming request",
		slog.String("method", http.MethodPost),
		slog.String("path", c.FullPath()),
		slog.String("reservation", reservation.String()),
	)

	err := app.Service.Reservation.Repository.AddReservation(reservation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"reservation": reservation,
	})
}
