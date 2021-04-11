package util

import "bytes"

// 头像相关方法

type AvatarHandler struct {
	Avatar        *bytes.Buffer
	ContentLength int64
	ContentType   string
	ExtraHeaders  map[string]string
}

func NewAvatarHandler(avatar []byte) *AvatarHandler {
	return &AvatarHandler{
		Avatar:        bytes.NewBuffer(avatar),
		ContentLength: int64(len(avatar)),
		ContentType:   "image/png",
		ExtraHeaders:  nil,
	}
}
func (h *AvatarHandler) SetHeaders(key, value string) {
	h.ExtraHeaders[key] = value
}

func (h *AvatarHandler) Read(p []byte) (int, error) {
	return h.Avatar.Read(p)
}
