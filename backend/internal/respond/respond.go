package respond

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON writes any payload with the given status code.
func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("respond: encode failed: %v", err)
	}
}

// errorBody is the single error envelope the whole API uses: {"error": "..."}.
// Fields is populated only for validation failures.
type errorBody struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// Error writes a plain JSON error.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, errorBody{Error: message})
}

// ValidationError writes a 422 listing exactly which fields were rejected, so
// the client can mark the offending inputs instead of showing a generic toast.
func ValidationError(w http.ResponseWriter, fields map[string]string) {
	JSON(w, http.StatusUnprocessableEntity, errorBody{
		Error:  "validation failed",
		Fields: fields,
	})
}
