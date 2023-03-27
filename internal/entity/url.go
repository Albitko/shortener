package entity

type (
	URLID       string
	OriginalURL string
)

type Config struct {
	ServerAddress        string `env:"SERVER_ADDRESS"`
	BaseURL              string `env:"BASE_URL"`
	FileStoragePath      string `env:"FILE_STORAGE_PATH"`
	CookiesStorageSecret string `env:"COOKIES_STORAGE_SECRET"`
	DatabaseDSN          string `env:"DATABASE_DSN"`
}

type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ModelURLBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ModelURLBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
