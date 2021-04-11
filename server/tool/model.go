package tool

type toolCommonRequest struct {
	Email string `json:"email"`
}

type saveImgResponse struct {
	Path string `json:"path"`
}
