package controller

import (
	"bytes"
	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/usecase/repo"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {

	type want struct {
		code          int
		requestBody   string
		requestMethod string
		response      string
		params        string
		location      string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "test POST request",
			want: want{
				code:          http.StatusCreated,
				requestBody:   "https://google.com",
				requestMethod: http.MethodPost,
				response:      "http://localhost:8080/cv6VxVduxj",
				params:        "",
			},
		},
		{
			name: "test GET request",
			want: want{
				code:          http.StatusTemporaryRedirect,
				requestBody:   "",
				requestMethod: http.MethodGet,
				response:      "",
				params:        "cv6VxVduxj",
				location:      "https://google.com",
			},
		},
		{
			name: "test Bad request",
			want: want{
				code:          http.StatusBadRequest,
				requestBody:   "230f2jdql",
				requestMethod: http.MethodPut,
				response:      "",
				params:        "",
			},
		},
	}

	repository := repo.NewRepository()
	uc := usecase.NewURLConverter(repository)
	hndl := NewURLHandler(uc)
	http.Handle("/", hndl)

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.want.requestMethod, "/"+tt.want.params, bytes.NewBuffer([]byte(tt.want.requestBody)))

			w := httptest.NewRecorder()

			hndl.ServeHTTP(w, request)
			res := w.Result()

			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			locationURL, err := res.Location()

			if err == nil && locationURL.String() != tt.want.location {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Errorf("Error closing body %s", err)
				}
			}(res.Body)
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}
		})
	}
}
