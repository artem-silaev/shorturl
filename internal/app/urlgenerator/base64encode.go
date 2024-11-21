package urlgenerator

import (
	"encoding/base64"
)

type Base64EncodeGenerator struct {
	URLGenerator
}

func NewBase64EncodeGenerator() *Base64EncodeGenerator {
	return &Base64EncodeGenerator{
	}
}

func (g *Base64EncodeGenerator) GenerateURL(longURL string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(longURL))
}

