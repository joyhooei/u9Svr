package createOrder

import (
	//"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	// "strings"
	"u9/models"
	// "u9/tool"
)

type YYH struct {
	Cr
}

func (this *YYH) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

	if err = this.Cr.Initial(lr, orderId, nil, nil, extParamStr, channelParams, ctx); err != nil {
		return
	}
	beego.Trace("extParamStr:",extParamStr)
	if(extParamStr!=""){
		this.Result = this.parsePayKey(extParamStr)
		beego.Trace("Result:",this.Result)
	}
	// content := fmt.Sprintf(this.extParamStr, this.orderId)
	// content = strings.Replace(content, `\`, ``, -1) //去json中的`\`
	return nil
}

func (this *YYH) InitParam() (err error) {
	return nil
}

func (this *YYH) GetResponse() (err error) {
	return nil
}

func (this *YYH) ParseChannelRet() (err error) {
	return nil
}

func (this *YYH) GetResult() (ret string) {
	format := "getResult: result:%s"
	msg := fmt.Sprintf(format, this.Result)
	beego.Trace(msg)
	return this.Result
}
