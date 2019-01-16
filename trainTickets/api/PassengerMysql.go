package api

import (
	"github.com/gin-gonic/gin" 
    "fmt"
    "trainTickets/utils"
    // "net/url"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    // "encoding/json"
    "time"
    "crypto/md5"
    "encoding/hex"
    "strconv"
    "trainTickets/libs"
)

func OpenDB() (success bool, db *sql.DB) {

    DB_Driver := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
        libs.Conf.Read("mysql", "username"),
        libs.Conf.Read("mysql", "password"),
        libs.Conf.Read("mysql", "host"),
        libs.Conf.Read("mysql", "dataname"))


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
//先检查数据是否存在，在插入
func CustomerDataisexist(ctx *gin.Context,uid string) (x bool){
    customer_id := uid

    code := ctx.PostForm("code")
    token := getAccess(code)//根据前端传来的code获取token
    _,_,_,_,istoken := utils.GetUserByAccess(token,ctx)
    if !istoken{
        fmt.Println("token无效")
        return
    }

    var isexist bool
    opend, db := OpenDB()
    if opend {
        fmt.Println("open success")
        /**********查询数据***********/
        rows, err := db.Query("SELECT * FROM train_ticket_search_record")
        CheckErr(err)
        defer rows.Close()
        
        //循环结果集 
        for rows.Next() {
            columns, _ := rows.Columns()

            scanArgs := make([]interface{}, len(columns))
            values := make([]interface{}, len(columns))

            for i := range values {
                scanArgs[i] = &values[i]
            }

            //将数据保存到 record 字典
            err = rows.Scan(scanArgs...)
            record := make(map[string]string)
            for i, col := range values {
                if col != nil {
                    record[columns[i]] = string(col.([]byte))
                }
            }
            // fmt.Println(record)
            //过滤数据
            if record["customer_id"] == customer_id{
                fmt.Println(record["customer_id"])  
                isexist = true 
            }else{
                isexist = false
            }
        }
        /**********查询数据***********/

    } else {
        fmt.Println("open faile:")
        ctx.JSON(404, gin.H{
        "error_code": "1",
        "message": "数据库连接失败,查询异常",
         })        
    }
    
    return isexist
}
//插入数据
func InsertSearchRecordToDB(ctx *gin.Context) {
    customer_id := "323232"
    opend, db := OpenDB()
    if opend {
        fmt.Println("open success")

        datain := CustomerDataisexist(ctx,customer_id)
        if datain {
             fmt.Println("该用户存在数据库")
            //存在数据库直接更新记录
            stmt, err := db.Prepare("update train_ticket_search_record set search_records=? where customer_id=?")
            CheckErr(err)
            res, err := stmt.Exec("上海-杭州", customer_id)
            affect, err := res.RowsAffected()
            fmt.Println("更新数据：", affect)
            CheckErr(err)

        }else{
             fmt.Println("该用户不存在数据库")

                /**********插入搜索记录*********/
                nowTimeStr := GetTime()
                stmt, err := db.Prepare("insert train_ticket_search_record set customer_id=?,search_records=?,update_time=?")
                CheckErr(err)
                res, err := stmt.Exec("323232", "上海-苏州", nowTimeStr)
                CheckErr(err)
                id, err := res.LastInsertId()
                CheckErr(err)
                if err != nil {
                    fmt.Println("插入数据失败")
                    ctx.JSON(200, gin.H{
                    "error_code": "1",
                    "message": "插入搜索记录失败",
                    })
                } else {
                    fmt.Println("插入数据成功：", id)
                    ctx.JSON(200, gin.H{
                    "error_code": "0",
                    "message": "插入搜索记录成功",
                    })
                }
                /**********插入搜索记录*********/
        }
        
       


    } else {
        fmt.Println("open faile:")
            ctx.JSON(404, gin.H{
            "error_code": "1",
            "message": "数据库打开失败",
            })
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



//插入数据
func InsertPassengerToDB(ctx *gin.Context) {

        passengerse_name := ctx.PostForm("passengerse_name")
        piao_type := ctx.PostForm("piao_type")
        piaotype_name := ctx.PostForm("piaotype_name")
        passporttypese_id := ctx.PostForm("passporttypese_id")
        passporttypeseid_name := ctx.PostForm("passporttypeseid_name")
        passportse_no := ctx.PostForm("passportse_no")

        code := ctx.PostForm("code")
        token := getAccess(code)//根据前端传来的code获取token
        var customerid string
        customer_id,_,_,_,istoken := utils.GetUserByAccess(token,ctx)
        //先检查token是否有效
        if !istoken{
            fmt.Println("token无效")
            return
        }
        customerid = fmt.Sprintf("%v",customer_id)


        uid := GetMD5Hash(passportse_no+customerid)
        nowTimeStr := GetTime()
         //打开数据库
         opend, db := OpenDB()
        if opend {
            fmt.Println("open success")
            
            datain := Dataisexist(ctx,uid)
            if datain {
                fmt.Println("存在")
                 ctx.JSON(200, gin.H{
                   "error_code": "1",
                   "message": "该乘客已经添加过了",
                })
                return
            }

            stmt, err := db.Prepare("insert passengers set passengerse_name=?,piao_type=?,piaotype_name=?,passporttypese_id=?,passporttypeseid_name=?,passportse_no=?,create_time=?,customer_id=?,uid=?")
            CheckErr(err)
            res, err := stmt.Exec(passengerse_name, piao_type, piaotype_name, passporttypese_id, passporttypeseid_name,passportse_no,nowTimeStr,customerid,uid)

            CheckErr(err)
            id, err := res.LastInsertId()
            CheckErr(err)
            if err != nil {
                fmt.Println("插入数据失败")
                ctx.JSON(200, gin.H{
                   "error_code": "1",
                    "message": "添加乘客异常，请稍后重试",
                })
            } else {
                fmt.Println("插入数据成功：", id)
                ctx.JSON(200, gin.H{
                   "error_code": "0",
                   "message": "添加乘客成功",
                })
            }
        } else {
            fmt.Println("open faile:")
            ctx.JSON(200, gin.H{
             "error_code": "1",
             "message": "添加乘客异常，请稍后重试",
            })
        }

}
//先检查数据是否存在，在插入
func Dataisexist(ctx *gin.Context,uid string) (x bool){
    customer_id := uid
    var isexist bool
    opend, db := OpenDB()
    if opend {
        fmt.Println("open success")
        /**********查询数据***********/
        rows, err := db.Query("SELECT * FROM passengers")
        CheckErr(err)
        defer rows.Close()
        
        //循环结果集 
        for rows.Next() {
            columns, _ := rows.Columns()

            scanArgs := make([]interface{}, len(columns))
            values := make([]interface{}, len(columns))

            for i := range values {
                scanArgs[i] = &values[i]
            }

            //将数据保存到 record 字典
            err = rows.Scan(scanArgs...)
            record := make(map[string]string)
            for i, col := range values {
                if col != nil {
                    record[columns[i]] = string(col.([]byte))
                }
            }
            // fmt.Println(record)
            //过滤数据
            if record["uid"] == customer_id{
                fmt.Println(record["uid"])  
                isexist = true 
            }else{
                isexist = false
            }
        }
        /**********查询数据***********/

    } else {
        fmt.Println("open faile:")
        ctx.JSON(404, gin.H{
        "error_code": "1",
        "message": "数据库连接失败,查询异常",
         })        
    }
    
    return isexist
}

//查询数据
func QueryPassengerFromDB(ctx *gin.Context) {

    code := ctx.PostForm("code")
    token := getAccess(code)//根据前端传来的code获取token
    var customerid string
    customer_id,realname,nickname,cellphone,istoken := utils.GetUserByAccess(token,ctx)
    if !istoken{
        fmt.Println("token无效")
        return
    }
    customerid = fmt.Sprintf("%v",customer_id)
    fmt.Println(customerid)
    fmt.Println(realname)
    fmt.Println(nickname)
    fmt.Println(cellphone)

    opend, db := OpenDB()
    if opend {
        fmt.Println("open success")
        /**********查询数据***********/
        rows, err := db.Query("SELECT * FROM passengers")
        CheckErr(err)
        defer rows.Close()
        var dataArray []map[string]string
        //循环结果集 
        for rows.Next() {
            columns, _ := rows.Columns()

            scanArgs := make([]interface{}, len(columns))
            values := make([]interface{}, len(columns))

            for i := range values {
                scanArgs[i] = &values[i]
            }

            //将数据保存到 record 字典
            err = rows.Scan(scanArgs...)
            record := make(map[string]string)
            for i, col := range values {
                if col != nil {
                    record[columns[i]] = string(col.([]byte))
                }
            }
            // fmt.Println(record)
            //过滤数据
            fmt.Println("record_customer_id:",record["customer_id"])
            if record["customer_id"] == customerid{
                fmt.Println(record["customer_id"])  
                dataArray = append(dataArray, record)
            }
        }

        ctx.JSON(200, gin.H{
             "error_code": "0",
             "message": "查询乘客成功",
             "result": dataArray,
        })

        /**********查询数据***********/

    } else {
        fmt.Println("open faile:")
        ctx.JSON(404, gin.H{
        "error_code": "1",
        "message": "数据库连接失败,查询异常",
         })        
    }
}
//更新数据
func UpdatePassengerToDB(ctx *gin.Context) {
    passengerse_name := ctx.PostForm("passengerse_name")
    piao_type := ctx.PostForm("piao_type")
    piaotype_name := ctx.PostForm("piaotype_name")
    passporttypese_id := ctx.PostForm("passporttypese_id")
    passporttypeseid_name := ctx.PostForm("passporttypeseid_name")
    passportse_no := ctx.PostForm("passportse_no")

    code := ctx.PostForm("code")
    token := getAccess(code)//根据前端传来的code获取token
    var customerid string
    customer_id,_,_,_,istoken := utils.GetUserByAccess(token,ctx)
    //先检查token是否有效
    if !istoken{
        fmt.Println("token无效")
        return
    }
    customerid = fmt.Sprintf("%v",customer_id)

    uid := GetMD5Hash(passportse_no+customerid)
    opend, db := OpenDB()
    if opend {
        fmt.Println("open success")

        //先判断数据是否存在
        datain := Dataisexist(ctx,uid)
        nowTimeStr := GetTime()
        if datain {
            fmt.Println("存在")

            stmt, err := db.Prepare("update passengers set passengerse_name=?, piao_type=?, piaotype_name=?, passporttypese_id=?, passporttypeseid_name=?, passportse_no=?,update_time=? where uid=?")
            CheckErr(err)
            res, err := stmt.Exec(passengerse_name,piao_type,piaotype_name,passporttypese_id,passporttypeseid_name,passportse_no,nowTimeStr,uid)
            affect, err := res.RowsAffected()
            fmt.Println("更新数据：", affect)
            CheckErr(err)

             ctx.JSON(200, gin.H{
               "error_code": "0",
               "message": "更新乘客信息成功",
             })
        }else{
            fmt.Println("不存在,修改的是身份证号码，重新插入记录")

            stmt, err := db.Prepare("insert passengers set passengerse_name=?,piao_type=?,piaotype_name=?,passporttypese_id=?,passporttypeseid_name=?,passportse_no=?,create_time=?,customer_id=?,uid=?")
            CheckErr(err)
            res, err := stmt.Exec(passengerse_name, piao_type, piaotype_name, passporttypese_id, passporttypeseid_name,passportse_no,nowTimeStr,customerid,uid)
            
            CheckErr(err)
            id, err := res.LastInsertId()
            CheckErr(err)
            if err != nil {
                fmt.Println("更新数据失败")
                ctx.JSON(200, gin.H{
                   "error_code": "1",
                    "message": "更新乘客异常，请稍后重试",
                })
            } else {
                fmt.Println("插入数据成功：", id)
                ctx.JSON(200, gin.H{
                   "error_code": "0",
                   "message": "更新乘客成功",
                })
            }

        }

    } else {
        fmt.Println("open faile:")
        ctx.JSON(200, gin.H{
        "error_code": "1",
        "message": "更新乘客信息异常",
        })
    }



}
//删除数据就是把里面的isdelete字段更新为1
func DeletePassengerFromDB(ctx *gin.Context) {
    passportse_no := ctx.PostForm("passportse_no")
    code := ctx.PostForm("code")

    token := getAccess(code)//根据前端传来的code获取token
    var customerid string
    customer_id,_,_,_,istoken := utils.GetUserByAccess(token,ctx)
    //先检查token是否有效
    if !istoken{
        fmt.Println("token无效")
        return
    }
    customerid = fmt.Sprintf("%v",customer_id)

    uid := GetMD5Hash(passportse_no+customerid)
    opend, db := OpenDB()
    if opend {
        fmt.Println("open success")

            stmt, err := db.Prepare("update passengers set isdelete=? where uid=?")
            CheckErr(err)
            res, err := stmt.Exec("1",uid)
            affect, err := res.RowsAffected()
            fmt.Println("删除数据：", affect)
            CheckErr(err)

             ctx.JSON(200, gin.H{
               "error_code": "0",
               "message": "乘客信息删除成功",
             })
        

    } else {
        fmt.Println("open faile:")
        ctx.JSON(200, gin.H{
        "error_code": "1",
        "message": "乘客信息删除异常",
        })
    }

}

func CheckErr(err error) {
    if err != nil {
        panic(err)
        fmt.Println("err:", err)
    }
}

// func GetTime() string {
//     const shortForm = "2006-01-02 15:04:05"
//     t := time.Now()
//     temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
//     str := temp.Format(shortForm)
//     fmt.Println(t)
//     return str
// }

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
    // opend, db := OpenDB()
    // if opend {
    //     fmt.Println("open success")
    //     // DeleteFromDB(db, 10)
    //     //QueryFromDB(db)
    //     //DeleteFromDB(db, 1)
    //     //UpdateDB(db, 5)
    //     //UpdateUID(db, 5)
    //     //UpdateTime(db, 4)
    //     insertToDB(db)
    // } else {
    //     fmt.Println("open faile:")
    // }

}