package main

import (
	"trainTickets/api"

	_ "trainTickets/docs"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
    api.SetPush()
	r.POST("/trainTickets/cityCode", api.CityCode)
	r.POST("/trainTickets/ticketsAvailable", api.TicketsAvailable)
	r.POST("/trainTickets/submit", api.Submit)
	r.POST("/trainTickets/pay", api.Pay)
	r.POST("/trainTickets/orderStatus", api.OrderStatus)
	r.POST("/trainTickets/refund", api.Refund)
	r.POST("/trainTickets/orders", api.Orders)
	r.POST("/trainTickets/cancel", api.Cancel)
	r.POST("/trainTickets/submit_callback", api.Submit_callback)
	r.POST("/trainTickets/pay_callback", api.Pay_callback)
	r.POST("/trainTickets/refund_callback", api.Refund_callback)

	r.Run()
}
