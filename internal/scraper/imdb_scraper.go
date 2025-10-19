package scraper

import (
	"4imdb-seasons-tracker/internal/config"
	"4imdb-seasons-tracker/internal/domain"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type IMDBScraper struct {
	config config.ScraperConfig
	client *http.Client
	logger *log.Logger
}

func NewIMDBScraper(config config.ScraperConfig, logger *log.Logger) *IMDBScraper {
	return &IMDBScraper{
		config: config,
		client: &http.Client{
			Timeout: config.RequestTimeout,
		},
		logger: logger,
	}
}

func (s *IMDBScraper) FetchEpisodeInfo(imdbID string, season int) (*domain.EpisodeInfo, error) {
	url := fmt.Sprintf("%s/title/%s/episodes/?season=%d", s.config.BaseURL, imdbID, season)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", s.config.UserAgent)
	req.Header.Set("Accept-Language", s.config.AcceptLanguage)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse HTML: %w", err)
	}

	return s.parseEpisodeInfo(doc), nil
}

func (s *IMDBScraper) parseEpisodeInfo(doc *goquery.Document) *domain.EpisodeInfo {
	info := &domain.EpisodeInfo{}

	doc.Find("article.episode-item-wrapper").First().Each(func(i int, sel *goquery.Selection) {
		titleText := sel.Find("h4[data-testid='slate-list-card-title'] .ipc-title__text").Text()
		info.Title = strings.TrimSpace(titleText)

		dateText := sel.Find("span.knzESm").Text()
		info.ReleaseDate = strings.TrimSpace(dateText)

		plotButton := sel.Find("a.ipc-text-button:contains('Add a plot')").Length()
		plotDiv := sel.Find(".ipc-html-content-inner-div").Text()

		info.HasPlot = plotButton == 0 && len(strings.TrimSpace(plotDiv)) > 0
	})

	return info
}
