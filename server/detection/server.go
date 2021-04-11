package detection

import (
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"github.com/gin-gonic/gin"
)

// 上传检测结果
func UploadDetectionResultServer(ctx *gin.Context) {
	var req uploadResultRequest
	err := ctx.Bind(&req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	var rows []db.DetectionResult
	for _, item := range req.Detections {
		rows = append(rows, db.DetectionResult{
			SrcVideoID:  req.VideoID,
			SrcDeviceID: req.DeviceID,
			At:          item.At,
			FrameCnt:    item.FrameCnt,
			Path:        item.Path,
			IdentifyCnt: item.IdentifyCnt,
			WearHatCnt:  item.WearHatCnt,
			NoHatCnt:    item.NoHatCnt,
		})
	}
	err = db.SaveDetectionResultRows(rows)
	if err != nil {
		logs.CtxError(ctx, "save detection result error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	util.EndJson(ctx, nil)
}
