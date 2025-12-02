package routes

import (
    "github.com/gin-gonic/gin"
    "backend/controllers"
)

func UserRoutes(r *gin.Engine) {
    r.POST("/register", controllers.Register)
    r.POST("/login", controllers.Login)

    r.POST("/auth/google", controllers.GoogleLogin)
    r.POST("/auth/facebook", controllers.FacebookLogin)

    r.GET("/test", controllers.Test)
}

