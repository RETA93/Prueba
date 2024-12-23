package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Hola Mundo!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
