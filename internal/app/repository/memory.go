package repository

import (
	"sync"

	e "github.com/artem-silaev/shorturl/internal/app/errors"
)

type InMemoryURLRepository struct {
	sync.RWMutex
	URLRepository
	values map[string]string
}

func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{
		values: map[string]string{},
	}
}

func (ms *InMemoryURLRepository) AddURL(shortURL string, longURL string) error {
	ms.Lock()
	defer ms.Unlock()
	_, ok := ms.values[shortURL]
	if ok {
		return nil
	}
	ms.values[shortURL] = longURL
	return nil
}

func (ms *InMemoryURLRepository) GetURL(shortURL string) (string, error) {
	ms.Lock()
	defer ms.Unlock()

	value, ok := ms.values[shortURL]
	if !ok {
		return "", e.ErrNotFound
	}
	return value, nil
}

func (ms *InMemoryURLRepository) DelURL(shortURL string) error {
	ms.Lock()
	defer ms.Unlock()
	delete(ms.values, shortURL)
	return nil
}

func (ms *InMemoryURLRepository) GetAllURL() map[string]string {
	ms.Lock()
	defer ms.Unlock()

	return ms.values
}
