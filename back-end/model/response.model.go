package model

type StatusEnum string

const (
	SUCCESS StatusEnum = "SUCCESS"
	FAILED  StatusEnum = "FAILED"
)

type ErrorResponse struct {
	Message string     `json:"message"`
	Status  StatusEnum `json:"status"`
	Code    int        `json:"error_code"`
}

type SuccessResponse struct {
	Message string     `json:"message"`
	Status  StatusEnum `json:"status"`
	Code    int        `json:"code"`
	Data    any        `json:"data"`
}
