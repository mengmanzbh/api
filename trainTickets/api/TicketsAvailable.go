package api

import (
	"github.com/gin-gonic/gin" 
    "fmt"
    "trainTickets/utils"
    "net/url"
    "encoding/json"
)

// const APPKEY = "5b433b1f92d41bba340a5bb47464ce32" //您申请的APPKEY

//余票查询
func TicketsAvailable(ctx *gin.Context) {
    train_date := ctx.PostForm("train_date")
    from_station := ctx.PostForm("from_station")
    to_station := ctx.PostForm("to_station")
    
    code := ctx.PostForm("code")
    token := getAccess(code)//根据前端传来的code获取token
    _,_,_,_,istoken := utils.GetUserByAccess(token,ctx)
    if !istoken{
        fmt.Println("token无效")
        return
    }
    	//请求地址
    juheURL :="http://op.juhe.cn/trainTickets/ticketsAvailable"
 
    //初始化参数
    param:=url.Values{}
 
    //配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
    param.Set("dtype","json") //返回的格式，json或xml，默认json
    param.Set("key",APPKEY)
    param.Set("train_date",train_date) //发车日期，如：2015-07-01（务必按照此格式）
    param.Set("from_station",from_station) //出发站简码，如：BJP
    param.Set("to_station",to_station) //到达站简码，如：SHH
 
 
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
