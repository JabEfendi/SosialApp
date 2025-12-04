package routes

import (
    "github.com/gin-gonic/gin"
    "backend/controllers"
)

func UserRoutes(r *gin.Engine) {
    // r.POST("/register", controllers.Register)
    r.POST("/login", controllers.Login)

    r.POST("/register/request-otp", controllers.RegisterRequest)
    r.POST("/register/verify-otp", controllers.RegisterVerify)
    r.POST("/register/resend-otp", controllers.RegisterResend)

    r.POST("/notif/test", controllers.SendTestNotification)
    r.POST("/save-fcm-token", controllers.SaveFCMToken)
    r.PUT("/user/:id", controllers.UpdateUser)

    r.POST("/auth/google", controllers.GoogleLogin)
    r.POST("/auth/facebook", controllers.FacebookLogin)

    r.GET("/test", controllers.Test)
    r.GET("/testlog", controllers.Ceklog)
}

