package handlers

import (
	"encoding/json"
	"gheadlines/db"
	"net/http"
)

// BreakingNewsAPIHandler returns the latest breaking news as JSON
func BreakingNewsAPIHandler(dbClient *db.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		news, err := dbClient.GetBreakingNews(r.Context(), 5)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(news)
	}
}
