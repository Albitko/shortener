package controller

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"
)

const serverAddr = "http://localhost:8080"

func Example() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	client := resty.New()

	// send the url that needs to be shortened and receive shorten
	client.R().SetContext(ctx).SetBody("https://google.com").Post(serverAddr + "/")
	// send the url that needs to be shortened and receive shorten in JSON
	client.R().SetContext(ctx).SetBody(`{"url":"https://google.com"}`).Post(serverAddr + "/api/shorten")
	// send  urls that needs to be shortened and receive shorten in JSON
	client.R().
		SetContext(ctx).
		SetBody(`[{"correlation_id": "q1", "original_url": "https://news.com"}, {"correlation_id": "q1", "original_url": "https://mail.com"}]`).
		Post(serverAddr + "/api/shorten/batch")
	// :id - short urls form post request to / or /api/shorten. Get original by short
	client.R().SetContext(ctx).Get(serverAddr + "/shorten_url")
	// get all urls short+original for user
	client.R().SetContext(ctx).Get(serverAddr + "/api/user/urls")
	// delete shorten urls for user
	client.R().SetContext(ctx).
		SetBody(`["short_url_1", "short_url_2", "short_url_3"]`).
		Delete(serverAddr + "/api/user/urls")
	// check DB connection
	client.R().SetContext(ctx).Get(serverAddr + "/ping")
}
