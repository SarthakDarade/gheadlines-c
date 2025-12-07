package handlers

import (
	"encoding/json"
	"fmt"
	"gheadlines/db"
	"gheadlines/models"
	"net/http"
	"strconv"
	"time"
)

// RealtimeHandler provides Server-Sent Events for real-time updates
func RealtimeHandler(dbClient *db.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Create a channel for updates
		updates := make(chan models.Article, 10)

		// Get access token from cookie
		var accessToken string
		if cookie, err := r.Cookie("access_token"); err == nil {
			accessToken = cookie.Value
		}

		// Start polling for new articles (every 10 seconds)
		go func() {
			lastCheck := time.Now()
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					// Fetch articles published after last check
					articles, err := dbClient.GetArticles(r.Context(), 10, 0, nil, accessToken)
					if err == nil {
						for i := range articles {
							if articles[i].CreatedAt.After(lastCheck) {
								updates <- articles[i]
							}
						}
					}
					lastCheck = time.Now()
				case <-r.Context().Done():
					return
				}
			}
		}()

		// Send updates to client
		for {
			select {
			case article := <-updates:
				data, _ := json.Marshal(map[string]interface{}{
					"type":    "new_article",
					"article": article,
				})
				fmt.Fprintf(w, "data: %s\n\n", data)
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			case <-r.Context().Done():
				return
			}
		}
	}
}

// LatestArticlesHandler returns latest articles as JSON (for polling fallback)
func LatestArticlesHandler(dbClient *db.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lastTimestamp := r.URL.Query().Get("last")
		var lastTime time.Time
		if lastTimestamp != "" {
			if ts, err := strconv.ParseInt(lastTimestamp, 10, 64); err == nil {
				lastTime = time.Unix(ts/1000, 0)
			}
		}

		// Get access token from cookie
		var accessToken string
		if cookie, err := r.Cookie("access_token"); err == nil {
			accessToken = cookie.Value
		}

		articles, err := dbClient.GetArticles(r.Context(), 10, 0, nil, accessToken)
		if err != nil {
			http.Error(w, "Failed to fetch articles", http.StatusInternalServerError)
			return
		}

		// Filter articles published after lastTime
		filtered := []models.Article{}
		for i := range articles {
			if articles[i].CreatedAt.After(lastTime) {
				filtered = append(filtered, articles[i])
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"articles": filtered,
		})
	}
}
