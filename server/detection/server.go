package detection

import (
	"fmt"
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"fsbm/util/redis"
	"github.com/gin-gonic/gin"
	"time"
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
	var rows []db.DetectionResultRecord
	for _, item := range req.Detections {
		rows = append(rows, db.DetectionResultRecord{
			SrcDeviceID: req.DeviceId,
			SrcShopID:   req.ShopId,
			VideoPath:   req.VideoPath,
			At:          item.At,
			FrameCnt:    item.FrameCnt,
			ImgPath:     item.ImgPath,
			IdentifyCnt: item.IdentifyCnt,
			WearHatCnt:  item.WearHatCnt,
			NoHatCnt:    item.NoHatCnt,
			ExtraJson:   item.ExtraJson,
			Status:      db.DetectionResultRecordStatus_NotScanYet,
		})
	}
	if len(rows) == 0 {
		util.EndJson(ctx, nil)
	}
	err = db.SaveDetectionResultRecordRows(rows)
	if err != nil {
		logs.CtxError(ctx, "save detection result error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	now := time.Now()
	key1 := fmt.Sprintf(util.DashboardRecordCnt, now.Format(util.YMD))
	_ = redis.IncrByWithRetry(ctx, key1, int64(len(rows)))
	key2 := fmt.Sprintf(util.DashboardLatestRecord, now.Format(util.YMD))
	_ = redis.SetWithRetry(ctx, key2, now.Format(util.YMDHMS)[12:], 0)
	util.EndJson(ctx, nil)
}

// 获取设备信息
func GetDeviceInfoServer(ctx *gin.Context) {
	req := &getDeviceInfoRequest{}
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	shopInfoList, err := getShopInfo(req.UserEmail, req.UserName, req.ShopName)
	if err != nil {
		logs.CtxError(ctx, "get shop info error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if len(shopInfoList) == 0 {
		logs.CtxInfo(ctx, "empty shop info list.")
		util.ErrorJson(ctx, util.ShopNotFound, "找不到店铺信息")
		return
	}
	if len(shopInfoList) > 1 {
		logs.CtxInfo(ctx, "shop info not unique.")
		util.ErrorJson(ctx, util.ShopNotUniq, "店铺不唯一")
		return
	}
	shopInfo := shopInfoList[0]
	deviceList, err := db.GetMonitorListByShopId(shopInfo.ID)
	if err != nil {
		logs.CtxError(ctx, "get device info list error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	if req.DeviceName != "" {
		var tmp []db.MonitorList
		for _, item := range deviceList {
			if item.Name == req.DeviceName {
				tmp = append(tmp, item)
			}
		}
		deviceList = tmp
	}
	var list []deviceInfo
	for _, row := range deviceList {
		list = append(list, deviceInfo{
			ShopId:     shopInfo.ID,
			DeviceId:   row.ID,
			DeviceName: row.Name,
		})
	}
	rsp := getDeviceInfoResponse{
		ShopId:   shopInfo.ID,
		ShopName: shopInfo.Name,
		List:     list,
	}
	util.EndJson(ctx, rsp)
}
