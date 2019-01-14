package api

import (
	"github.com/gin-gonic/gin" 
    // "fmt"
    "trainTickets/utils"
    // "net/url"
    "encoding/json"
    // "bytes"
    // "io/ioutil"
    // "net/http"
    // "strconv"
    "github.com/streadway/amqp"
    "log"
)
// rabbitmq使用的错误输出
func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

//站点简码查询
func Pay(ctx *gin.Context) {
	//获取参数
	orderid_origin := ctx.PostForm("orderid")
    totalprice_origin := ctx.PostForm("totalprice")
    code_origin := ctx.PostForm("code")

    code := ctx.PostForm("code")
    token := getAccess(code)//根据前端传来的code获取token
    _,_,_,_,istoken := utils.GetUserByAccess(token,ctx)
    if !istoken{
        fmt.Println("token无效")
        return
    }

    /*****************数据发送到消息队列***************/
    conn, err := amqp.Dial("amqp://guest:guest@3.81.214.206:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "pay_queue", // name
        true,         // durable
        false,        // delete when unused
        false,        // exclusive
        false,        // no-wait
        nil,          // arguments
    )
    failOnError(err, "Failed to declare a queue")

    orderid := orderid_origin
    totalprice := totalprice_origin
    code := code_origin
    data := make(map[string]interface{})
    data["orderid"] = orderid
    data["totalprice"] = totalprice
    data["code"] = code
    
    b,_:=json.Marshal(data)

    body := b
    err = ch.Publish(
        "",     // exchange
        q.Name, // routing key
        false,  // mandatory
        false,
        amqp.Publishing{
            DeliveryMode: amqp.Persistent,
            ContentType:  "text/plain",
            Body:         []byte(body),
        })

    if err!=nil{
        ctx.JSON(404, gin.H{
            "error_code": "1",
            "message": "数据提交失败",
        })

    }else{
        ctx.JSON(200, gin.H{
            "error_code": "0",
            "message": "数据提交成功",
        })
    }

    failOnError(err, "Failed to publish a message")
    log.Printf(" [x] Sent %s", body)
    /*****************数据发送到消息队列完毕***************/

}