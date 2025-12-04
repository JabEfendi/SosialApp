package routes

import (
    "github.com/gin-gonic/gin"
    "backend/controllers"
)

func AuthRoutes(r *gin.Engine) {
    auth := r.Group("/auth")

    auth.POST("/google", controllers.GoogleLogin)
    auth.POST("/facebook", controllers.FacebookLogin)
}
