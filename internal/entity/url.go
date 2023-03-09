package entity

type URLID string
type OriginalURL string

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

type UserURL struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}
