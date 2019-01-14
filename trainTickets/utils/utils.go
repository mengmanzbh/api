package utils

import (
	"time"
	"math/rand"
    "net/http"
    "net/url"
    "io/ioutil"
    "bytes"
    "fmt"
    "crypto/md5"
    "encoding/hex"
    "github.com/gin-gonic/gin" 
    "encoding/json"
)

// 根据授权码获取用户信息
func GetUserByAccess(access string,ctx *gin.Context)(x float64,y string,z string,q string,istoken bool) {

  var islogin bool
  data := make(map[string]interface{})
    data["access"] = access
    bytesData, err := json.Marshal(data)
    if err != nil {
        fmt.Println(err.Error() )
        return
    }
    reader := bytes.NewReader(bytesData)
    url := "http://47.52.47.73:8112/auth/auth/GetUserByAccess"
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
    var customer_id float64
    var realname string
    var nickname string
    var cellphone string

    //返回结果
    var netReturn map[string]interface{}
    json.Unmarshal(respBytes,&netReturn)
    if netReturn["isSuccess"]==true{
        userdata := netReturn["data"]
        islogin = true
        // fmt.Printf("获取用户信息:\r\n%v",userdata)
        customer_id = userdata.(map[string]interface{})["customer_id"].(float64)
        realname = userdata.(map[string]interface{})["realname"].(string)
        nickname = userdata.(map[string]interface{})["nickname"].(string)
        cellphone = userdata.(map[string]interface{})["cellphone"].(string)
        
    }else{
        fmt.Println("获取token失败")
        islogin = false
        ctx.JSON(404, gin.H{
            "error_code": "404",
            "message": "token失效，请重新登录",
        })
        
    }
    return customer_id,realname,nickname,cellphone,islogin
}

// 生成32位MD5
func MD5(text string) string {
    ctx := md5.New()
    ctx.Write([]byte(text))
    return hex.EncodeToString(ctx.Sum(nil))
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

// get 网络请求
func Get(apiURL string,params url.Values)(rs[]byte ,err error){
    var Url *url.URL
    Url,err=url.Parse(apiURL)
    if err!=nil{
        fmt.Printf("解析url错误:\r\n%v",err)
        return nil,err
    }
    //如果参数中有中文参数,这个方法会进行URLEncode
    Url.RawQuery=params.Encode()
    resp,err:=http.Get(Url.String())
    if err!=nil{
        fmt.Println("err:",err)
        return nil,err
    }
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
}
 
// post 网络请求 ,params 是url.Values类型
func Post(apiURL string, params url.Values)(rs[]byte,err error){
    resp,err:=http.PostForm(apiURL, params)
    if err!=nil{
        return nil ,err
    }
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
}
