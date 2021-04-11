package detection

type uploadResultRequest struct {
	DeviceID   int64           `json:"device_id"`
	VideoID    int64           `json:"video_id"`
	Detections []detectionInfo `json:"detections"`
}

type detectionInfo struct {
	At          int64  `json:"at"`
	FrameCnt    int64  `json:"frame_cnt"`
	Path        string `json:"path"`
	IdentifyCnt int64  `json:"identify_cnt"`
	WearHatCnt  int64  `json:"wear_hat_cnt"`
	NoHatCnt    int64  `json:"no_hat_cnt"`
}

