package routes

import (
    "github.com/gin-gonic/gin"
    "backend/controllers"
)

func SystemRoutes(r *gin.Engine) {
    r.GET("/test", controllers.Test)
    r.GET("/testlog", controllers.Ceklog)
    r.POST("/notif/test", controllers.SendTestNotification)

    r.POST("/kyc", controllers.SubmitOrUpdateKyc)
    r.POST("/room", controllers.CreateRoom)


    user := r.Group("/user")
        user.POST("/login", controllers.Login)
        user.POST("/upload-avatar", controllers.UploadAvatar)
        user.POST("/save-fcm-token", controllers.SaveFCMToken)
        user.PUT("/:id", controllers.UpdateUser)
        user.PUT("/changepassword/:id", controllers.ChangePass)

    reg := r.Group("/register")
        reg.POST("/request-otp", controllers.RegisterRequest)
        reg.POST("/verify-otp", controllers.RegisterVerify)
        reg.POST("/resend-otp", controllers.RegisterResend)

    auth := r.Group("/auth")
        auth.POST("/google", controllers.GoogleLogin)
        auth.POST("/facebook", controllers.FacebookLogin)
        
    room := r.Group("/room")
        room.POST("/join", controllers.JoinRoom)
        room.GET("/:id/participants", controllers.GetRoomParticipants)
        room.GET("/list", controllers.ListRoom)

    chat := r.Group("/chat")
	{
		chat.POST("/send", controllers.SendMessage)
		chat.GET("/:roomID/messages", controllers.GetMessages)
		chat.GET("/:roomID/stream", controllers.GetRealtimeStream)
	}

    direct := r.Group("/direct")
    {
        direct.POST("/send", controllers.SendDirectMessage)
        direct.GET("/:threadID/messages", controllers.GetDirectMessages)
        //note remember ini untuk update status pesan anjay tapi ini dari receivernya
        direct.PUT("/:threadID/delivered/:userID", controllers.MarkDirectDelivered)
    }

}


