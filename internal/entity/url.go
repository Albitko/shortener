package entity

// type for short and original url.
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
	EnableHTTPS          bool   `env:"ENABLE_HTTPS"`
	Config               string `env:"CONFIG"`
	TrustedSubnet        string `env:"TRUSTED_SUBNET"`
}

// JSONConfig type that define app configurations
type JSONConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDsn     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	TrustedSubnet   string `json:"trusted_subnet"`
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

// URLStats describes count URLs and users in service
type URLStats struct {
	URLsCount  int `json:"urls"`
	UsersCount int `json:"users"`
}
