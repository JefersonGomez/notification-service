package routes

import (
	"notification-service/controller"
	"notification-service/middlewares"

	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine) {

	r.POST("/registro", controller.Registro)
	r.POST("/login", controller.Login)

	p := r.Group("/")
	p.Use(middlewares.ValidarToken)
	{
		p.POST("/eventos", controller.PublicarEvento)
		p.GET("/eventos", controller.ObtenerEventos)
		p.PATCH("/eventos/:id", controller.MarcarLeido)
		p.GET("/ws", controller.WsHandler)

	}

}
