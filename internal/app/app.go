package app

import (
	"log"
	"net/http"
)

func Run() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
