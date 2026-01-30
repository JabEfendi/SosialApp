package controllers

import (
	"net/http"
	"backend/db"
	"backend/helpers"
	"backend/models"
	"time"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
)

//ADMIN
func GetCorporates(c *gin.Context) {
	var corporates []models.Corporate

	if err := db.DB.Order("id desc").Find(&corporates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": corporates,
	})
}

func ChangeCorporateStatus(c *gin.Context) {
	adminID := c.MustGet("admin_id").(uint)
	id := c.Param("id")

	var corporate models.Corporate
	if err := db.DB.First(&corporate, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "corporate not found",
		})
		return
	}

	before := corporate

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	corporate.Status = input.Status
	corporate.UpdatedBy = &adminID

	if err := db.DB.Save(&corporate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	_ = helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    adminID,
		Action:     "UPDATE_STATUS",
		TargetType: "corporate",
		TargetID:   &corporate.ID,
		Before:     before,
		After:      corporate,
		Context:    c,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "status updated",
	})
}

func CreateCorporate(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=6"`
		Description string `json:"description"`
		Logo        string `json:"logo"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	adminIDValue, exists := c.Get("admin_id")
	if !exists {
		c.JSON(403, gin.H{"error": "OTP verification required"})
		return
	}

	adminID := adminIDValue.(uint)

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	referralCode, _ := helpers.GenerateUniqueCorporateReferral(db.DB)

	corporate := models.Corporate{
		Name:          input.Name,
		Email:         input.Email,
		Password:      string(hashed),
		Description:   input.Description,
		Logo:          input.Logo,
		Status:        "active",
		CreatedBy:     &adminID,
		Reffcorporate: referralCode,
	}

	db.DB.Create(&corporate)

	helpers.CreateAdminAuditLog(helpers.AuditPayload{
		AdminID:    adminID,
		Action:     "CREATE",
		TargetType: "corporate",
		TargetID:   &corporate.ID,
		After:      corporate,
		Context:    c,
	})

	c.JSON(201, gin.H{
		"message": "Corporate created by admin",
		"data":    corporate,
	})
}




//CORPORATE
type TempCorporate struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Email       string
	Password    string
	Description string
	Logo        string
	CreatedAt   time.Time
}

func CreateCorporateRequest(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=6"`
		Description string `json:"description"`
		Logo        string `json:"logo"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existing models.Corporate
	if err := db.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(400, gin.H{"error": "Email already registered"})
		return
	}

	db.DB.Where("email = ?", input.Email).Delete(&models.Corporate{})
	db.DB.Where("email = ?", input.Email).Delete(&models.OTPVerification{})

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	temp := models.Corporate{
		Name:        input.Name,
		Email:       input.Email,
		Password:    string(hashed),
		Description: input.Description,
		Logo:        input.Logo,
	}

	db.DB.Create(&temp)

	otp := helpers.GenerateOTP()
	db.DB.Create(&models.OTPVerification{
		Email:     input.Email,
		OTP:       otp,
		ExpiredAt: time.Now().Add(10 * time.Minute),
	})

	helpers.SendEmail(
		input.Email,
		"OTP Verifikasi Corporate",
		fmt.Sprintf("<h2>%s</h2><p>Berlaku 10 menit</p>", otp),
	)

	c.JSON(200, gin.H{
		"message": "OTP sent to corporate email",
	})
}

func VerifyCorporateOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var otpData models.OTPVerification
	if err := db.DB.
		Where("email = ? AND otp = ?", input.Email, input.OTP).
		First(&otpData).Error; err != nil {
		c.JSON(400, gin.H{"error": "Invalid OTP"})
		return
	}

	if time.Now().After(otpData.ExpiredAt) {
		c.JSON(400, gin.H{"error": "OTP expired"})
		return
	}

	var temp models.Corporate
	if err := db.DB.Where("email = ?", input.Email).First(&temp).Error; err != nil {
		c.JSON(404, gin.H{"error": "Data not found"})
		return
	}

	referralCode, _ := helpers.GenerateUniqueCorporateReferral(db.DB)

	corporate := models.Corporate{
		Name:          temp.Name,
		Email:         temp.Email,
		Password:      temp.Password,
		Description:   temp.Description,
		Logo:          temp.Logo,
		Reffcorporate: referralCode,
		Status:        "active",
	}

	db.DB.Create(&corporate)

	db.DB.Delete(&temp)
	db.DB.Delete(&otpData)

	c.JSON(200, gin.H{
		"message": "Corporate account created successfully",
		"data":    corporate,
	})
}

func CorporateDetail(c *gin.Context) {
	corporateID := c.MustGet("corporate_id").(uint)

	var corporate models.Corporate
	if err := db.DB.First(&corporate, corporateID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "corporate not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": corporate,
	})
}

func UpdateCorporate(c *gin.Context) {
	corporateID := c.MustGet("corporate_id").(uint)

	var corporate models.Corporate
	if err := db.DB.First(&corporate, corporateID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "corporate not found",
		})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Logo        string `json:"logo"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if input.Name != "" {
		corporate.Name = input.Name
	}
	if input.Description != "" {
		corporate.Description = input.Description
	}
	if input.Logo != "" {
		corporate.Logo = input.Logo
	}

	if err := db.DB.Save(&corporate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated",
		"data": gin.H{
			"id":          corporate.ID,
			"name":        corporate.Name,
			"description": corporate.Description,
			"logo":        corporate.Logo,
		},
	})
}

