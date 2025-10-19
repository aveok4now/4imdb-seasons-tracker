package repository

import (
	"4imdb-seasons-tracker/internal/domain"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type JSONRepository struct {
	mu       sync.RWMutex
	data     map[string]*domain.Series
	filePath string
	logger   *log.Logger
}

func NewJSONRepository(filePath string, logger *log.Logger) *JSONRepository {
	return &JSONRepository{
		data:     make(map[string]*domain.Series),
		filePath: filePath,
		logger:   logger,
	}
}

func (r *JSONRepository) Load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if len(data) == 0 {
		return nil
	}

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read file: %w", err)
	}

	var series []domain.Series
	if err := json.Unmarshal(data, &series); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	r.data = make(map[string]*domain.Series)
	for i := range series {
		r.data[series[i].Key()] = &series[i]
	}

	return nil
}

func (r *JSONRepository) Save() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	series := make([]domain.Series, 0, len(r.data))
	for _, s := range r.data {
		series = append(series, *s)
	}

	data, err := json.MarshalIndent(series, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (r *JSONRepository) GetAll() ([]*domain.Series, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Series, 0, len(r.data))
	for _, s := range r.data {
		result = append(result, s)
	}
	return result, nil
}

func (r *JSONRepository) Get(imdbID string, season int) (*domain.Series, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := domain.FormatKey(imdbID, season)
	series, exists := r.data[key]
	if !exists {
		return nil, fmt.Errorf("series not found")
	}
	return series, nil
}

func (r *JSONRepository) Add(series *domain.Series) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := series.Key()
	if _, exists := r.data[key]; exists {
		return fmt.Errorf("series already exists")
	}

	r.data[key] = series
	return nil
}

func (r *JSONRepository) Update(series *domain.Series) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := series.Key()
	if _, exists := r.data[key]; !exists {
		return fmt.Errorf("series not found")
	}

	r.data[key] = series
	return nil
}

func (r *JSONRepository) Delete(imdbID string, season int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := domain.FormatKey(imdbID, season)
	if _, exists := r.data[key]; !exists {
		return fmt.Errorf("series not found")
	}

	delete(r.data, key)
	return nil
}
