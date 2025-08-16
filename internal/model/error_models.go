package model

type APIError struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}
