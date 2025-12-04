package routes

import (
    "github.com/gin-gonic/gin"
    "backend/controllers"
)

func UserRoutes(r *gin.Engine) {
    user := r.Group("/user")

    user.POST("/login", controllers.Login)

    reg := user.Group("/register")
    reg.POST("/request-otp", controllers.RegisterRequest)
    reg.POST("/verify-otp", controllers.RegisterVerify)
    reg.POST("/resend-otp", controllers.RegisterResend)

    user.POST("/save-fcm-token", controllers.SaveFCMToken)
    user.PUT("/:id", controllers.UpdateUser)
}


