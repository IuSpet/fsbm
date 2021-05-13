package cronjob

import (
	"context"
	"fsbm/db"
	"fsbm/util"
)

// 监控任务

func infoNoticeTask(ctx context.Context) error {
	// 查询未扫描过的记录
	recordList, err := db.GetDetectionResultsByNoticeLevel(util.InfoNotice)
	if err != nil {
		return err
	}
	// 根据shop_id聚合
	shopRecords := make(map[int64][]db.DetectionResult)
	for i := range recordList {
		shopRecords[recordList[i].SrcShopID] = append(shopRecords[recordList[i].SrcShopID], recordList[i])
	}
	// 对shop_id遍历
	for shopId, record := range shopRecords {
		// 查询shop报警阈值

	}
}
