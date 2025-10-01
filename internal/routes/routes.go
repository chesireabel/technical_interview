package routes

import (
	"github.com/chesireabel/Technical-Interview/internal/handlers"
	"github.com/chesireabel/Technical-Interview/internal/middleware"


	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, customerHandler *handlers.CustomerHandler, orderHandler *handlers.OrderHandler,oidc *middleware.OIDC,returnToURL string) {
	//Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Server",
			"status":  "Server is running",
		})
	})

	//Auth routes
		auth := r.Group("/auth")
	{
		auth.GET("/login", oidc.LoginHandler())
		auth.GET("/callback", oidc.CallbackHandler())
		auth.GET("/logout", oidc.LogoutHandler(returnToURL))
	}


	//Customers routes
	customers := r.Group("/customers")
	{
		customers.POST("", customerHandler.CreateCustomer)
		customers.GET("", customerHandler.GetAllCustomers)
		customers.GET("/:id", customerHandler.GetCustomer)
		customers.PUT("/:id", customerHandler.UpdateCustomer)
		customers.DELETE("/:id", customerHandler.DeleteCustomer)
	}

	//Orders routes
	orders := r.Group("/orders")
	{
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("", orderHandler.GetAllOrders)
		orders.GET("/:id", orderHandler.GetOrder)
		orders.PUT("/:id", orderHandler.UpdateOrder)
		orders.DELETE("/:id", orderHandler.DeleteOrder)
	}

	//Get orders made by customer
	r.GET("/customers/:id/orders", orderHandler.GetOrdersByCustomer)
}