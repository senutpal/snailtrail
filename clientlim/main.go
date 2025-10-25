package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}



func endPointHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	message := Message{
		Status: "success",
		Body:   "Request processed",
	}

	err := json.NewEncoder(writer).Encode(&message)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return
	}
}

func main() {
	http.Handle("/ping", perClientRateLimiter(http.HandlerFunc(endPointHandler)))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error listening on port 8080:", err)
	}
}
