package urlgenerator

type URLGenerator interface {
	GenerateURL(longURL string) string
}