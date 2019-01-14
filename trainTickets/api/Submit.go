package api

import (
	"github.com/gin-gonic/gin" 
    "fmt"
    "trainTickets/utils"
    "net/url"
    "encoding/json"
    // "bytes"
    // "net/http"
    // "io/ioutil"
    // "strconv"
)
// const APPKEY = "5b433b1f92d41bba340a5bb47464ce32" //您申请的APPKEY
//提交订单
func Submit(ctx *gin.Context) {
	
	code := ctx.PostForm("code")
	token := getAccess(code)//根据前端传来的code获取token
	customer_id,realname,nickname,cellphone,istoken := utils.GetUserByAccess(token,ctx)
	if !istoken{
		fmt.Println("token无效")
	    return
    }
	fmt.Println(customer_id)
	fmt.Println(realname)
	fmt.Println(nickname)
	fmt.Println(cellphone)
     
    randomstring := utils.GetRandomString(6)
	user_orderid := fmt.Sprintf("%s%v", randomstring, customer_id)

	fmt.Println("用户订单号",user_orderid)
	train_date := ctx.PostForm("train_date")
	is_accept_standing := ctx.PostForm("is_accept_standing")
	choose_seats := ctx.PostForm("choose_seats")
	from_station_code := ctx.PostForm("from_station_code")
	to_station_code := ctx.PostForm("to_station_code")
	checi := ctx.PostForm("checi")
	passengers := ctx.PostForm("passengers")
	//请求地址
	juheURL :="http://op.juhe.cn/trainTickets/submit"
	
	//初始化参数
	param:=url.Values{}
	
	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("dtype","json") //返回的格式，json或xml，默认json
	param.Set("key",APPKEY)
	param.Set("user_orderid",user_orderid) //您自定义的订单号，如：12345678，不要超过50个字符
	param.Set("train_date",train_date) //乘车日期，如：2015-07-01
	param.Set("is_accept_standing",is_accept_standing) //是否接受无座（最低座次无票时自动尝试抢无座票），传值yes则接受，传值no则不接受，默认yes
	param.Set("choose_seats",choose_seats) //需要选的座位 如：1A2B2C
	param.Set("from_station_code",from_station_code) //出发站简码，如：SZH
	param.Set("to_station_code",to_station_code) //到达站简码，如：SHH
	param.Set("checi",checi) //车次 如：G7027，请注意：出发站、到达站信息必须属实，例如G101车次不经过北京（经过北京南），出发站信息中不能填北京
	param.Set("passengers",passengers)//乘车人信息，每个订单最多5个乘客，json格式的字符串，如：[{"passengerid":1,"passengersename":"张三","piaotype":"1","piaotypename":"成人票","passporttypeseid":"1","passporttypeseidname":"二代身份证","passportseno":"420205199207231234","price":"763.5","zwcode":"M","zwname":"一等座"}]请想清楚，这里传json不代表所有的参数（如key、checi等）都是通过json传递过来！您测试的时候请将乘客信息换成真实信息，如使用非真实信息（张三、李四等）提交的订单过多，将会被禁止使用各个字段的解释：  "passengerid":乘客的顺序号，如：1，当一个订单有多个乘客时，用来唯一标识每个乘客，建议设为1-5（因每单最多5个乘客）的正整数；请自定义此参数。  "passengersename":乘车人姓名  "piaotype":如：1。其中，1 :成人票,2 :儿童票,4 :残军票  "piaotypename":如：成人票。票种名称，和上面的piaotype对应  "passporttypeseid": 如：1。其中，1:二代身份证,2:一代身份证,C:港澳通行证,B:护照,G:台湾通行证  "passporttypeseidname":如：二代身份证。证件类型名称，和上面的passporttypeseid对应  "passportseno":如：420205199207231234。乘客证件号码  "price": 票价，即当前乘客选择的座位的价格  "zwcode":如：1。表示座位编码，其中    F:动卧(新增),    9:商务座,    P:特等座,    M:一等座,    O（大写字母O，不是数字0）:二等座,    6:高级软卧,     4:软卧,    3:硬卧,    2:软座,    1:硬座。    注意:此处规则与12306不同，无座没有zwcode，当最低座位无票时,购买选该座位, 买下的就是无座。请务必阅读火车票订票接口常见问题中序号为7、8、9、9.1的内容  "zwname":如：硬座。表示座位名称，和上面的座位编码对应，注意：无座没有zwname，即不能由用户指定购买无座（无座仅作为最低座次的备选票）
    
    var juheorderid string
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
	
	    juheorderid = netReturn["result"].(map[string]interface{})["orderid"].(string)
		ctx.JSON(200, gin.H{
			"error_code": netReturn["error_code"],
			"message": netReturn["reason"],
			"result":netReturn["result"],
		})
	}

    /****************订单入库****************/
     fmt.Println("订单入库")


    opend, db := OpenDB()
    if opend {
        fmt.Println("open success")

	    nowTimeStr := GetTime()
	    stmt, err := db.Prepare("insert train_tickets_orders set realname=?,nickname=?,cellphone=?,customer_id=?,user_orderid=?")
	    CheckErr(err)
	    res, err := stmt.Exec(realname, nickname, cellphone, customer_id, user_orderid)
	    CheckErr(err)
	    id, err := res.LastInsertId()
	    CheckErr(err)
	    if err != nil {
	        fmt.Println("插入数据失败")
	    } else {
	        fmt.Println("插入数据成功：", id)
	    }

    } else {
        fmt.Println("open faile:")
    }


    /****************订单入库****************/
}



