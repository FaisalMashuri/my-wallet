package shared

type BaseReponse struct {
	Code    string      `json:"responseCode"`
	Status  string      `json:"responseStatus"`
	Message string      `json:"responseMessage,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(status string, message string, data interface{}) *BaseReponse {
	return &BaseReponse{
		Code:    "0000",
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(code string, status string, message string) *BaseReponse {
	return &BaseReponse{
		Code:    code,
		Status:  status,
		Message: message,
	}
}
