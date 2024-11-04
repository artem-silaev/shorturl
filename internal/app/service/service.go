package service

import (
	"errors"
	e "github.com/artem-silaev/shorturl/internal/app/errors"
	"github.com/artem-silaev/shorturl/internal/app/repository"
	"github.com/artem-silaev/shorturl/internal/app/storage"
	"github.com/artem-silaev/shorturl/internal/app/urlgenerator"
	"github.com/google/uuid"
)

type ShortenerService struct {
	URLShortener
	repo         repository.URLRepository
	urlGenerator urlgenerator.URLGenerator
	storage      *storage.Storage
}

func NewShortenerService(filePath string) *ShortenerService {
	return &ShortenerService{
		repo:         repository.NewInMemoryURLRepository(),
		urlGenerator: urlgenerator.NewBase64EncodeGenerator(),
		storage:      storage.NewStorage(filePath),
	}
}

func (s *ShortenerService) ShortenURL(longURL string) (string, error) {
	shortURL := s.urlGenerator.GenerateURL(longURL)
	if err := s.repo.AddURL(shortURL, longURL); err != nil {
		return "", e.ErrInternal
	}

	s.storage.SaveURLs(storage.URL{ShortURL: shortURL, OriginalURL: longURL, UUID: uuid.NewString()})

	return shortURL, nil
}

func (s *ShortenerService) GetOriginalURL(shortURL string) (string, error) {
	longURL, err := s.repo.GetURL(shortURL)

	if errors.Is(err, e.ErrNotFound) {
		return "", e.ErrNotFound
	}

	if err != nil {
		return "", e.ErrInternal
	}

	return longURL, nil
}
