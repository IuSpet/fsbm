package userAccount

import (
	"fsbm/util"
	"fsbm/util/logs"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func SetAvatarServer(ctx *gin.Context) {
	bodyReader := ctx.Request.Body
	data, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		logs.CtxError(ctx, "read body error. err: %+v", err)
		util.ErrorJson(ctx, util.ParamError, "请求内容读取失败")
		return
	}
	if len(data) > 64*64 {
		logs.CtxError(ctx, "body too large. len: %d", len(data))
		util.ErrorJson(ctx, util.ParamError, "内容太大")
		return
	}
	email := ctx.GetHeader("email")

}
