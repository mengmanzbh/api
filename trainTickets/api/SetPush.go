package api

import (
	"github.com/gin-gonic/gin" 
    "fmt"
    "trainTickets/utils"
    "net/url"
    "encoding/json"
)
// const APPKEY = "5b433b1f92d41bba340a5bb47464ce32" //您申请的APPKEY
//设置回调

//订单提交，占座回调地址，处理完占座后会将订单信息回调到此地址
func Submit_callback(ctx *gin.Context) {
	fmt.Println("Submit_callback")
	data := ctx.PostForm("data")
	fmt.Println(data)
	fmt.Println("success")
	ctx.JSON(200, gin.H{
	"message": "success",
	})
}

//出票回调地址，处理完出票请求后会将订单信息回调到此地址
func Pay_callback(ctx *gin.Context){
	fmt.Println("Pay_callback")
	data := ctx.PostForm("data")
	fmt.Println(data)
	ctx.JSON(200, gin.H{
	"message": "success",
	})
}

//退款回调地址，有退款发生时，会将订单信息回调到此地址
func Refund_callback(ctx *gin.Context){
	fmt.Println("Refund_callback")
	data := ctx.PostForm("data")
	fmt.Println(data)
	ctx.JSON(200, gin.H{
	"message": "success",
	})
}

func SetPush() {

    baseurl := "http://3.81.214.206:9000"
	submit_callback := baseurl + "/trainTickets/submit_callback"
	pay_callback := baseurl + "/trainTickets/pay_callback"
	refund_callback := baseurl + "trainTickets/refund_callback"
	
	//请求地址
	juheURL :="http://op.juhe.cn/trainTickets/setPush"
	
	//初始化参数
	param:=url.Values{}
	
	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("dtype","json") //返回的格式，json或xml，默认json
	param.Set("key",APPKEY)
	param.Set("submit_callback",submit_callback) //占座回调地址，处理完占座后会将订单信息回调到此地址
	param.Set("pay_callback",pay_callback) //出票回调地址，处理完出票请求后会将订单信息回调到此地址
	param.Set("refund_callback",refund_callback) //退款回调地址，有退款发生时，会将订单信息回调到此地址
	
	
	//发送请求
	data,err:=utils.Post(juheURL,param)
	if err!=nil{
		fmt.Errorf("请求失败,错误信息:\r\n%v",err)
	}else{
		var netReturn map[string]interface{}
		json.Unmarshal(data,&netReturn)
	    fmt.Println("回调成功,返回信息:\r\n%v",netReturn)
	}
}

