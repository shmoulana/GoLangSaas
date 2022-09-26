package dto

type BaseResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

type ErrorResponse struct {
	HTTPCode int32  `json:"http_code"`
	Code     int32  `json:"code"`
	Message  string `json:"message"`
}
