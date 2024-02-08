package main

import (
	"fmt"
	"log"
	"net/http"
)

func GetApiKeyFromHeader(r *http.Request) (string, error) {
	apiKey := r.Header.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("Authorization header is missing")
	}

	log.Printf("Retrieved token from header -> %v\n", apiKey)
	return apiKey, nil
}
