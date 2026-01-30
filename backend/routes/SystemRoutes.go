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
        user.POST("/agreements/join", controllers.JoinCorporate)
		user.POST("/agreements/:id/revision", controllers.RequestAgreementRevision)
		user.POST("/agreements/:id/approve", controllers.UserApproveAgreement)
        user.POST("/agreements/:id/terminate", controllers.UserRequestTermination)
        user.POST("/agreements/:id/terminate/approve", controllers.UserApproveTermination)
            
    kyc := r.Group("/kyc")
        kyc.POST("/", middlewares.AuthMiddleware(), controllers.SubmitOrUpdateKyc)
        kyc.POST("/approve", middlewares.AdminAuth(), controllers.ApproveKyc)
        kyc.POST("/reject", middlewares.AdminAuth(), controllers.RejectKyc)
        
    room := r.Group("/room")
        room.Use(middlewares.AuthMiddleware())
        {
            room.POST("/", controllers.CreateRoom)
            room.POST("/join", controllers.JoinRoom)
            room.POST("/paidroom/join", controllers.JoinPaidRoom)
            room.GET("/:id/participants", controllers.GetRoomParticipants)
            room.GET("/:id", controllers.DetailRoom)
            room.GET("/list", controllers.ListRoom)
            room.PUT("/update/:id", controllers.UpdateRoom)
        }
        room.POST("/revenue/release", middlewares.AdminAuth(), func(c *gin.Context){
            err := controllers.ReleaseRoomRevenue()
            if err != nil {
                c.JSON(500, gin.H{"error": err.Error()})
                return
            }
            c.JSON(200, gin.H{"message": "Revenue released successfully"})
        })
        room.POST("/corporate-approve-payout", middlewares.CorporateAuth(), controllers.CorporateApprovePayout)

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
            admin.GET("/reports", controllers.AdminListUserReports)
            admin.POST("/reports/:id/verify", controllers.AdminVerifyReport)
            admin.POST("/users/:id/status", controllers.AdminSuspendUser)
            admin.GET("/biostar", middlewares.SuperAdminOnly(), controllers.GetAllUsers) // karna namanya beda ini untuk get all user tapi khusus untuk superadmin
            admin.GET("/allstar", controllers.GetPublicUsers) // karna namanya beda ini untuk get all user tapi khusus untuk admin
            admin.GET("/all-community", controllers.GetAllCommunity)
            admin.GET("/dashboard", controllers.AdminDashboard)
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

    corporate := r.Group("/corporate")
        corporate.POST("/corpo-create", middlewares.SuperAdminOnly(), controllers.CreateCorporate)
        corporate.POST("/register", controllers.CreateCorporateRequest)
        corporate.POST("/verify-otp", controllers.VerifyCorporateOTP)
        corporate.POST("/login", controllers.CorporateLogin)
        corporate.GET("/detail", middlewares.CorporateAuth(), controllers.CorporateDetail)
        corporate.PUT("/profile-update", middlewares.CorporateAuth(), controllers.UpdateCorporate)
        corporate.GET("/users", middlewares.CorporateAuth(), controllers.GetCorporateUsers)
        corporate.PATCH("/corporates/:id/status", middlewares.SuperAdminOnly(),  controllers.ChangeCorporateStatus)
        corporate.POST("/agreements/:id/approve", middlewares.CorporateAuth(), controllers.CorporateApproveAgreement)
		corporate.POST("/agreements/:id/send", middlewares.CorporateAuth(), controllers.SendAgreementToUser)
        corporate.POST("/agreements/:id/terminate", controllers.CorporateRequestTermination)
        corporate.POST("/agreements/:id/terminate/approve", controllers.CorporateApproveTermination)
    
    dashboard := r.Group("/api/corporate/dashboard")
        dashboard.Use(middlewares.CorporateAuth())
        {
            dashboard.GET("/income", controllers.CorporateIncomeChart)
            dashboard.GET("/withdraw-request", controllers.CorporateWithdrawRequestChart)
            dashboard.GET("/user-withdraw", controllers.UserWithdrawAccumulationChart)

            // untuk hit seluruh chart
            dashboard.GET("", controllers.CorporateDashboard)
        }
}