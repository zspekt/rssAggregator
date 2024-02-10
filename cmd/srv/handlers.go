package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zspekt/rssAggregator/internal/jsonGenerics"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING readinessHandler...")

	respStruct := readinessResp{
		Status: "ok",
	}

	jsongenerics.RespondWithJSON(w, 200, respStruct)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n\n")
	log.Println("RUNNING errorHandler...")

	respStruct := errorResp{
		Error: "Internal Server Error",
	}

	jsongenerics.RespondWithJSON(w, 500, respStruct)
}