func CorporateLogin(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var corporate models.Corporate
	if err := db.DB.Where("email = ?", input.Email).First(&corporate).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "email or password is incorrect",
		})
		return
	}

	if corporate.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "account is not active",
		})
		return
	}

	if !helpers.CheckPasswordHash(input.Password, corporate.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "email or password is incorrect",
		})
		return
	}

	token, err := helpers.GenerateCorporateToken(corporate.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"data": gin.H{
			"token": token,
			"corporate": gin.H{
				"id":    corporate.ID,
				"name":  corporate.Name,
				"email": corporate.Email,
			},
		},
	})
}

func ForgotPasswordRequestCorporate(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var corporate models.Corporate
	if err := db.DB.Where("email = ?", input.Email).First(&corporate).Error; err != nil {
		c.JSON(404, gin.H{"error": "Email corporate tidak terdaftar"})
		return
	}

	// hapus OTP lama
	db.DB.Where("email = ?", input.Email).Delete(&models.OTPVerification{})

	otp := helpers.GenerateOTP()
	expired := time.Now().Add(10 * time.Minute)

	db.DB.Create(&models.OTPVerification{
		Email:     input.Email,
		OTP:       otp,
		ExpiredAt: expired,
	})

	emailBody := fmt.Sprintf(`
		<h3>Reset Password Corporate</h3>
		<p>OTP kamu adalah:</p>
		<h2>%s</h2>
		<p>Berlaku selama 10 menit</p>
	`, otp)

	helpers.SendEmail(input.Email, "Reset Password Corporate", emailBody)

	fmt.Println("FORGOT OTP CORPORATE:", otp)

	c.JSON(200, gin.H{
		"message": "OTP sent to corporate email",
	})
}

func ForgotPasswordVerifyCorporate(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var otpData models.OTPVerification
	if err := db.DB.
		Where("email = ? AND otp = ?", input.Email, input.OTP).
		First(&otpData).Error; err != nil {
		c.JSON(400, gin.H{"error": "OTP tidak valid"})
		return
	}

	if time.Now().After(otpData.ExpiredAt) {
		c.JSON(400, gin.H{"error": "OTP sudah expired"})
		return
	}

	c.JSON(200, gin.H{
		"message": "OTP verified",
	})
}

