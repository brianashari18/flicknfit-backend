package utils

// ErrorResponse is a generic struct for sending a consistent
// JSON error message to the client.
type ErrorResponse struct {
	Message string `json:"message"`
}
