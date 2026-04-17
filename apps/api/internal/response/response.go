package response

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) ApiResponse {
	return ApiResponse{
		Code:    0,
		Message: "ok",
		Data:    data,
	}
}

func Error(message string) ApiResponse {
	return ApiResponse{
		Code:    1,
		Message: message,
	}
}

type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}
