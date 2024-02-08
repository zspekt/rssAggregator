package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	fmt.Print("\n\n\n")
	log.Printf("Responding with error code -> %v and the following message:\n\t%v", code, msg)
	respondWithJSON(w, code, msg)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	fmt.Print("\n\n\n")
	log.Printf("Responding with code -> %v and provided payload...\n", code)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return
	}

	// fmt.Println(string(jsonPayload))
	w.Write(jsonPayload)
}

func decodeJson[T any](r io.Reader, st *T) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(st)
	if err != nil {
		return err
	}
	return nil
}
