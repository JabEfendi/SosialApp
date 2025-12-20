package routes

import (
    "github.com/gin-gonic/gin"
    "backend/controllers"
    "backend/middlewares"
)

func SystemRoutes(r *gin.Engine) {
    r.GET("/test", controllers.Test)
    r.GET("/testlog", controllers.Ceklog)
    r.POST("/notif/test", controllers.SendTestNotification)

    reg := r.Group("/register")
        reg.POST("/request-otp", controllers.RegisterRequest)
        reg.POST("/verify-otp", controllers.RegisterVerify)
        reg.POST("/resend-otp", controllers.RegisterResend)

    auth := r.Group("/auth")
        auth.POST("/login", controllers.Login)
        // auth.POST("/upload-avatar", controllers.UploadAvatar)
        auth.POST("/save-fcm-token", controllers.SaveFCMToken)
        auth.PUT("/:id", controllers.UpdateUser)
        auth.PUT("/changepassword/:id", controllers.ChangePass)
        auth.POST("/google", controllers.GoogleLogin)
        auth.POST("/facebook", controllers.FacebookLogin)
        auth.POST("/forgot-password", controllers.ForgotPasswordRequest)
        auth.POST("/forgot-password/verify", controllers.ForgotPasswordVerify)
        auth.POST("/reset-password", controllers.ResetPassword)
        
    user := r.Group("/user")
        user.GET("/:id", controllers.GetUserDetail)
        user.POST("/photos/upload", controllers.UploadAvatar)
        user.GET("/photos/:id", controllers.GetUserPhotos)
        user.POST("/photos/set-avatar", controllers.SetProfilePhoto)
        user.DELETE("/photos/:id", controllers.DeleteUserPhoto)
    
    kyc := r.Group("/kyc")
        kyc.POST("/", controllers.SubmitOrUpdateKyc)
        kyc.POST("/approve", controllers.ApproveKyc)
        kyc.POST("/reject", controllers.RejectKyc)
        
    room := r.Group("/room")
        room.POST("/", controllers.CreateRoom)
        room.POST("/join", controllers.JoinRoom)
        room.GET("/:id/participants", controllers.GetRoomParticipants)
        room.GET("/:id", controllers.DetailRoom)
        room.GET("/list", controllers.ListRoom)
        room.PUT("/update/:id", controllers.UpdateRoom)

    roomchat := r.Group("/chatroom")
		roomchat.POST("/send", controllers.SendMessage)
		roomchat.GET("/:roomID/messages", controllers.GetMessages)
		roomchat.GET("/:roomID/stream", controllers.GetRealtimeStream)

    directchat := r.Group("/directchat")
        directchat.POST("/send", controllers.SendDirectMessage)
        directchat.GET("/:threadID/messages", controllers.GetDirectMessages)
        //note remember ini untuk update status pesan anjay tapi ini dari receivernya
        directchat.PUT("/:threadID/delivered/:userID", controllers.MarkDirectDelivered)

    topup := r.Group("/topup")
        topup.Use(middlewares.AuthMiddleware())
        {
            topup.POST("", controllers.CreateTopUp)
            topup.POST("/callback", controllers.PaymentCallback)
            topup.GET("/history", controllers.GetTopUpHistory)
        }

    community := r.Group("/community")
        community.POST("/", controllers.CreateCommunity)
        community.POST("/join", controllers.JoinCommunity)
        community.GET("/:id", controllers.GetCommunityDetail)
        community.GET("/:id/members", controllers.GetCommunityMembers)
        community.POST("/chat/send", controllers.SendCommunityMessage)
        community.GET("/chat/:communityID", controllers.GetCommunityMessages)
        community.POST("/photos/upload", controllers.UploadCommunityAvatar)
        community.GET("/:id/photos", controllers.GetPhotosCommunity)
        community.POST("/photos/set-avatar", controllers.SetProfilePhotoCommunity)
        community.DELETE("/photos/:id", controllers.DeletePhotoCommunity)

}