func ResetPasswordCorporate(c *gin.Context) {
	var input struct {
		Email       string `json:"email" binding:"required,email"`
		OTP         string `json:"otp" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var otpData models.OTPVerification
	if err := db.DB.
		Where("email = ? AND otp = ?", input.Email, input.OTP).
		First(&otpData).Error; err != nil {
		c.JSON(400, gin.H{"error": "OTP tidak valid"})
		return
	}

	if time.Now().After(otpData.ExpiredAt) {
		c.JSON(400, gin.H{"error": "OTP sudah expired"})
		return
	}

	var corporate models.Corporate
	if err := db.DB.Where("email = ?", input.Email).First(&corporate).Error; err != nil {
		c.JSON(404, gin.H{"error": "Corporate account not found"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	corporate.Password = string(hashed)
	db.DB.Save(&corporate)
	db.DB.Delete(&otpData)

	helpers.SendEmail(
		corporate.Email,
		"Password Corporate Berhasil Diubah",
		"<p>Password akun corporate kamu berhasil diubah.</p>",
	)

	c.JSON(200, gin.H{
		"message": "Corporate password reset successfully",
	})
}

func GetCorporateUsers(c *gin.Context) {
	corporateID, exists := c.Get("corporate_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var users []models.User

	if err := db.DB.
		Where("corporate_id = ?", corporateID).
		Preload("Corporate").
		Find(&users).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func CorporateApproveAgreement(c *gin.Context) {
    corporateID := c.MustGet("corporate_id").(uint)
    id := c.Param("id")

    var agreement models.Agreement
    db.DB.First(&agreement, id)

    if agreement.CorporateID != corporateID {
        c.JSON(403, gin.H{"error": "Unauthorized"})
        return
    }

    now := time.Now()

    db.DB.Model(&agreement).Updates(map[string]interface{}{
        "corporate_approved_at": now,
        "status": "corporate_sent",
    })

    helpers.CreateNotification(
        agreement.UserID,
        "Agreement Approved by Corporate",
        "Please review and approve the agreement",
        "agreement_sent",
    )

    c.JSON(200, gin.H{"message": "Agreement approved by corporate"})
}

func SendAgreementToUser(c *gin.Context) {
	corporateID := c.MustGet("corporate_id").(uint)
	id := c.Param("id")

	var input struct {
		StartDate time.Time  `json:"start_date" binding:"required"`
		EndDate   *time.Time `json:"end_date"`
		RevenueUserPercent float64 `json:"revenue_user_percent" binding:"required"`
		RevenueCorporatePercent float64 `json:"revenue_corporate_percent" binding:"required"`
		PaymentPeriod string `json:"payment_period" binding:"required"`
		ScopeDescription string `json:"scope_description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	var agreement models.Agreement
	if err := db.DB.First(&agreement, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Agreement not found"})
		return
	}

	if agreement.CorporateID != corporateID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	if agreement.Status != "requested" && agreement.Status != "revision_requested" {
		c.JSON(400, gin.H{"error": "Agreement cannot be sent in current status"})
		return
	}

	now := time.Now()

	if err := db.DB.Model(&agreement).Updates(map[string]interface{}{
		"start_date": input.StartDate,
		"end_date": input.EndDate,
		"revenue_user_percent": input.RevenueUserPercent,
		"revenue_corporate_percent": input.RevenueCorporatePercent,
		"payment_period": input.PaymentPeriod,
		"scope_description": input.ScopeDescription,
		"status": "corporate_sent",
		"corporate_approved_at": now,
		"user_approved_at": nil,
		"updated_at": now,
	}).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to send agreement"})
		return
	}

	helpers.CreateNotification(
		agreement.UserID,
		"New Agreement",
		"Corporate has sent you an agreement for review",
		"agreement_sent",
	)

	c.JSON(200, gin.H{"message": "Agreement sent to user"})
}

func CorporateRequestTermination(c *gin.Context) {
	corporateID := c.MustGet("corporate_id").(uint)
	id := c.Param("id")

	var input struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Reason is required"})
		return
	}

	var agreement models.Agreement
	if err := db.DB.First(&agreement, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Agreement not found"})
		return
	}

	if agreement.CorporateID != corporateID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	if agreement.Status != "active" {
		c.JSON(400, gin.H{"error": "Agreement is not active"})
		return
	}

	db.DB.Model(&agreement).Updates(map[string]interface{}{
		"status": "termination_requested_by_corporate",
		"termination_reason": input.Reason,
		"termination_requested_at": time.Now(),
	})

	helpers.CreateNotification(
		agreement.UserID,
		"Termination Requested",
		"Corporate requested to terminate the agreement",
		"agreement_termination_request",
	)

	c.JSON(200, gin.H{"message": "Termination request submitted"})
}

