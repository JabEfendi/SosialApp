package routes

import (
    "github.com/gin-gonic/gin"
    "backend/controllers"
    "backend/middlewares"
)

func SystemRoutes(r *gin.Engine) {
    r.GET("/testlog", controllers.Ceklog)
    r.POST("/notif/test", controllers.SendTestNotification)
    r.GET("/account-banks", controllers.GetAccountBanks)

    reg := r.Group("/register")
        reg.POST("/request-otp", controllers.RegisterRequest)
        reg.POST("/verify-otp", controllers.RegisterVerify)
        reg.POST("/resend-otp", controllers.RegisterResend)

    auth := r.Group("/auth")
        auth.POST("/login", controllers.Login)
        auth.POST("/save-fcm-token", controllers.SaveFCMToken)
        auth.POST("/google", controllers.GoogleLogin)
        auth.POST("/facebook", controllers.FacebookLogin)
        auth.POST("/forgot-password", controllers.ForgotPasswordRequest)
        auth.POST("/forgot-password/verify", controllers.ForgotPasswordVerify)
        auth.POST("/reset-password", controllers.ResetPassword)
        
    user := r.Group("/user").Use(middlewares.AuthMiddleware())
        user.PUT("/changepassword", controllers.ChangePass)
        user.PUT("/updateuser", controllers.UpdateUser)
        user.GET("/bio", controllers.GetUserDetail)
        user.POST("/photos/upload", controllers.UploadAvatar)
        user.GET("/photos", controllers.GetUserPhotos)
        user.POST("/photos/set-avatar", controllers.SetProfilePhoto)
        user.DELETE("/photos", controllers.DeleteUserPhoto)
        user.POST("/withdraw-account", controllers.CreateWithdrawAccount)
        user.PUT("/withdraw-account/request-update", controllers.RequestUpdateWithdrawAccount)
        user.GET("/withdraw-account", controllers.GetMyWithdrawAccount)
        user.POST("/withdraw", controllers.WithdrawCommission)
            
    kyc := r.Group("/kyc")
        kyc.POST("/", middlewares.AuthMiddleware(), controllers.SubmitOrUpdateKyc)
        kyc.POST("/approve", middlewares.AdminAuth(), controllers.ApproveKyc)
        kyc.POST("/reject", middlewares.AdminAuth(), controllers.RejectKyc)
        
    room := r.Group("/room").Use(middlewares.AuthMiddleware())
        room.POST("/", controllers.CreateRoom)
        room.POST("/join", controllers.JoinRoom)
        room.POST("/paidroom/join", controllers.JoinPaidRoom)
        room.GET("/:id/participants", controllers.GetRoomParticipants)
        room.GET("/:id", controllers.DetailRoom)
        room.GET("/list", controllers.ListRoom)
        room.PUT("/update/:id", controllers.UpdateRoom)

    roomchat := r.Group("/chatroom").Use(middlewares.AuthMiddleware())
        roomchat.POST("/send", controllers.SendMessage)
        roomchat.GET("/:roomID/messages", controllers.GetMessages)
        roomchat.GET("/:roomID/stream", controllers.GetRealtimeStream)

    directchat := r.Group("/directchat").Use(middlewares.AuthMiddleware())
        directchat.POST("/send", controllers.SendDirectMessage)
        directchat.GET("/:threadID/messages", controllers.GetDirectMessages)
        //note remember ini untuk update status pesan anjay tapi ini dari receivernya
        directchat.PUT("/:threadID/delivered/:userID", controllers.MarkDirectDelivered)

    topup := r.Group("/topup").Use(middlewares.AuthMiddleware())
        topup.POST("/", controllers.CreateTopUp)
        topup.POST("/callback", controllers.PaymentCallback)
        topup.GET("/history", controllers.GetTopUpHistory)

    community := r.Group("/community").Use(middlewares.AuthMiddleware())
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
        admin.Use(middlewares.AdminAuth())
        {
            admin.GET("/me", controllers.Me)
            admin.POST("/approve/:id", middlewares.SuperAdminOnly(), controllers.ApproveAdmin)
            admin.POST("/reject/:id", middlewares.SuperAdminOnly(), controllers.RejectAdmin)
            // admin.PUT("/:id", middlewares.SuperAdminOnly(), controllers.UpdateAdmin)
            admin.PUT("/profile", controllers.UpdateMyProfile)
            admin.POST("/:id/reset-password", middlewares.SuperAdminOnly(), controllers.ResetAdminPassword)
            admin.GET("/packages", controllers.GetPackages)
            admin.POST("/packages", middlewares.SuperAdminOnly(), controllers.CreatePackage)
            admin.PUT("/packages/:id", middlewares.SuperAdminOnly(), controllers.UpdatePackage)
            admin.PATCH("/packages/:id/status", middlewares.SuperAdminOnly(), controllers.ChangePackageStatus)
            admin.GET("/notifications", controllers.GetMyNotifications)
            admin.GET("/notifications/unread-count", controllers.GetUnreadNotificationCount)
            admin.PATCH("/notifications/:id/read", controllers.MarkNotificationAsRead)
            admin.PATCH("/notifications/read-all", controllers.MarkAllNotificationsAsRead)
            admin.GET("/corporates", controllers.GetCorporates)
            admin.POST("/corporates", middlewares.SuperAdminOnly(), controllers.CreateCorporate)
            admin.PUT("/corporates/:id", middlewares.SuperAdminOnly(), controllers.UpdateCorporate)
            admin.PATCH("/corporates/:id/status", middlewares.SuperAdminOnly(),  controllers.ChangeCorporateStatus)
            admin.GET("/reports", controllers.AdminListUserReports)
            admin.POST("/reports/:id/verify", controllers.AdminVerifyReport)
            admin.POST("/users/:id/status", controllers.AdminSuspendUser)
            admin.GET("/biostar", middlewares.SuperAdminOnly(), controllers.GetAllUsers) // karna namanya beda ini untuk get all user tapi khusus untuk superadmin
            admin.GET("/allstar", controllers.GetPublicUsers) // karna namanya beda ini untuk get all user tapi khusus untuk admin
            admin.GET("/all-community", controllers.GetAllCommunity)
        }

    sys := r.Group("/system").Use(middlewares.AdminAuth())
        sys.POST("/permissions", middlewares.RequirePermission("system.permission.create"), controllers.CreatePermission)
        sys.PUT("/roles/:role_id/permissions", middlewares.RequirePermission("system.role.permission.update"), controllers.UpdateRolePermissions)
        sys.GET("/notifications", middlewares.RequirePermission("system.notification.view"), controllers.GetNotificationSettings)
        sys.PUT("/notifications/:id", middlewares.RequirePermission("system.notification.update"), controllers.UpdateNotificationSetting)
        sys.POST("/legal", middlewares.RequirePermission("system.legal.update"), controllers.CreateLegalDocument)
        sys.POST("/email-campaigns", middlewares.RequirePermission("system.email.schedule"), controllers.CreateEmailCampaign)
        sys.POST("/maintenance", middlewares.RequirePermission("system.maintenance.create"), controllers.CreateMaintenance)

    reports := r.Group("/reports").Use(middlewares.AuthMiddleware())
        reports.GET("/reasons", controllers.GetReportReasons)
        reports.POST("/", controllers.ReportUser)
}