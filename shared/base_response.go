package shared

type BaseResponse struct {
	Code    string      `json:"responseCode"`
	Status  string      `json:"responseStatus"`
	Message string      `json:"responseMessage,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(status string, message string, data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:    "0000",
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(code string, status string, message string) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Status:  status,
		Message: message,
	}
}
