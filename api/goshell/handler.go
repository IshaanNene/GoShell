package goshell

import (
	"net/http"
	"encoding/json"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Hello, GoShell!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
