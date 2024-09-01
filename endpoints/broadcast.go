package endpoints

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"websocket-gateway/ws"
)

type BroadcastRequest struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

var allowedIPs = strings.Split(os.Getenv("ALLOWED_IP_CIDR"), ",")

func isAllowedIP(remoteAddr string) bool {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		ip = remoteAddr
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		log.Printf("Invalid IP address: %v", ip)
		return false
	}

	for _, cidr := range allowedIPs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Printf("Error parsing CIDR: %v", err)
			continue
		}
		if ipNet.Contains(parsedIP) {
			return true
		}
	}
	return false
}

func HandleBroadcast(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	if !isAllowedIP(clientIP) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}

	var req BroadcastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if req.Channel == "" || req.Message == "" {
		http.Error(w, "Missing channel or message", http.StatusBadRequest)
		return
	}

	BroadcastMessage(req.Channel, []byte(req.Message))

	w.WriteHeader(http.StatusOK)
}

func BroadcastMessage(channel string, message []byte) {
	hub.Broadcast <- ws.BroadcastMessage{Channel: channel, Message: message}
}