func CorporateApproveTermination(c *gin.Context) {
	corporateID := c.MustGet("corporate_id").(uint)
	id := c.Param("id")

	var agreement models.Agreement
	if err := db.DB.First(&agreement, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Agreement not found"})
		return
	}

	if agreement.CorporateID != corporateID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	if agreement.Status != "termination_requested_by_user" {
		c.JSON(400, gin.H{"error": "No termination request from user"})
		return
	}

	now := time.Now()
	tx := db.DB.Begin()

	if err := tx.Model(&agreement).Updates(map[string]interface{}{
		"status": "terminated",
		"terminated_at": now,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to terminate agreement"})
		return
	}

	// lepas user dari corporate
	tx.Model(&models.User{}).
		Where("id = ?", agreement.UserID).
		Updates(map[string]interface{}{
			"corporate_id": nil,
			"joined_corporate_at": nil,
		})

	tx.Commit()

	helpers.CreateNotification(
		agreement.UserID,
		"Agreement Terminated",
		"Corporate approved termination request",
		"agreement_terminated",
	)

	c.JSON(200, gin.H{"message": "Agreement terminated"})
}



//USER
func JoinCorporate(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if user.CorporateID != nil {
		c.JSON(403, gin.H{"error": "You are already joined to a corporate"})
		return
	}

	var input struct {
		ReffCorporate string `json:"reffcorporate" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	// cek agreement existing
	var count int64
	db.DB.Model(&models.Agreement{}).
		Where("user_id = ? AND status IN ?", userID,
			[]string{"requested", "revision_requested", "corporate_sent", "active"}).
		Count(&count)

	if count > 0 {
		c.JSON(409, gin.H{"error": "You already have an ongoing agreement"})
		return
	}

	var corporate models.Corporate
	if err := db.DB.Where("reffcorporate = ?", input.ReffCorporate).
		First(&corporate).Error; err != nil {
		c.JSON(404, gin.H{"error": "Corporate referral code not found"})
		return
	}

	agreement := models.Agreement{
		UserID:      userID,
		CorporateID: corporate.ID,
		AgreementNumber: fmt.Sprintf(
			"AGR-%d-%d-%d",
			corporate.ID, userID, time.Now().Unix(),
		),
		Status:        "requested",
		RevenueType:   "percentage",
		PaymentPeriod:"monthly",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.DB.Create(&agreement).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create agreement"})
		return
	}

	helpers.CreateNotification(
		corporate.ID,
		"New Partnership Request",
		"A user requested to join your corporate",
		"agreement_request",
	)

	c.JSON(200, gin.H{"message": "Partnership request submitted"})
}

func RequestAgreementRevision(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id := c.Param("id")

	var input struct {
		Note string `json:"note" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Note is required"})
		return
	}

	var agreement models.Agreement
	if err := db.DB.First(&agreement, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Agreement not found"})
		return
	}

	if agreement.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	if agreement.Status != "corporate_sent" {
		c.JSON(400, gin.H{"error": "Agreement is not revisable"})
		return
	}

	db.DB.Create(&models.AgreementRevision{
		AgreementID: agreement.ID,
		RequestedBy: "user",
		Note: input.Note,
		CreatedAt: time.Now(),
	})

	db.DB.Model(&agreement).Update("status", "revision_requested")

	helpers.CreateNotification(
		agreement.CorporateID,
		"Agreement Revision Requested",
		"User requested changes to the agreement",
		"agreement_revision",
	)

	c.JSON(200, gin.H{"message": "Revision request sent"})
}

func UserApproveAgreement(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id := c.Param("id")

	var agreement models.Agreement
	if err := db.DB.First(&agreement, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Agreement not found"})
		return
	}

	if agreement.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	if agreement.Status != "corporate_sent" {
		c.JSON(400, gin.H{"error": "Agreement is not ready for approval"})
		return
	}

	now := time.Now()
	tx := db.DB.Begin()

	if err := tx.Model(&agreement).Updates(map[string]interface{}{
		"user_approved_at": now,
		"status": "active",
		"updated_at": now,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to approve agreement"})
		return
	}

	if err := tx.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"corporate_id": agreement.CorporateID,
			"joined_corporate_at": now,
		}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	tx.Commit()

	helpers.CreateNotification(
		agreement.CorporateID,
		"Agreement Activated",
		"User has approved the agreement",
		"agreement_active",
	)

	c.JSON(200, gin.H{"message": "Agreement activated"})
}

func UserRequestTermination(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id := c.Param("id")

	var input struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Reason is required"})
		return
	}

	var agreement models.Agreement
	if err := db.DB.First(&agreement, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Agreement not found"})
		return
	}

	if agreement.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	if agreement.Status != "active" {
		c.JSON(400, gin.H{"error": "Agreement is not active"})
		return
	}

	db.DB.Model(&agreement).Updates(map[string]interface{}{
		"status": "termination_requested_by_user",
		"termination_reason": input.Reason,
		"termination_requested_at": time.Now(),
	})

	helpers.CreateNotification(
		agreement.CorporateID,
		"Termination Requested",
		"User requested to terminate the agreement",
		"agreement_termination_request",
	)

	c.JSON(200, gin.H{"message": "Termination request submitted"})
}

func UserApproveTermination(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id := c.Param("id")

	var agreement models.Agreement
	if err := db.DB.First(&agreement, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Agreement not found"})
		return
	}

	if agreement.UserID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}

	if agreement.Status != "termination_requested_by_corporate" {
		c.JSON(400, gin.H{"error": "No termination request from corporate"})
		return
	}

	now := time.Now()
	tx := db.DB.Begin()

	tx.Model(&agreement).Updates(map[string]interface{}{
		"status": "terminated",
		"terminated_at": now,
	})

	tx.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"corporate_id": nil,
			"joined_corporate_at": nil,
		})

	tx.Commit()

	helpers.CreateNotification(
		agreement.CorporateID,
		"Agreement Terminated",
		"User approved termination request",
		"agreement_terminated",
	)

	c.JSON(200, gin.H{"message": "Agreement terminated"})
}