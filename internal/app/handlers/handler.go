package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"wedding_website/internal/app/models"
	"wedding_website/internal/app/repository"
	"wedding_website/internal/app/telegram"
)

type RSVPHandler struct {
	repo     *repository.RSVPRepository
	telegram *telegram.TelegramService
}

func NewRSVPHandler(db *sql.DB, telegramService *telegram.TelegramService) *RSVPHandler {
	return &RSVPHandler{
		repo:     repository.NewRSVPRepository(db),
		telegram: telegramService,
	}
}

func (h *RSVPHandler) HandleRSVP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RSVPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(req.Name)
	attendance := req.Attendance
	companion := strings.TrimSpace(req.Companion)

	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if attendance == "" {
		http.Error(w, "Attendance is required", http.StatusBadRequest)
		return
	}

	rsvp := &models.RSVP{
		Name:       name,
		Attendance: attendance == "yes",
		Companion:  companion,
	}

	if err := h.repo.CreateRSVP(rsvp); err != nil {
		log.Printf("Error saving RSVP: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	go func() {
		if err := h.telegram.SendRSVPNotification(rsvp); err != nil {
			log.Printf("Failed to send Telegram notification: %v", err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Спасибо за ответ!",
		"status":  "success",
	})
}
