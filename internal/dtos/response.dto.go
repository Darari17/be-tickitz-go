package dtos

type Response struct {
	Code    int         `json:"code" example:"200"`
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"request berhasil"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrResponse struct {
	Code    int         `json:"code" example:"500"`
	Success bool        `json:"success" example:"false"`
	Message string      `json:"message" example:"internal server error"`
	Data    interface{} `json:"data,omitempty"`
}
