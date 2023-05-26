package controller

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const postMainURL = "http://localhost:8080/"

func BenchmarkPostMain(b *testing.B) {
	client := resty.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := uuid.New().String()
		client.R().SetBody("http://" + id).Post(postMainURL)
	}
}
