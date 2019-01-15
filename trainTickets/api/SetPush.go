package api

import (
	"github.com/gin-gonic/gin" 
    "fmt"
    "trainTickets/utils"
    "net/url"
    "encoding/json"
    // "os/exec"
    // "strings"
)
// const APPKEY = "5b433b1f92d41bba340a5bb47464ce32" //您申请的APPKEY
//设置回调

//订单提交，占座回调地址，处理完占座后会将订单信息回调到此地址
func Submit_callback(ctx *gin.Context) {
	fmt.Println("Submit_callback")
	data := ctx.PostForm("data")
    var dicdata map[string]interface{}
    json.Unmarshal([]byte(data),&dicdata)	
    fmt.Println(dicdata["from_station_name"])

	// fmt.Println("success")
	// ctx.JSON(200, gin.H{
	// "message": "success",
	// })


    /****************订单入库****************/
    //  fmt.Println("订单入库")

    //这个OpenDB()方法在PassengerMysql里面
    opend, db := OpenDB()
    if opend {
        fmt.Println("open success")

		stmt, err := db.Prepare("update train_tickets_orders set orderid=?,msg=?,orderamount=?,status=?,ordernumber=?,submit_time=?,deal_time=?,cancel_time=?,pay_time=?,finished_time=?,refund_time=?,juhe_refund_time=?,train_date=?,from_station_name=?,to_station_name=?,refund_money=? where user_orderid=?")
	    CheckErr(err)
	    res, err := stmt.Exec(dicdata["orderid"],dicdata["msg"],dicdata["orderamount"],dicdata["status"],dicdata["ordernumber"],dicdata["submit_time"],dicdata["deal_time"],dicdata["cancel_time"],dicdata["pay_time"],dicdata["finished_time"],dicdata["refund_time"],dicdata["juhe_refund_time"],dicdata["train_date"],dicdata["from_station_name"],dicdata["to_station_name"],dicdata["refund_money"],dicdata["user_orderid"])
	    affect, err := res.RowsAffected()
	    fmt.Println("更新数据：", affect)
	    CheckErr(err)

    } else {
        fmt.Println("open faile:")
    }
    /****************订单入库****************/

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

