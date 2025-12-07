package handlers

import (
	"encoding/json"
	"gheadlines/db"
	"gheadlines/models"
	"log"
	"net/http"
)

// NewsletterSubscribeHandler handles newsletter subscriptions
func NewsletterSubscribeHandler(dbClient *db.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Email string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		err := dbClient.CreateNewsletterSubscriber(r.Context(), req.Email)
		if err != nil {
			http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

// ContactFormHandler handles contact form submissions
func ContactFormHandler(dbClient *db.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Name    string `json:"name"`
			Email   string `json:"email"`
			Phone   string `json:"phone"`
			Message string `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.Email == "" || req.Message == "" {
			http.Error(w, "Name, Email and Message are required", http.StatusBadRequest)
			return
		}

		msg := &models.ContactMessage{
			Name:    req.Name,
			Email:   req.Email,
			Phone:   req.Phone,
			Message: req.Message,
		}

		err := dbClient.CreateContactMessage(r.Context(), msg)
		if err != nil {
			log.Printf("Contact Form DB Error: %v", err)
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}
