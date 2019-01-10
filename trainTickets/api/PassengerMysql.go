package api

import (
	"github.com/gin-gonic/gin" 
    "fmt"
    "trainTickets/utils"
    "net/url"
    "encoding/json"
)
const (
    DB_Driver = "root:my-secret-pw@tcp(3.81.214.206:3306)/data?charset=utf8"
)
func OpenDB() (success bool, db *sql.DB) {
    var isOpen bool
    db, err := sql.Open("mysql", DB_Driver)
    if err != nil {
        isOpen = false
    } else {
        isOpen = true
    }
    CheckErr(err)
    return isOpen, db
}
//插入数据
func InsertPassengerToDB(ctx *gin.Context) {
	        ctx.JSON(200, gin.H{
            "error_code": "0",
            "message": "添加乘客成功",
            })

}
//查询数据
func QueryPassengerFromDB(ctx *gin.Context) {
	        ctx.JSON(200, gin.H{
            "error_code": "0",
            "message": "查询乘客成功",
            })

}
//更新数据
func UpdatePassengerToDB(ctx *gin.Context) {
	        ctx.JSON(200, gin.H{
            "error_code": "0",
            "message": "更新乘客成功",
            })

}
//删除数据
func DeletePassengerFromDB(ctx *gin.Context) {
	        ctx.JSON(200, gin.H{
            "error_code": "0",
            "message": "删除乘客成功",
            })

}

func CheckErr(err error) {
    if err != nil {
        panic(err)
        fmt.Println("err:", err)
    }
}

func GetTime() string {
    const shortForm = "2006-01-02 15:04:05"
    t := time.Now()
    temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
    str := temp.Format(shortForm)
    fmt.Println(t)
    return str
}

func GetMD5Hash(text string) string {
    haser := md5.New()
    haser.Write([]byte(text))
    return hex.EncodeToString(haser.Sum(nil))
}

func GetNowtimeMD5() string {
    t := time.Now()
    timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
    return GetMD5Hash(timestamp)
}
func OpenAndInsertToDB() {
    opend, db := OpenDB()
    if opend {
        fmt.Println("open success")
        // DeleteFromDB(db, 10)
        //QueryFromDB(db)
        //DeleteFromDB(db, 1)
        //UpdateDB(db, 5)
        //UpdateUID(db, 5)
        //UpdateTime(db, 4)
        // insertToDB(db)
    } else {
        fmt.Println("open faile:")
    }

}