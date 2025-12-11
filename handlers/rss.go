package handlers

import (
	"encoding/xml"
	"gheadlines/db"
	"net/http"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	Item          []Item `xml:"item"`
}

type Item struct {
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	Description string     `xml:"description"`
	PubDate     string     `xml:"pubDate"`
	GUID        string     `xml:"guid"`
	Enclosure   *Enclosure `xml:"enclosure,omitempty"`
}

type Enclosure struct {
	URL  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

// RSSHandler generates and serves RSS feed
func RSSHandler(dbClient *db.Client, siteURL string, siteName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch latest 100 articles
		articles, err := dbClient.GetArticles(r.Context(), 100, 0, nil, "")
		if err != nil {
			http.Error(w, "Failed to fetch articles", http.StatusInternalServerError)
			return
		}

		rss := RSS{
			Version: "2.0",
			Channel: Channel{
				Title:         siteName,
				Link:          siteURL,
				Description:   "Latest news and updates from " + siteName,
				Language:      "en-us",
				LastBuildDate: time.Now().Format(time.RFC1123),
				Item:          []Item{},
			},
		}

		for _, a := range articles {
			item := Item{
				Title:       a.Title,
				Link:        siteURL + "/article/" + a.ID,
				Description: a.Excerpt,
				PubDate:     a.CreatedAt.Format(time.RFC1123),
				GUID:        siteURL + "/article/" + a.ID,
			}

			// Add image enclosure if available
			if a.ImageURL != "" {
				item.Enclosure = &Enclosure{
					URL:  a.ImageURL,
					Type: "image/jpeg", // Assuming JPEG for simplicity, or strict check
				}
			}

			rss.Channel.Item = append(rss.Channel.Item, item)
		}

		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		encoder := xml.NewEncoder(w)
		encoder.Indent("", "  ")
		if err := encoder.Encode(rss); err != nil {
			http.Error(w, "Failed to encode RSS", http.StatusInternalServerError)
		}
	}
}
