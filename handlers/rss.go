package handlers

import (
	"encoding/xml"
	"fmt"
	"gheadlines/db"
	"net/http"
	"time"
)

type RSS struct {
	XMLName   xml.Name `xml:"rss"`
	Version   string   `xml:"version,attr"`
	XmlnsAtom string   `xml:"xmlns:atom,attr"`
	Channel   Channel  `xml:"channel"`
}

type Channel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	AtomLink      *AtomLink `xml:"atom:link,omitempty"`
	Item          []Item    `xml:"item"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
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
			Version:   "2.0",
			XmlnsAtom: "http://www.w3.org/2005/Atom",
			Channel: Channel{
				Title:         siteName,
				Link:          siteURL,
				Description:   "Latest news and updates from " + siteName,
				Language:      "en-us",
				LastBuildDate: time.Now().Format(time.RFC1123),
				AtomLink: &AtomLink{
					Href: siteURL + "/rss",
					Rel:  "self",
					Type: "application/rss+xml",
				},
				Item: []Item{},
			},
		}

		for _, a := range articles {
			link := fmt.Sprintf("%s/article/%s", siteURL, a.Slug)
			if a.Slug == "" {
				link = fmt.Sprintf("%s/article/%s", siteURL, a.ID)
			}

			item := Item{
				Title:       a.Title,
				Link:        link,
				Description: a.Excerpt,
				PubDate:     a.CreatedAt.Format(time.RFC1123),
				GUID:        link,
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

		w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "  ")
		if err := encoder.Encode(rss); err != nil {
			http.Error(w, "Failed to encode RSS", http.StatusInternalServerError)
		}
	}
}
