package routes

import (
    "github.com/gin-gonic/gin"
    "backend/controllers"
)

func SystemRoutes(r *gin.Engine) {
    r.GET("/test", controllers.Test)
    r.GET("/testlog", controllers.Ceklog)
    r.POST("/notif/test", controllers.SendTestNotification)
}

