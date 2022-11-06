package http_response

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func NewErrorResp(code int, msg string) ErrorResponse {
	resp := ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = msg
	return resp
}
