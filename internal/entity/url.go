package entity

type (
	// URLID type for shorten URL.
	URLID string
	// OriginalURL type for common URL.
	OriginalURL string
)

// Config type that represent application settings.
type Config struct {
	ServerAddress        string `env:"SERVER_ADDRESS"`
	BaseURL              string `env:"BASE_URL"`
	FileStoragePath      string `env:"FILE_STORAGE_PATH"`
	CookiesStorageSecret string `env:"COOKIES_STORAGE_SECRET"`
	DatabaseDSN          string `env:"DATABASE_DSN"`
}

// ModelURLForDelete type that represents JSON struct for deleting via
// DELETE /api/user/urls.
type ModelURLForDelete struct {
	UserID   string
	ShortURL string
}

// UserURL describes a pair of urls (shorten and full) for the user.
type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ModelURLBatchRequest type that represents JSON struct for multiple URL shorten.
type ModelURLBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ModelURLBatchResponse type that represents JSON response with multiple shorten URLs..
type ModelURLBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
