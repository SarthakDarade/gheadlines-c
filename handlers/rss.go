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
	XmlnsDc      string   `xml:"xmlns:dc,attr"`
	XmlnsMedia   string   `xml:"xmlns:media,attr"`   // Keeping for Google Discover support
	XmlnsContent string   `xml:"xmlns:content,attr"` // Keeping for full content support
	Channel      Channel  `xml:"channel"`
}

type Channel struct {
	AtomLink      *AtomLink     `xml:"atom:link,omitempty"`
	Title         string        `xml:"title"`
	Link          string        `xml:"link"`
	Description   string        `xml:"description"`
	Language      string        `xml:"language"`
	Copyright     string        `xml:"copyright"`
	Docs          string        `xml:"docs"`
	Image         *ChannelImage `xml:"image,omitempty"`
	LastBuildDate string        `xml:"lastBuildDate"`
	Item          []Item        `xml:"item"`
}

type ChannelImage struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	URL   string `xml:"url"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Item struct {
	Title          CDATA         `xml:"title"`
	Description    CDATA         `xml:"description"`
	Link           CDATA         `xml:"link"`
	GUID           CDATA         `xml:"guid"`
	PubDate        string        `xml:"pubDate"`
	Creator        string        `xml:"dc:creator,omitempty"`
	Enclosure      *Enclosure    `xml:"enclosure,omitempty"`
	MediaContent   *MediaContent `xml:"media:content,omitempty"`
	ContentEncoded *CDATA        `xml:"content:encoded,omitempty"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr,omitempty"` // Added length
	Type   string `xml:"type,attr"`
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
			XmlnsDc:      "http://purl.org/dc/elements/1.1/",
			XmlnsMedia:   "http://search.yahoo.com/mrss/",
			XmlnsContent: "http://purl.org/rss/1.0/modules/content/",
			Channel: Channel{
				Title:         siteName,
				Link:          siteURL,
				Description:   "Latest news and updates from " + siteName,
				Language:      "en-us", // TOI uses en-gb, we can stick to en-us or make dynamic
				Copyright:     fmt.Sprintf("Copyright (C) %d %s", time.Now().Year(), siteName),
				Docs:          siteURL + "/rss",
				LastBuildDate: time.Now().Format(time.RFC1123Z),
				AtomLink: &AtomLink{
					Href: siteURL + "/rss",
					Rel:  "self",
					Type: "application/rss+xml",
				},
				// Add a default logo/image for the channel
				Image: &ChannelImage{
					Title: siteName,
					Link:  siteURL,
					URL:   siteURL + "/static/gheadlineicon.png", // Ensure this path is correct
				},
				Item: []Item{},
			},
		}

		for _, a := range articles {
			link := fmt.Sprintf("%s/article/%s", siteURL, a.Slug)
			if a.Slug == "" {
				link = fmt.Sprintf("%s/article/%s", siteURL, a.ID)
			}

			// Sanitize or fallback for author
			author := a.Author
			if author == "" {
				author = siteName
			}

			item := Item{
				Title:       CDATA{Value: a.Title},
				Description: CDATA{Value: a.Excerpt},
				Link:        CDATA{Value: link},
				GUID:        CDATA{Value: link},
				PubDate:     a.CreatedAt.Format(time.RFC1123Z),
				Creator:     author,
			}

			if a.Content != "" {
				item.ContentEncoded = &CDATA{Value: a.Content}
			}

			// Add image enclosure and media:content
			if a.ImageURL != "" {
				// Standard Enclosure
				item.Enclosure = &Enclosure{
					URL:    a.ImageURL,
					Type:   "image/jpeg",
					Length: "0", // Supabase doesn't easily give us size, using 0 or omitted is better than wrong. TOI has it.
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
