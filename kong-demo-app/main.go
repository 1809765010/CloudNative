package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Message   string `json:"message"`
	Port      string `json:"port"`
	Timestamp string `json:"timestamp"`
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := Response{
		Message:   "Hello from Demo App!",
		Port:      ":" + os.Getenv("PORT"),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	port = ":" + port

	fmt.Printf("Starting server on port %s\n", port)

	// 注册路由
	http.HandleFunc("/", handleRoot)

	log.Fatal(http.ListenAndServe(port, nil))
}
