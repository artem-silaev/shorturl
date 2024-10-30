package service

import (
	"errors"
	e "github.com/artem-silaev/shorturl/internal/app/errors"
	"github.com/artem-silaev/shorturl/internal/app/repository"
	"github.com/artem-silaev/shorturl/internal/app/urlgenerator"
)

type ShortenerService struct {
	URLShortener
	repo         repository.URLRepository
	urlGenerator urlgenerator.URLGenerator
}

func NewShortenerService() *ShortenerService {
	return &ShortenerService{
		repo:         repository.NewInMemoryURLRepository(),
		urlGenerator: urlgenerator.NewBase64EncodeGenerator(),
	}
}

func (s *ShortenerService) ShortenURL(longURL string) (string, error) {
	var shortURL string

	shortURL = s.urlGenerator.GenerateURL(longURL)

	if err := s.repo.AddURL(shortURL, longURL); err != nil {
		return "", e.ErrInternal
	}

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
