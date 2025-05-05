package models

type BaseResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
	ErrorCode int    `json:"errorCode"`
}