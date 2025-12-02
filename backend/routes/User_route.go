package routes

import (
    "github.com/gin-gonic/gin"
    "backend/controllers"
)

func UserRoutes(r *gin.Engine) {
    r.POST("/register", controllers.Register)
    r.POST("/login", controllers.Login)
    r.GET("/test", controllers.Test)
}
