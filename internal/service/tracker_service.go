package service

import (
	"4imdb-seasons-tracker/internal/domain"
	"4imdb-seasons-tracker/internal/repository"
	"4imdb-seasons-tracker/internal/scraper"
	"fmt"
	"log"
	"strings"
	"time"
)

type TrackerService struct {
	repo    repository.Repository
	scraper scraper.Scraper
	logger  *log.Logger
}

func NewTrackerService(repo repository.Repository, scraper scraper.Scraper, logger *log.Logger) *TrackerService {
	return &TrackerService{
		repo:    repo,
		scraper: scraper,
		logger:  logger,
	}
}

func (s *TrackerService) AddSeries(imdbID string, season int) (string, error) {
	info, err := s.scraper.FetchEpisodeInfo(imdbID, season)
	if err != nil {
		return "", fmt.Errorf("fetch episode info: %w", err)
	}

	if s.isAnnounced(info) {
		return fmt.Sprintf("Season %d of %s is already announced:\nTitle: %s\nRelease: %s",
			season, imdbID, info.Title, info.ReleaseDate), nil
	}

	series := &domain.Series{
		ImdbID:      imdbID,
		Season:      season,
		LastChecked: time.Now(),
		Status:      domain.StatusPlaceholder,
	}

	if err := s.repo.Add(series); err != nil {
		return "", fmt.Errorf("add to repository: %w", err)
	}

	if err := s.repo.Save(); err != nil {
		return "", fmt.Errorf("save repository: %w", err)
	}

	return fmt.Sprintf("Now tracking %s season %d", imdbID, season), nil
}

func (s *TrackerService) GetAll() ([]*domain.Series, error) {
	return s.repo.GetAll()
}

func (s *TrackerService) CheckAll() error {
	series, err := s.repo.GetAll()
	if err != nil {
		return fmt.Errorf("get all series: %w", err)
	}

	for _, ser := range series {
		if err := s.checkSeries(ser); err != nil {
			s.logger.Printf("Error checking %s season %d: %v", ser.ImdbID, ser.Season, err)
			continue
		}
		time.Sleep(2 * time.Second)
	}

	if err := s.repo.Save(); err != nil {
		return fmt.Errorf("save repository: %w", err)
	}

	return nil
}

func (s *TrackerService) checkSeries(series *domain.Series) error {
	info, err := s.scraper.FetchEpisodeInfo(series.ImdbID, series.Season)
	if err != nil {
		return err
	}

	oldStatus := series.Status
	series.LastChecked = time.Now()

	if s.isAnnounced(info) {
		series.Status = domain.StatusAnnounced
	} else {
		series.Status = domain.StatusPlaceholder
	}

	if err := s.repo.Update(series); err != nil {
		return err
	}

	if oldStatus != domain.StatusAnnounced && series.Status == domain.StatusAnnounced {
		s.logger.Printf("🎉 NEW ANNOUNCEMENT: %s Season %d", series.ImdbID, series.Season)
		s.logger.Printf("   Title: %s", info.Title)
		s.logger.Printf("   Release: %s", info.ReleaseDate)
	}

	return nil
}

func (s *TrackerService) isAnnounced(info *domain.EpisodeInfo) bool {
	if info.Title == "" {
		return false
	}

	isPlaceholder := strings.Contains(info.Title, "Episode #")
	hasRealDate := s.hasSpecificDate(info.ReleaseDate)

	return !isPlaceholder || hasRealDate
}

func (s *TrackerService) hasSpecificDate(dateStr string) bool {
	if dateStr == "" {
		return false
	}
	dateStr = strings.ToLower(dateStr)
	months := []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}
	for _, month := range months {
		if strings.Contains(dateStr, month) {
			return true
		}
	}
	return strings.Contains(dateStr, ",")
}
