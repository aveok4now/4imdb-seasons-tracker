package domain

import (
	"strconv"
	"time"
)

type Series struct {
	ImdbID      string    `json:"imdb_id"`
	Season      int       `json:"season"`
	LastChecked time.Time `json:"last_checked"`
	Status      Status    `json:"status"`
}

type Status string

const (
	StatusPlaceholder Status = "placeholder"
	StatusAnnounced   Status = "announced"
	StatusUnknown     Status = "unknown"
)

type EpisodeInfo struct {
	Title       string
	ReleaseDate string
	HasPlot     bool
}

func (s *Series) Key() string {
	return FormatKey(s.ImdbID, s.Season)
}

func FormatKey(imdbID string, season int) string {
	return imdbID + "-s" + strconv.Itoa(season)
}
