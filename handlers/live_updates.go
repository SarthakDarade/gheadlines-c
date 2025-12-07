package handlers

import (
	"encoding/json"
	"gheadlines/db"
	"net/http"
)

// LiveUpdatesAPIHandler returns the latest live updates as JSON
func LiveUpdatesAPIHandler(dbClient *db.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch latest 10 updates
		updates, err := dbClient.GetLiveUpdates(r.Context(), 10, "")
		if err != nil {
			http.Error(w, "Failed to fetch updates", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updates)
	}
}
