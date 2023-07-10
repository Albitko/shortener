package http

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const (
	mainURL         = "http://localhost:8080/"
	shortenURL      = "http://localhost:8080/api/shorten"
	shortenBatchURL = "http://localhost:8080/api/shorten/batch"
	userURL         = "http://localhost:8080/api/user/urls"
)

func BenchmarkHandlers(b *testing.B) {
	client := resty.New()
	b.ResetTimer()
	b.Run("POST /", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			id := uuid.New().String()
			b.StartTimer()
			client.R().SetBody("http://" + id).Post(mainURL)
		}
	})
	b.Run("POST /api/shorten", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			id := uuid.New().String()
			b.StartTimer()
			client.R().SetBody(`{"url":"https://` + id + `"}`).Post(shortenURL)
		}
	})
	b.Run("POST /api/shorten/batch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			correlationID := uuid.New().String()
			firstSite := uuid.New().String()
			secondSite := uuid.New().String()
			b.StartTimer()
			client.R().SetBody(
				`[{"correlation_id": "` + correlationID + `", "original_url": "https://` + firstSite +
					`"},{"correlation_id": "` + correlationID + `", "original_url": "https://` + secondSite + `"}]`,
			).Post(shortenBatchURL)
		}
	})
	b.Run("GET /:id", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			id := uuid.New().String()
			b.StartTimer()
			client.R().Get(mainURL + id)
		}
	})
	b.Run("GET /api/user/urls", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			id := uuid.New().String()
			client.R().SetBody("http://" + id).Post(mainURL)
			b.StartTimer()
			client.R().Get(userURL)
		}
	})
}
