package api

import (
	"github.com/gin-gonic/gin" 
    "fmt"
    // "trainTickets/utils"
    // "net/url"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    // "encoding/json"
    "time"
    "crypto/md5"
    "encoding/hex"
    "strconv"
)

//插入数据
func InsertSearchRecordToDB(ctx *gin.Context) {
    opend, db := OpenDB()
    if opend {
        fmt.Println("open success")
        
            ctx.JSON(200, gin.H{
            "error_code": "0",
            "message": "插入搜索记录成功",
            })
    } else {
        fmt.Println("open faile:")
    }

}
//查询数据
func QuerySearchRecordFromDB(ctx *gin.Context) {
	        ctx.JSON(200, gin.H{
            "error_code": "0",
            "message": "查询搜索记录成功",
            })
}
//更新数据
func UpdateSearchRecordToDB(ctx *gin.Context) {
   	        ctx.JSON(200, gin.H{
            "error_code": "0",
            "message": "更新搜索记录成功",
            })
}
//删除数据
func DeleteSearchRecordFromDB(ctx *gin.Context) {
	        ctx.JSON(200, gin.H{
            "error_code": "0",
            "message": "删除搜索记录成功",
            })

}