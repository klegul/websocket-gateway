package main

import (
	"log"
	"net/http"
	"websocket-gateway/endpoints"
)

func main() {
	endpoints.Run()

	http.HandleFunc("/ws", endpoints.HandleConnections)
	http.HandleFunc("/broadcast", endpoints.HandleBroadcast)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
