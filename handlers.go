package main

import (
	"fmt"
	"log"
	"net/http"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING readinessHandler...")

	respStruct := readinessResp{
		Status: "ok",
	}

	respondWithJSON(w, 200, respStruct)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING errorHandler...")

	respStruct := errorResp{
		Error: "Internal Server Error",
	}

	respondWithJSON(w, 500, respStruct)
}
