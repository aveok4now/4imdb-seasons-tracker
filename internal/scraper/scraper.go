package scraper

import "4imdb-seasons-tracker/internal/domain"

type Scraper interface {
	FetchEpisodeInfo(imdbID string, season int) (*domain.EpisodeInfo, error)
}
