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
        kyc.POST("/approve", middlewares.AdminAuth(), controllers.ApproveKyc)
        kyc.POST("/reject", middlewares.AdminAuth(), controllers.RejectKyc)
        
    room := r.Group("/room")
        room.POST("/", controllers.CreateRoom)
        room.POST("/join", controllers.JoinRoom)
        room.POST("/paidroom/join", controllers.JoinPaidRoom)
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

    admin := r.Group("/admin")
        admin.POST("/register", controllers.RegisterAdmin)
        admin.POST("/login", controllers.LoginAdmin)
        admin.GET("/me", controllers.Me)
        admin.POST("/approve/:id", middlewares.AdminAuth(), middlewares.SuperAdminOnly(), controllers.ApproveAdmin)
		admin.POST("/reject/:id", middlewares.AdminAuth(), middlewares.SuperAdminOnly(), controllers.RejectAdmin)
        // admin.PUT("/:id", middlewares.AdminAuth(), middlewares.SuperAdminOnly(), controllers.UpdateAdmin)
        admin.PUT("/profile", middlewares.AdminAuth(), controllers.UpdateMyProfile)
        admin.POST("/:id/reset-password", middlewares.AdminAuth(), middlewares.SuperAdminOnly(), controllers.ResetAdminPassword)
        admin.GET("/packages", middlewares.AdminAuth(), controllers.GetPackages)
        admin.POST("/packages", middlewares.AdminAuth(), middlewares.SuperAdminOnly(), controllers.CreatePackage)
        admin.PUT("/packages/:id", middlewares.AdminAuth(), middlewares.SuperAdminOnly(), controllers.UpdatePackage)
        admin.PATCH("/packages/:id/status", middlewares.AdminAuth(), middlewares.SuperAdminOnly(), controllers.ChangePackageStatus)
        admin.GET("/notifications", middlewares.AdminAuth(), controllers.GetMyNotifications)
        admin.GET("/notifications/unread-count", middlewares.AdminAuth(), controllers.GetUnreadNotificationCount)
        admin.PATCH("/notifications/:id/read", middlewares.AdminAuth(), controllers.MarkNotificationAsRead)
        admin.PATCH("/notifications/read-all", middlewares.AdminAuth(), controllers.MarkAllNotificationsAsRead)
        admin.GET("/corporates", middlewares.AdminAuth(), controllers.GetCorporates)
        admin.POST("/corporates", middlewares.AdminAuth(), middlewares.SuperAdminOnly(), controllers.CreateCorporate)
        admin.PUT("/corporates/:id", middlewares.AdminAuth(), middlewares.SuperAdminOnly(), controllers.UpdateCorporate)
        admin.PATCH("/corporates/:id/status", middlewares.AdminAuth(), middlewares.SuperAdminOnly(),  controllers.ChangeCorporateStatus)

    sys := admin.Group("/system")
        sys.POST("/permissions", middlewares.AdminAuth(), middlewares.RequirePermission("system.permission.create"), controllers.CreatePermission)
        sys.PUT("/roles/:role_id/permissions", middlewares.AdminAuth(), middlewares.RequirePermission("system.role.permission.update"), controllers.UpdateRolePermissions)
        sys.GET("/notifications", middlewares.AdminAuth(), middlewares.RequirePermission("system.notification.view"), controllers.GetNotificationSettings)
        sys.PUT("/notifications/:id", middlewares.AdminAuth(), middlewares.RequirePermission("system.notification.update"), controllers.UpdateNotificationSetting)
        sys.POST("/legal", middlewares.AdminAuth(), middlewares.RequirePermission("system.legal.update"), controllers.CreateLegalDocument)
        sys.POST("/email-campaigns", middlewares.AdminAuth(), middlewares.RequirePermission("system.email.schedule"), controllers.CreateEmailCampaign)
        sys.POST("/maintenance", middlewares.AdminAuth(), middlewares.RequirePermission("system.maintenance.create"), controllers.CreateMaintenance)
}


