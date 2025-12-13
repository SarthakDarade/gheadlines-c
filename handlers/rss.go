package handlers

import (
	"encoding/xml"
	"fmt"
	"gheadlines/db"
	"net/http"
	"time"
)

type RSS struct {
	XMLName      xml.Name `xml:"rss"`
	Version      string   `xml:"version,attr"`
	XmlnsAtom    string   `xml:"xmlns:atom,attr"`
	XmlnsMedia   string   `xml:"xmlns:media,attr"`
	XmlnsDc      string   `xml:"xmlns:dc,attr"`
	XmlnsContent string   `xml:"xmlns:content,attr"`
	Channel      Channel  `xml:"channel"`
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
	Title          string        `xml:"title"`
	Link           string        `xml:"link"`
	Description    string        `xml:"description"`
	PubDate        string        `xml:"pubDate"`
	GUID           string        `xml:"guid"`
	Enclosure      *Enclosure    `xml:"enclosure,omitempty"`
	Creator        string        `xml:"dc:creator,omitempty"`
	Category       string        `xml:"category,omitempty"`
	ContentEncoded *CDATA        `xml:"content:encoded,omitempty"`
	MediaContent   *MediaContent `xml:"media:content,omitempty"`
}

type Enclosure struct {
	URL  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

type MediaContent struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Medium string `xml:"medium,attr"`
}

type CDATA struct {
	Value string `xml:",cdata"`
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
			Version:      "2.0",
			XmlnsAtom:    "http://www.w3.org/2005/Atom",
			XmlnsMedia:   "http://search.yahoo.com/mrss/",
			XmlnsDc:      "http://purl.org/dc/elements/1.1/",
			XmlnsContent: "http://purl.org/rss/1.0/modules/content/",
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
				Creator:     a.Author,
				Category:    a.Category,
			}

			if a.Content != "" {
				item.ContentEncoded = &CDATA{Value: a.Content}
			}

			// Add image enclosure and media:content for better compatibility
			if a.ImageURL != "" {
				// Standard Enclosure (Podcasts/Legacy)
				item.Enclosure = &Enclosure{
					URL:  a.ImageURL,
					Type: "image/jpeg",
				}
				// Media RSS (Google News/Discovery)
				item.MediaContent = &MediaContent{
					URL:    a.ImageURL,
					Type:   "image/jpeg",
					Medium: "image",
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
