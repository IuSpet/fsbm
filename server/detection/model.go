package detection

type uploadResultRequest struct {
	ShopId     int64           `json:"shop_id"`
	DeviceId   int64           `json:"device_id"`
	VideoPath  string          `json:"video_path"`
	Detections []detectionInfo `json:"detections"`
}

type detectionInfo struct {
	At          int64  `json:"at"`
	FrameCnt    int64  `json:"frame_cnt"`
	ImgPath     string `json:"img_path"`
	IdentifyCnt int64  `json:"identify_cnt"`
	WearHatCnt  int64  `json:"wear_hat_cnt"`
	NoHatCnt    int64  `json:"no_hat_cnt"`
	ExtraJson   string `json:"extra_json"`
}

type getDeviceInfoRequest struct {
	ShopName   string `json:"shop_name"`
	UserName   string `json:"user_name"`
	UserEmail  string `json:"user_email"`
	DeviceName string `json:"device_name"`
}

type getDeviceInfoResponse struct {
	ShopId   int64        `json:"shop_id"`
	ShopName string       `json:"shop_name"`
	List     []deviceInfo `json:"list"`
}

type deviceInfo struct {
	ShopId     int64  `json:"shop_id"`
	DeviceId   int64  `json:"device_id"`
	DeviceName string `json:"device_name"`
}
