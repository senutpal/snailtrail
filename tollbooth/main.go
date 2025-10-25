package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
)


type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func endPointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Message{
		Status: "success",
		Body:   "Request processed",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	errorMessage := Message{
		Status: "Request Failed",
		Body:   "The API is at capacity",
	}
	jsonError, err := json.Marshal(errorMessage)
	if err != nil {
		log.Fatalf("Failed to encode rate limit error message: %v", err)
	}

	
	tollboothLimiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Hour,
	})

	tollboothLimiter.SetMessageContentType("application/json")
	tollboothLimiter.SetMessage(string(jsonError))

	http.Handle("/ping", tollbooth.LimitFuncHandler(tollboothLimiter, endPointHandler))

	fmt.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
