package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/artem-silaev/shorturl/internal/app/service"
)

type MockShortenerService struct {
	service.URLShortener
	mock.Mock
}

func NewMockShortenerService() *MockShortenerService {
	return &MockShortenerService{}
}

func (m *MockShortenerService) ShortenURL(longURL string) (string, error) {
	args := m.Called(longURL)
	return args.String(0), args.Error(1)
}

func (m *MockShortenerService) GetOriginalURL(shortURL string) (string, error) {
	args := m.Called(shortURL)
	return args.String(0), args.Error(1)
}
