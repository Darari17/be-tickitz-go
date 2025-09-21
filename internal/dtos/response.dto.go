package dtos

type Response struct {
	Code    int         `json:"code" example:"200"`
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"request berhasil"`
	Data    interface{} `json:"data,omitempty"`
}
