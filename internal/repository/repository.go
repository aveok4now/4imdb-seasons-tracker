package repository

import "4imdb-seasons-tracker/internal/domain"

type Repository interface {
	Load() error
	Save() error
	GetAll() ([]*domain.Series, error)
	Get(imdbID string, season int) (*domain.Series, error)
	Add(*domain.Series) error
	Update(*domain.Series) error
	Delete(imdbID string, season int) error
}
