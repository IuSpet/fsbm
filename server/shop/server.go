package shop

import (
	"encoding/json"
	"errors"
	"fsbm/db"
	"fsbm/util"
	"fsbm/util/logs"
	"github.com/gin-gonic/gin"
	"time"
)

// 获取店铺列表
func GetShopListServer(ctx *gin.Context) {
	req := newGetShopListRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
}

// 注册新店铺
func AddShopServer(ctx *gin.Context) {
	req := newAddShopRequest()
	err := ctx.Bind(req)
	if err != nil {
		logs.CtxError(ctx, "bind req error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	logs.CtxInfo(ctx, "req: %+v", req)
	if !isValidShopInfo(req) {
		logs.CtxInfo(ctx, "shop info invalid")
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}
	email := ctx.GetString("email")
	user, err := db.GetUserByEmail(email)
	if err != nil {
		logs.CtxError(ctx, "get user by email error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	noticeCfg, err := json.Marshal(req.NoticeCfg)
	if err != nil {
		logs.CtxError(ctx, "marshal error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "参数错误")
		return
	}

	row := &db.ShopList{
		Name:         req.Name,
		UserID:       user.ID,
		Addr:         req.Addr,
		NoticeConfig: string(noticeCfg),
		Status:       0,
		Remark:       req.Remark,
	}
	err = db.SaveShopListRow(row)
	if err != nil {
		logs.CtxError(ctx, "save shop_list row error. err: %+v", err)
		util.ErrorJson(ctx, util.DbError, "数据库错误")
		return
	}
	util.EndJson(ctx, nil)
}

func getShopListData(req *getShopListRequest) ([]shopInfo, error) {
	begin, err1 := time.Parse(util.YMD, req.CreateBegin)
	end, err2 := time.Parse(util.YMD, req.CreateEnd)
	if err1 != nil || err2 != nil {
		return nil, errors.New("time parse error")
	}
	rows, err := getShopListRows(req.Name, req.Addr, req.Admin, begin, end)
	if err != nil {
		return nil, err
	}
	_ = rows
	return nil, nil
}

func newGetShopListRequest() *getShopListRequest {
	return &getShopListRequest{
		Name:        "",
		Admin:       "",
		Addr:        "",
		CreateBegin: time.Now().Format(util.YMD),
		CreateEnd:   time.Now().AddDate(0, 0, -7).Format(util.YMD),
	}
}

func newAddShopRequest() *addShopRequest {
	return &addShopRequest{
		Name:      "",
		Addr:      "",
		Latitude:  0,
		Longitude: 0,
		NoticeCfg: make(map[string]db.ShopNoticeConfigBase),
		Remark:    "",
	}
}

func isValidShopInfo(req *addShopRequest) bool {
	if req.Name == "" {
		return false
	}
	if req.Addr == "" {
		return false
	}
	return true
}
