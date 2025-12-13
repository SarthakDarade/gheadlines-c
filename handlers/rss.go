package handlers

import (
	"bytes"
	"gheadlines/db"
	"net/http"
	"strings"
	"text/template"
	"time"
)

// RSSHandler generates and serves RSS feed using templates for precise control
func RSSHandler(dbClient *db.Client, siteURL string, siteName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch latest 100 articles
		articles, err := dbClient.GetArticles(r.Context(), 100, 0, nil, "")
		if err != nil {
			http.Error(w, "Failed to fetch articles", http.StatusInternalServerError)
			return
		}

		// Prepare data for template
		type RSSItem struct {
			Title       string
			Link        string
			Description string
			Author      string
			PubDate     string
			GUID        string
			ImageURL    string
			Content     string
		}

		type RSSData struct {
			SiteName      string
			SiteURL       string
			Description   string
			Language      string
			Copyright     string
			Docs          string
			LastBuildDate string
			AtomLink      string
			Items         []RSSItem
		}

		// Use RFC1123Z (e.g. Mon, 02 Jan 2006 15:04:05 -0700) which is standard for RSS 2.0
		now := time.Now().Format(time.RFC1123Z)

		data := RSSData{
			SiteName:      siteName,
			SiteURL:       siteURL,
			Description:   "Latest news and updates from " + siteName,
			Language:      "en-us",
			Copyright:     "Copyright (C) " + time.Now().Format("2006") + " " + siteName,
			Docs:          siteURL + "/rss",
			LastBuildDate: now,
			AtomLink:      siteURL + "/rss",
			Items:         make([]RSSItem, 0, len(articles)),
		}

		for _, a := range articles {
			link := siteURL + "/article/" + a.Slug
			if a.Slug == "" {
				link = siteURL + "/article/" + a.ID
			}

			author := a.Author
			if author == "" {
				author = siteName
			}

			// Clean potentially problematic characters
			title := strings.TrimSpace(a.Title)
			desc := strings.TrimSpace(a.Excerpt)
			// Remove newlines from description to keep it clean in snippet view, though CDATA handles it
			desc = strings.ReplaceAll(desc, "\n", " ")
			desc = strings.ReplaceAll(desc, "\r", "")

			item := RSSItem{
				Title:       title,
				Link:        link,
				Description: desc,
				Author:      author,
				PubDate:     a.CreatedAt.Format(time.RFC1123Z),
				GUID:        link,
				ImageURL:    a.ImageURL,
				Content:     a.Content,
			}
			data.Items = append(data.Items, item)
		}

		// Defined template with NO whitespace inside CDATA for URLs
		// This is critical for feed validation. " http://..." is an invalid URL.
		const rssTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:media="http://search.yahoo.com/mrss/">
<channel>
<atom:link href="{{.AtomLink}}" rel="self" type="application/rss+xml" />
<title>{{.SiteName}}</title>
<link>{{.SiteURL}}</link>
<description>{{.Description}}</description>
<language>{{.Language}}</language>
<copyright>{{.Copyright}}</copyright>
<docs>{{.Docs}}</docs>
<image>
<title>{{.SiteName}}</title>
<link>{{.SiteURL}}</link>
<url>{{.SiteURL}}/static/gheadlineicon.png</url>
</image>
<lastBuildDate>{{.LastBuildDate}}</lastBuildDate>
{{range .Items}}<item>
<title><![CDATA[{{.Title}}]]></title>
<description><![CDATA[{{.Description}}]]></description>
<link><![CDATA[{{.Link}}]]></link>
<guid><![CDATA[{{.GUID}}]]></guid>
<pubDate>{{.PubDate}}</pubDate>
<dc:creator><![CDATA[{{.Author}}]]></dc:creator>
{{if .ImageURL}}<enclosure url="{{.ImageURL}}" type="image/jpeg" />{{end}}
{{if .ImageURL}}<media:content url="{{.ImageURL}}" type="image/jpeg" medium="image" />{{end}}
{{if .Content}}<content:encoded><![CDATA[{{.Content}}]]></content:encoded>{{end}}
</item>
{{end}}</channel>
</rss>`

		tmpl, err := template.New("rss").Parse(rssTemplate)
		if err != nil {
			http.Error(w, "Failed to parse RSS template", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		var buffer bytes.Buffer
		if err := tmpl.Execute(&buffer, data); err != nil {
			// In production, log error
			return
		}
		w.Write(buffer.Bytes())
	}
}
