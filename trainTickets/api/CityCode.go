package api

import (
	"github.com/gin-gonic/gin" 
    "fmt"
    "trainTickets/utils"
    "trainTickets/libs"
    "net/url"
    "encoding/json"
)
const APPKEY = "5b433b1f92d41bba340a5bb47464ce32" //您申请的APPKEY


// @Summary 站点简码查询
// @Description 站点简码查询
// @Accept json
// @Produce json
// @Param stationName path string true "站点名，如苏州、苏州北，不需要加“站”字"
// @Param all path string false "如果需要全部站点简码，请将此参数设为1"
// @Router /trainTickets/cityCode [post]

//站点简码查询
func CityCode(ctx *gin.Context) {
	stationName := ctx.PostForm("stationName")
    all := ctx.PostForm("all")
	//请求地址

    juheURL := libs.Conf.Read("api", "juhebaseurl") + "/trainTickets/cityCode"
 
    //初始化参数
    param:=url.Values{}
 
    //配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
    param.Set("dtype","json") //返回的格式，json或xml，默认json
    param.Set("key",APPKEY)
    param.Set("stationName",stationName) //站点名，如苏州、苏州北，不需要加“站”字
    param.Set("all",all) //如果需要全部站点简码，请将此参数设为1
 
 
    //发送请求
    data,err:=utils.Post(juheURL,param)
    if err!=nil{
        fmt.Errorf("请求失败,错误信息:\r\n%v",err)
        ctx.JSON(404, gin.H{
            "error_code": "404",
            "message": err,
            })
    }else{
        var netReturn map[string]interface{}
        json.Unmarshal(data,&netReturn)

        ctx.JSON(200, gin.H{
            "error_code": netReturn["error_code"],
            "message": netReturn["reason"],
            "result":netReturn["result"],
        })
    }
}

