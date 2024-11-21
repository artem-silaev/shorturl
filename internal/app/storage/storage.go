package storage

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/artem-silaev/shorturl/internal/app/repository"
)

type Storage struct {
	filePath string
}

type URL struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewStorage(filePath string) *Storage {
	return &Storage{
		filePath: filePath,
	}
}

func (s *Storage) LoadURLs(repo repository.URLRepository) error {
	file, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	lines := string(file)
	for _, line := range strings.Split(lines, "\n") {
		if line != "" {
			var url URL
			if err := json.Unmarshal([]byte(line), &url); err != nil {
				return err
			}
			repo.AddURL(url.ShortURL, url.OriginalURL)
		}
	}

	return nil
}

func (s *Storage) SaveURLs(url URL) error {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)

	if err != nil {
		return err
	}
	defer file.Close()

	line, err := json.Marshal(url)
	if err != nil {
		return err
	}
	_, err = file.Write(append(line, '\n'))
	if err != nil {
		return err
	}

	return nil
}
