package controllers

import (
	"net/http"
	"time"

	"backend/db"
	"backend/helpers"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func RegisterAdmin(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		RoleID   uint   `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := helpers.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	admin := models.Admin{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		RoleID:   2,
		Status:   "pending",
	}

	if err := db.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var creatorID uint
	if v, exists := c.Get("admin_id"); exists {
		creatorID = v.(uint)
	}
	_ = helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    creatorID,
		Action:     "CREATE",
		TargetType: "admin",
		TargetID:   &admin.ID,
		After:      admin,
		Context:    c,
	})

	var superAdmins []models.Admin
	db.DB.
		Where("role_id = ?", 1).
		Where("status = ?", "active").
		Find(&superAdmins)

	for _, sa := range superAdmins {
		_ = helpers.CreateAdminNotification(helpers.AdminNotificationPayload{
			AdminID: sa.ID,
			Type:    "account",
			Title:   "New Admin Registration",
			Message: "A new admin has registered and is awaiting approval.",
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "admin registered successfully",
		"admin":   admin.ID,
	})
}

func LoginAdmin(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.Admin
	if err := db.DB.Preload("Role").
		Where("email = ?", input.Email).
		First(&admin).Error; err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if admin.Status != "active" {
		message := "account is not active"

		switch admin.Status {
		case "pending":
			message = "account is pending approval"
		case "paused":
			message = "account is temporarily suspended"
		case "inactive":
			message = "account is inactive"
		case "rejected":
			message = "account registration was rejected"
		}

		c.JSON(http.StatusForbidden, gin.H{"error": message})
		return
	}

	if !helpers.CheckPasswordHash(input.Password, admin.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	now := time.Now()
	db.DB.Model(&admin).Update("last_login_at", &now)

	token, err := helpers.GenerateAdminToken(admin.ID, admin.Role.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"admin": gin.H{
			"id":    admin.ID,
			"name":  admin.Name,
			"email": admin.Email,
			"role":  admin.Role.Name,
		},
	})
}

func Me(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var admin models.Admin
	if err := db.DB.Preload("Role").
		First(&admin, adminID).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    admin.ID,
		"name":  admin.Name,
		"email": admin.Email,
		"role":  admin.Role.Name,
	})
}

func ApproveAdmin(c *gin.Context) {
	adminID := c.Param("id")
	superAdminID, _ := c.Get("admin_id")

	var admin models.Admin
	if err := db.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	if admin.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "admin already processed"})
		return
	}

	before := admin
	now := time.Now()
	admin.Status = "active"
	admin.ApprovedAt = &now
	admin.ApprovedBy = helpers.UintPtr(superAdminID.(uint))

	db.DB.Save(&admin)

	_ = helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    superAdminID.(uint),
		Action:     "UPDATE",
		TargetType: "admin",
		TargetID:   &admin.ID,
		Before:     before,
		After:      admin,
		Context:    c,
	})
	
	_ = helpers.CreateAdminNotification(helpers.AdminNotificationPayload{
		AdminID: admin.ID,
		Type:    "account",
		Title:   "Account Approved",
		Message: "Your admin account has been approved and is now active.",
	})

	c.JSON(http.StatusOK, gin.H{"message": "admin approved"})
}

func RejectAdmin(c *gin.Context) {
	adminID := c.Param("id")
	rejectorID := c.MustGet("admin_id").(uint)

	var admin models.Admin
	if err := db.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	if admin.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "admin already processed"})
		return
	}

	before := admin

	db.DB.Model(&admin).Update("status", "rejected")
	admin.Status = "rejected"

	adminIDUint := admin.ID

	_ = helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    rejectorID,
		Action:     "UPDATE",
		TargetType: "admin",
		TargetID:   &adminIDUint,
		Before:     before,
		After:      admin,
		Context:    c,
	})

	_ = helpers.CreateAdminNotification(helpers.AdminNotificationPayload{
		AdminID: admin.ID,
		Type:    "account",
		Title:   "Account Rejected",
		Message: "Your admin account registration has been rejected.",
	})

	c.JSON(http.StatusOK, gin.H{"message": "admin rejected"})
}


func ResetAdminPassword(c *gin.Context) {
	adminID := c.Param("id")
	superAdminID := c.MustGet("admin_id").(uint)

	var input struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.Admin
	if err := db.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	before := admin

	hashedPassword, err := helpers.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	admin.Password = hashedPassword
	db.DB.Save(&admin)

	_ = helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    superAdminID,
		Action:     "RESET_PASSWORD",
		TargetType: "admin",
		TargetID:   &admin.ID,
		Before:     before,
		After:      admin,
		Context:    c,
	})

	_ = helpers.CreateAdminNotification(helpers.AdminNotificationPayload{
		AdminID: admin.ID,
		Type:    "security",
		Title:   "Password Reset",
		Message: "Your password has been reset by a superadmin.",
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "password reset successfully",
	})
}


func UpdateMyProfile(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email" binding:"omitempty,email"`
		Password string `json:"password" binding:"omitempty,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var admin models.Admin
	if err := db.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "admin not found",
		})
		return
	}

	// simpan BEFORE untuk audit
	before := admin

	// update field yang dikirim saja
	if input.Name != "" {
		admin.Name = input.Name
	}

	if input.Email != "" {
		admin.Email = input.Email
	}

	if input.Password != "" {
		hashed, err := helpers.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to hash password",
			})
			return
		}
		admin.Password = hashed
	}

	if err := db.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// AUDIT LOG
	_ = helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    adminID,
		Action:     "UPDATE_PROFILE",
		TargetType: "admin",
		TargetID:   &admin.ID,
		Before:     before,
		After:      admin,
		Context:    c,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
	})
}


func GetPackages(c *gin.Context) {
	var packages []models.AdminPackage

	db.DB.Order("id desc").Find(&packages)

	c.JSON(http.StatusOK, gin.H{
		"data": packages,
	})
}

func CreatePackage(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)

	var input struct {
		Name       string  `json:"name" binding:"required"`
		CoinAmount int     `json:"coin_amount" binding:"required"`
		Price      float64 `json:"price" binding:"required"`
		BonusCoin  int     `json:"bonus_coin"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pkg := models.AdminPackage{
		Name:       input.Name,
		CoinAmount: input.CoinAmount,
		Price:      input.Price,
		BonusCoin:  input.BonusCoin,
		Status:     "active",
		CreatedBy:  &adminID,
	}

	if err := db.DB.Create(&pkg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// AUDIT
	_ = helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    adminID,
		Action:     "CREATE",
		TargetType: "package",
		TargetID:   &pkg.ID,
		After:      pkg,
		Context:    c,
	})

	c.JSON(http.StatusCreated, gin.H{
		"message": "package created",
		"data":    pkg,
	})
}

