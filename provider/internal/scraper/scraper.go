package scraper

import (
	"net/url"
	"time"

	"github.com/gocolly/colly/v2"
	"go.uber.org/atomic"

	"github.com/metatube-community/metatube-sdk-go/provider"
)

var _ provider.Provider = (*Scraper)(nil)

// Scraper implements basic Provider interface.
type Scraper struct {
	name     string
	priority *atomic.Int64
	baseURL  *url.URL
	c        *colly.Collector
}

// NewScraper returns Provider implemented *Scraper.
func NewScraper(name, baseURL string, priority int, opts ...Option) *Scraper {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}
	s := &Scraper{
		name:     name,
		priority: atomic.NewInt64(int64(priority)),
		baseURL:  u,
		c:        colly.NewCollector(),
	}
	for _, opt := range opts {
		// Apply options.
		if err := opt(s); err != nil {
			panic(err)
		}
	}
	return s
}

// NewDefaultScraper returns a *Scraper with default options enabled.
func NewDefaultScraper(name, baseURL string, priority int, opts ...Option) *Scraper {
	return NewScraper(name, baseURL, priority, append([]Option{
		WithAllowURLRevisit(),
		WithIgnoreRobotsTxt(),
		WithRandomUserAgent(),
	}, opts...)...)
}

func (s *Scraper) Name() string { return s.name }

func (s *Scraper) URL() *url.URL { return s.baseURL }

func (s *Scraper) Priority() int64 { return s.priority.Load() }

func (s *Scraper) SetPriority(v int64) { s.priority.Store(v) }

func (s *Scraper) NormalizeMovieID(id string) string { return id /* AS IS */ }

func (s *Scraper) ParseMovieIDFromURL(string) (string, error) { panic("unimplemented") }

func (s *Scraper) NormalizeActorID(id string) string { return id /* AS IS */ }

func (s *Scraper) ParseActorIDFromURL(string) (string, error) { panic("unimplemented") }

// ClonedCollector returns cloned internal collector.
func (s *Scraper) ClonedCollector() *colly.Collector { return s.c.Clone() }

// SetRequestTimeout sets timeout for HTTP requests.
func (s *Scraper) SetRequestTimeout(timeout time.Duration) { s.c.SetRequestTimeout(timeout) }
