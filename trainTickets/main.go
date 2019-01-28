package main

import (
	"trainTickets/api"

	_ "trainTickets/docs"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.New()

    // Logging to a file.
    // f, _ := os.Create("gin.log")
    // gin.DefaultWriter = io.MultiWriter(f)

	
    api.SetPush()
	r.POST("/trainTickets/cityCode", api.CityCode)
	r.POST("/trainTickets/ticketsAvailable", api.TicketsAvailable)
	r.POST("/trainTickets/submit", api.Submit)
	r.POST("/trainTickets/pay", api.Pay)
	r.POST("/trainTickets/orderStatus", api.OrderStatus)
	r.POST("/trainTickets/refund", api.Refund)
	r.POST("/trainTickets/orders", api.Orders)
	r.POST("/trainTickets/cancel", api.Cancel)
	r.POST("/trainTickets/insertPassenger", api.InsertPassengerToDB)
	r.POST("/trainTickets/queryPassenger", api.QueryPassengerFromDB)
	r.POST("/trainTickets/updatePassenger", api.UpdatePassengerToDB)
	r.POST("/trainTickets/deletePassenger", api.DeletePassengerFromDB)
	r.POST("/trainTickets/insertSearchRecordToDB", api.InsertSearchRecordToDB)
	r.POST("/trainTickets/querySearchRecordFromDB", api.QuerySearchRecordFromDB)
	r.POST("/trainTickets/updateSearchRecordToDB", api.UpdateSearchRecordToDB)
	r.POST("/trainTickets/deleteSearchRecordFromDB", api.DeleteSearchRecordFromDB)
	r.POST("/trainTickets/submit_callback", api.Submit_callback)
	r.POST("/trainTickets/pay_callback", api.Pay_callback)
	r.POST("/trainTickets/refund_callback", api.Refund_callback)


    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run()
}
