# Global Headlines - World-Class SEO News Website

A production-ready, ultra-fast, SEO-optimized global news website built with Go, TailwindCSS, and Supabase.

## Features

- **Ultra-Fast Performance**: Server-side rendering with sub-200ms response times
- **Perfect SEO**: Complete JSON-LD structured data, sitemap.xml, canonical URLs, OpenGraph tags
- **Beautiful Design**: NYTimes/BBC-inspired clean editorial design with TailwindCSS
- **Scalable Backend**: Go handles thousands of concurrent users
- **Real-time Updates**: Supabase Realtime integration with minimal JavaScript
- **Mobile-First**: Fully responsive design optimized for all devices

## Tech Stack

- **Backend**: Go (Golang) with server-side rendering
- **Frontend**: HTML templates with TailwindCSS
- **Database**: Supabase (PostgreSQL)
- **CDN**: Cloudflare for images
- **Automation**: Ready for n8n integration

## Project Structure

```
gheadlines/
├── main.go                 # Go HTTP server with SSR routes
├── config/                 # Configuration loader
├── models/                 # Data models (Article, Author, Category)
├── db/                     # Supabase connection & helpers
├── services/               # Business logic (SEO, news processing)
├── handlers/               # HTTP handlers
├── templates/              # HTML templates
│   ├── base.html
│   ├── index.html
│   ├── article.html
│   └── partials/
└── static/                 # Static files (robots.txt)
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Supabase account (optional for MVP - uses dummy data if not configured)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd gheadlines
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure environment variables (optional for MVP):
```bash
cp env.example .env
# Edit .env with your Supabase credentials
```

4. Run the server:
```bash
go run main.go
# Or build and run:
go build -o gheadlines.exe .
./gheadlines.exe
```

5. Open your browser:
```
http://localhost:8080
```

## Configuration

Create a `.env` file with the following variables:

```env
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-anon-key-here
PORT=8080
HOST=localhost
SITE_URL=http://localhost:8080
SITE_NAME=Global Headlines
```

If Supabase is not configured, the application will use dummy data for MVP testing.

## Routes

- `/` - Homepage with article grid
- `/article/:slug` - Individual article page
- `/category/:name` - Category filtered articles
- `/sitemap.xml` - Dynamic sitemap
- `/robots.txt` - Robots.txt file
- `/api/realtime` - Server-Sent Events for real-time updates
- `/health` - Health check endpoint

## SEO Features

- Complete JSON-LD structured data (Article, NewsArticle schemas)
- Dynamic sitemap.xml with lastmod dates
- Canonical URLs on every page
- OpenGraph and Twitter Card meta tags
- Semantic HTML5 structure
- Fast Core Web Vitals optimization

## Performance Optimizations

- Gzip compression for all text content
- HTTP/2 support
- Cache-Control headers for static assets
- Lazy loading images
- Preconnect hints for CDN
- Minimal JavaScript (only for real-time updates)

## Design

The design is inspired by NYTimes and BBC, featuring:
- Clean editorial typography (serif headlines, sans-serif body)
- Professional color palette
- Responsive grid layouts
- Sticky navigation
- Smooth hover effects and transitions
- Mobile-first responsive design

## Development

### Building

```bash
go build -o gheadlines.exe .
```

### Running Tests

```bash
go test ./...
```

## Production Deployment

1. Set environment variables in your deployment platform
2. Build the binary: `go build -o gheadlines .`
3. Run the server with proper process management (systemd, PM2, etc.)
4. Configure reverse proxy (nginx, Caddy) for HTTPS
5. Set up Cloudflare CDN for images
6. Configure Supabase database with proper schema

## Future Enhancements

- n8n automation for news fetching and rewriting
- AI-based headline optimization
- Multi-language support (i18n)
- AMP pages for mobile speed
- Advanced analytics integration
- User authentication and comments

## License

MIT License