func UpdatePackage(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)
	id := c.Param("id")

	var pkg models.AdminPackage
	if err := db.DB.First(&pkg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "package not found"})
		return
	}

	before := pkg

	var input struct {
		Name       string  `json:"name"`
		CoinAmount int     `json:"coin_amount"`
		Price      float64 `json:"price"`
		BonusCoin  int     `json:"bonus_coin"`
		Status     string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		pkg.Name = input.Name
	}
	if input.CoinAmount > 0 {
		pkg.CoinAmount = input.CoinAmount
	}
	if input.Price > 0 {
		pkg.Price = input.Price
	}
	if input.BonusCoin >= 0 {
		pkg.BonusCoin = input.BonusCoin
	}
	if input.Status != "" {
		pkg.Status = input.Status
	}

	pkg.UpdatedBy = &adminID

	db.DB.Save(&pkg)

	// AUDIT
	_ = helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    adminID,
		Action:     "UPDATE",
		TargetType: "package",
		TargetID:   &pkg.ID,
		Before:     before,
		After:      pkg,
		Context:    c,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "package updated",
	})
}

func ChangePackageStatus(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)
	id := c.Param("id")

	var pkg models.AdminPackage
	if err := db.DB.First(&pkg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "package not found"})
		return
	}

	before := pkg

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pkg.Status = input.Status
	pkg.UpdatedBy = &adminID
	db.DB.Save(&pkg)

	_ = helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    adminID,
		Action:     "UPDATE_STATUS",
		TargetType: "package",
		TargetID:   &pkg.ID,
		Before:     before,
		After:      pkg,
		Context:    c,
	})

	_ = helpers.CreateAdminNotification(helpers.AdminNotificationPayload{
		AdminID: *pkg.CreatedBy,
		Type:    "package",
		Title:   "Package Updated",
		Message: "A package you created has been updated or changed status.",
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "status updated",
	})
}


func GetMyNotifications(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)

	var notifs []models.AdminNotification

	db.DB.
		Where("admin_id = ?", adminID).
		Order("id desc").
		Find(&notifs)

	c.JSON(http.StatusOK, gin.H{
		"data": notifs,
	})
}

func GetUnreadNotificationCount(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)

	var count int64
	db.DB.
		Model(&models.AdminNotification{}).
		Where("admin_id = ? AND is_read = false", adminID).
		Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"unread": count,
	})
}

func MarkNotificationAsRead(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)
	id := c.Param("id")

	result := db.DB.
		Model(&models.AdminNotification{}).
		Where("id = ? AND admin_id = ?", id, adminID).
		Update("is_read", true)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "notification not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "notification marked as read",
	})
}

func MarkAllNotificationsAsRead(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)

	db.DB.
		Model(&models.AdminNotification{}).
		Where("admin_id = ? AND is_read = false", adminID).
		Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{
		"message": "all notifications marked as read",
	})
}
