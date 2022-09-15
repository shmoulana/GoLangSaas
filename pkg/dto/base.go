package dto

type BaseResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}
