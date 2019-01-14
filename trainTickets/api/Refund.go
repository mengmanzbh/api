package api

import (
    "bytes"
	"github.com/gin-gonic/gin" 
    "fmt"
    "trainTickets/utils"
    "net/url"
    "encoding/json"
    "math/rand"
    "time"
    "strconv"
    "net/http"
    "io/ioutil"
)
// const APPKEY = "5b433b1f92d41bba340a5bb47464ce32" //您申请的APPKEY
//线上退票
func Refund(ctx *gin.Context) {
	orderid := ctx.PostForm("orderid")
	tickets := ctx.PostForm("tickets")
    totalprice_origin := ctx.PostForm("totalprice")
    code_origin := ctx.PostForm("code")
    
    code := ctx.PostForm("code")
    token := getAccess(code)//根据前端传来的code获取token
    _,_,_,_,istoken := utils.GetUserByAccess(token,ctx)
    if !istoken{
        fmt.Println("token无效")
        return
    }
	//请求地址
	juheURL :="http://op.juhe.cn/trainTickets/refund"
	
	//初始化参数
	param:=url.Values{}
	
	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("dtype","json") //返回的格式，json或xml，默认json
	param.Set("key",APPKEY)
	param.Set("orderid",orderid) //订单号
	param.Set("tickets",tickets) //要退的火车票的信息。请注意，退票针对乘客（火车票），不针对订单，所以需要您提供乘客（车票）信息。所以一个订单下可以存在已出票和已退票的乘客。json格式的字符串，如：[{"ticket_no":"E374529237102001D","passengername":"王二","passporttypeseid":"1","passportseno":"370827199109121234"}]各个字段的解释：ticket_no：票号，是票号，不是12306的订单号（ordernumber），请先看常见问题中对各个订单号的说明passengername：该票号对应的乘客passporttypeseid：证件类型，1:二代身份证,2:一代身份证,C:港澳通行证,G:台湾通 行证,B:护照passportseno：证件号码
	
	
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
	    //返回前端数据
		ctx.JSON(200, gin.H{
			"error_code": netReturn["error_code"],
			"message": netReturn["reason"],
			"result":netReturn["result"],
		})
        

        fmt.Print(netReturn)
        if netReturn["result"] != nil{

			status := netReturn["result"].(map[string]interface{})["status"].(string)
			fmt.Print("退票状态:",status)

			//7：（有乘客）退票成功*/
			if status == "7" {
					//退票成功，给用户退币
					fmt.Print("退票成功，给用户加币")
					//****************给用户退币****************//
			                token := getAccess(code_origin)//根据前端传来的code获取token
			                lastprice :=  getLastprice()//请求接口获取最新价格
			                /**************将价格转成浮点数**************/
			                total, _ := strconv.ParseFloat(totalprice_origin, 64)
			                last, _ := strconv.ParseFloat(lastprice, 64)
			                resultdata := total/last
			                resultstr := fmt.Sprintf("%.4f",resultdata)
			                resultNum, _:= strconv.ParseFloat(resultstr, 64)
			                fmt.Println("resultNum:",resultNum)
			                /**************将价格转成浮点数**************/
    		                   isSuccess :=addNum(token,resultNum)
    		                   fmt.Println("加币结果:",isSuccess)
    		                   if isSuccess == true {
    		                   	fmt.Println("加币成功,钱已经退回给用户")
    		                   }else{
    		                   	fmt.Println("加币失败,钱没有退回给用户")
    		                   }
					//****************给用户退币****************//
			}else{
		         	fmt.Print("退票不成功")
			}
			
        } 


	}
}

//获取最新汇率
func getLastprice()(lastprice string) {
    huili := make(map[string]interface{})
    huili["currency_mark"] = "USDT"
    huili["currency_trade_mark"] = "TGV"
    bytesData, err := json.Marshal(huili)
    if err != nil {
        fmt.Println(err.Error() )
        return
    }
    reader := bytes.NewReader(bytesData)
    url := "https://m.51tg.vip/api/currency/currency/getLastData"
    request, err := http.NewRequest("POST", url, reader)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    request.Header.Set("Content-Type", "application/json;charset=UTF-8")
    client := http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    respBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    var data string
    //返回结果
    var netReturn map[string]interface{}
    json.Unmarshal(respBytes,&netReturn)
    if netReturn["isSuccess"]==true{
        last_price := netReturn["data"].(map[string]interface{})["lastData"].(map[string]interface{})["last_price"]
        // fmt.Println("获取最新价是:\r\n",last_price)
        data = last_price.(string)
    }else{
        fmt.Println("获取最新价失败")
    }
    return data
}


//根据前端传来的code获取token
func getAccess(code string)(token string) {
    data := make(map[string]interface{})
    data["code"] = code
    bytesData, err := json.Marshal(data)
    if err != nil {
        fmt.Println(err.Error() )
        return
    }
    reader := bytes.NewReader(bytesData)
    url := "http://47.52.47.73:8112/auth/auth/getAccess"
    request, err := http.NewRequest("POST", url, reader)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    request.Header.Set("Content-Type", "application/json;charset=UTF-8")
    client := http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    respBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    var tokendata string
    //返回结果
    var netReturn map[string]interface{}
    json.Unmarshal(respBytes,&netReturn)
    if netReturn["isSuccess"]==true{
        fmt.Printf("获取token:\r\n%v",netReturn["data"].(map[string]interface{})["auth"].(map[string]interface{})["access"])
        access := netReturn["data"].(map[string]interface{})["auth"].(map[string]interface{})["access"]
        tokendata = access.(string)
    }else{
        fmt.Println("获取token失败")
    }
    
    return tokendata
}
func addNum(taken string,num float64)(y bool) {
    data := make(map[string]interface{})
    data["access"] = taken
    data["mark"] = "TGV"
    data["num"] = num
    data["order_no"] = GetRandomString(6)
    bytesData, err := json.Marshal(data)
    if err != nil {
        fmt.Println(err.Error() )
        return
    }
    reader := bytes.NewReader(bytesData)
    url := "http://47.52.47.73:8112/auth/auth/subNum"
    request, err := http.NewRequest("POST", url, reader)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    request.Header.Set("Content-Type", "application/json;charset=UTF-8")
    client := http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    respBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    //返回结果
    var netReturn map[string]interface{}
    json.Unmarshal(respBytes,&netReturn)
     
    var isSuccess bool
    
    isSuccess = netReturn["isSuccess"].(bool)
    return isSuccess
}

//获取当前时间
func GetTime() string {
    const shortForm = "20060102150405"
    t := time.Now()
    temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
    str := temp.Format(shortForm)
    return str
}

// 随机生成置顶位数的大写字母和数字的组合
func  GetRandomString(l int) string {
    //str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    str := "0123456789"
    bytes := []byte(str)
    result := []byte{}
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < l; i++ {
        result = append(result, bytes[r.Intn(len(bytes))])
    }
    return GetTime() + string(result)
}
