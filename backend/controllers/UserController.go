package controllers

import (
	"backend/db"
	"backend/models"
	"backend/helpers"
	"net/http"
  	"bytes"
	"net/url"
	"encoding/json"
	"io"
	"strings"
	"time"
	"fmt"
    "os"
    "path/filepath"

    "github.com/google/uuid"
    "github.com/spf13/cast"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register 
type RegisterInput struct {
    Name      string `json:"name" binding:"required"`
    Username  string `json:"username" binding:"required"`
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=6"`
    Gender    string `json:"gender"`
    Phone     string `json:"phone"`
    Bio       string `json:"bio"`
    Country   string `json:"country"`
    Address   string `json:"address"`
    ReferralCode string `json:"referral_code"`
}

func RegisterRequest(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var existing models.User
    if err := db.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
        c.JSON(400, gin.H{"error": "Email is already in use"})
        return
    }

    db.DB.Where("email = ?", input.Email).Delete(&models.TempUser{})
    db.DB.Where("email = ?", input.Email).Delete(&models.OTPVerification{})

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

    var emptyDate *time.Time = nil
    temp := models.TempUser{
        Email:    input.Email,
        Name:     input.Name,
        Username: input.Username,
        Password: string(hashedPassword),
        Gender:   input.Gender,
        Birthdate: emptyDate,
        Phone:    input.Phone,
        Bio:      input.Bio,
        Country:  input.Country,
        Address:  input.Address,
        ReferralCode: strings.TrimSpace(input.ReferralCode),
    }

    db.DB.Create(&temp)

    otp := helpers.GenerateOTP()
    expired := time.Now().Add(10 * time.Minute)

    db.DB.Create(&models.OTPVerification{
        Email:     input.Email,
        OTP:       otp,
        ExpiredAt: expired,
    })

    fmt.Println("OTP DIKIRIM:", otp)

    c.JSON(200, gin.H{
        "message": "OTP has been sent to email",
        "email":   input.Email,
    })
}

func generateUniqueReferralCode(name string) string {
	for {
		code := helpers.GenerateReferralCode(name)

		var count int64
		db.DB.Model(&models.User{}).
			Where("referral_code = ?", code).
			Count(&count)

		if count == 0 {
			return code
		}
	}
}

func RegisterVerify(c *gin.Context) {
    var input struct {
        Email string `json:"email" binding:"required,email"`
        OTP   string `json:"otp" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var otpData models.OTPVerification
    if err := db.DB.Where("email = ? AND otp = ?", input.Email, input.OTP).
        First(&otpData).Error; err != nil {
        c.JSON(400, gin.H{"error": "Wrong OTP"})
        return
    }

    if time.Now().After(otpData.ExpiredAt) {
        c.JSON(400, gin.H{"error": "OTP expired, please resend"})
        return
    }

    var temp models.TempUser
    if err := db.DB.Where("email = ?", input.Email).First(&temp).Error; err != nil {
        c.JSON(400, gin.H{"error": "Data not found"})
        return
    }

    referralCode := generateUniqueReferralCode(temp.Name)
    user := models.User{
        Name:     temp.Name,
        Username: temp.Username,
        Email:    temp.Email,
        Password: temp.Password,
        Gender:   temp.Gender,
        Birthdate: temp.Birthdate,
        Phone:    temp.Phone,
        Bio:      temp.Bio,
        Country:  temp.Country,
        Address:  temp.Address,
        Provider: "local",
        ReferralCode: referralCode,
    }

    db.DB.Create(&user)

    if temp.ReferralCode != "" {
        var referrer models.User

        err := db.DB.
            Where("referral_code = ?", temp.ReferralCode).
            First(&referrer).Error

        if err == nil && referrer.ID != user.ID {
            user.ReferredBy = &referrer.ID
            db.DB.Save(&user)

            referral := models.Referral{
                ReferrerID: referrer.ID,
                ReferredID: user.ID,
                Status:     "pending",
            }
            db.DB.Create(&referral)
        }
    }

    db.DB.Delete(&temp)
    db.DB.Delete(&otpData)

    emailBody := fmt.Sprintf(`
        <h2>Congratulations %s 🎉</h2>
        <p>Your account has been successfully created on <b>%s</b>.</p>
    `, user.Name, time.Now().Format("02 Jan 2006"))

    err := helpers.SendEmail(user.Email, "🎉 Account Created Successfully", emailBody)
    if err != nil {
        fmt.Println("Email send error:", err)
    }

    var token models.UserFCMToken
    db.DB.Where("user_id = ?", user.ID).Last(&token)

    message := fmt.Sprintf("Halo %s! your account has been successfully created on %s",
            user.Name,
            time.Now().Format("02 Jan 2006"),
    )

    notif := models.Notification{
            UserID: user.ID,
            Title:  "Account Created Successfully 🎉",
            Message: message,
    }
    db.DB.Create(&notif)

    if token.FCMToken != "" {
            helpers.SendFCMToken(token.FCMToken,
                    "Account Created Successfully 🎉",
                    message,
            )
    }

    c.JSON(200, gin.H{
        "message": "Account Created Successfully",
        "user":    user,
    })
}

func RegisterResend(c *gin.Context) {
    var input struct {
        Email string `json:"email" binding:"required,email"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    db.DB.Where("email = ?", input.Email).Delete(&models.OTPVerification{})

    otp := helpers.GenerateOTP()

    db.DB.Create(&models.OTPVerification{
        Email:     input.Email,
        OTP:       otp,
        ExpiredAt: time.Now().Add(10 * time.Minute),
    })

    fmt.Println("OTP RESEND:", otp)

    c.JSON(200, gin.H{"message": "New OTP sent"})
}


// _________________________________________________________________________________________________
// Notif
func SendTestNotification(c *gin.Context) {
    var input struct {
        UserID uint   `json:"user_id" binding:"required"`
        Title  string `json:"title" binding:"required"`
        Body   string `json:"body" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var token models.UserFCMToken
    if err := db.DB.Where("user_id = ?", input.UserID).
        Order("id desc").
        First(&token).Error; err != nil {
        c.JSON(400, gin.H{"error": "FCM token user not found"})
        return
    }

    notif := models.Notification{
        UserID: input.UserID,
        Title:  input.Title,
        Message: input.Body,
    }
    db.DB.Create(&notif)

    if token.FCMToken != "" {
        helpers.SendFCMToken(token.FCMToken, input.Title, input.Body)
    }

    c.JSON(200, gin.H{
        "message": "Notification sent",
        "sent_to": token.FCMToken,
    })
}


// _________________________________________________________________________________________________
// token fcm untuk notif
type SaveFCMTokenInput struct {
    UserID   uint   `json:"user_id" binding:"required"`
    FCMToken string `json:"fcm_token" binding:"required"`
    Device   string `json:"device" binding:"required"`
}

func SaveFCMToken(c *gin.Context) {

    fmt.Println("➡️ SaveFCMToken() DIPANGGIL")

    var input SaveFCMTokenInput

    body, _ := io.ReadAll(c.Request.Body)
    fmt.Println("📥 Raw Body:", string(body))
    c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

    if err := c.ShouldBindJSON(&input); err != nil {
        fmt.Println("❌ ERROR BIND JSON:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    fmt.Println("✔️ JSON input terbind:", input)

    data := models.UserFCMToken{
        UserID:   input.UserID,
        FCMToken: input.FCMToken,
        Device:   input.Device,
    }

    fmt.Println("🟦 Data sebelum simpan:", data)

    err := db.DB.
        Where("user_id = ? AND device = ?", input.UserID, input.Device).
        Assign(data).
        FirstOrCreate(&data).Error

    if err != nil {
        fmt.Println("❌ ERROR QUERY DB:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save FCM token"})
        return
    }

    fmt.Println("✅ SUCCESS — Data tersimpan:", data)

    c.JSON(http.StatusOK, gin.H{
        "message": "FCM token saved successfully",
        "data":    data,
    })
}


// _________________________________________________________________________________________________
// Login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
    var input LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    db.DB.Where("email = ? AND provider = 'local'", input.Email).First(&user)

    if user.ID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
        return
    }

    createLoginLog(c, user.ID)

    c.JSON(http.StatusOK, gin.H{"message": "Login success", "user": user})
}


// _________________________________________________________________________________________________
// change password
func ChangePass(c *gin.Context) {
    userID := c.Param("id")

    var input struct {
        OldPassword string `json:"old_password"`
        NewPassword string `json:"new_password"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
        c.JSON(400, gin.H{"error": "User not found"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
        c.JSON(400, gin.H{"error": "Old password wrong"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to hash new password"})
        return
    }

    user.Password = string(hashedPassword)
    if err := db.DB.Save(&user).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to update password"})
        return
    }

    c.JSON(200, gin.H{
        "message": "Password updated successfully!",
    })
}



// _________________________________________________________________________________________________
// log last login
// Login log
func createLoginLog(c *gin.Context, userID uint) {
    log := models.LoginLog{
        UserID:    userID,
        IPAddress: c.ClientIP(),
        Device:    c.Request.Header.Get("X-Device"),
        Location:  c.Request.Header.Get("X-Location"),
        UserAgent: c.Request.UserAgent(),
    }

    db.DB.Create(&log)
}


// _________________________________________________________________________________________________
// auth
// Google
type GoogleCodeInput struct {
	Code string `json:"code" binding:"required"`
}

func GoogleLogin(c *gin.Context) {

	var input GoogleCodeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenURL := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("code", input.Code)
	data.Set("client_id", "34179357585-ojeb7n28gu1doapa3drn6db8hsjdhfpk.apps.googleusercontent.com")
	data.Set("client_secret", "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")
	data.Set("redirect_uri", "https://developers.google.com/oauthplayground")
	data.Set("grant_type", "authorization_code")

	req, _ := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed request to Google"})
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var tokenResponse map[string]interface{}
	json.Unmarshal(body, &tokenResponse)

	fmt.Println("Google Token Response: ", string(body))

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":            "Failed to get access token",
			"google_response":  tokenResponse,
		})
		return
	}

	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken
	res2, err := http.Get(userInfoURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed get Google user"})
		return
	}

	body2, _ := io.ReadAll(res2.Body)

	var googleUser map[string]interface{}
	json.Unmarshal(body2, &googleUser)

	email, _ := googleUser["email"].(string)
	name, _ := googleUser["name"].(string)
	avatar, _ := googleUser["picture"].(string)
	providerID, _ := googleUser["id"].(string)

	var user models.User
	db.DB.Where("email = ?", email).First(&user)

	if user.ID == 0 {
		user = models.User{
			Name:       name,
			Email:      email,
			Password:   "",
			Provider:   "google",
			Birthdate:  nil,
			ProviderID: providerID,
			Avatar:     avatar,
		}

		db.DB.Create(&user)
	}

	db.DB.Exec(`
		INSERT INTO login_logs (user_id, device, ip_address, location, logged_in_at)
		VALUES (?, ?, ?, ?, ?)
	`, user.ID, c.GetHeader("User-Agent"), c.ClientIP(), "Unknown", time.Now())

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Google Success",
		"user":    user,
	})
}

// FB
type FacebookTokenInput struct {
	AccessToken string `json:"access_token" binding:"required"`
}

func FacebookLogin(c *gin.Context) {

	var input FacebookTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInfoURL := "https://graph.facebook.com/me?fields=id,name,email,picture&access_token=" + input.AccessToken

	res, err := http.Get(userInfoURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed request to Facebook"})
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println("Facebook User Response: ", string(body))

	var fbUser map[string]interface{}
	json.Unmarshal(body, &fbUser)

	if fbUser["error"] != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "Facebook token invalid",
			"facebook_error": fbUser["error"],
		})
		return
	}

	email, _ := fbUser["email"].(string)
	name, _ := fbUser["name"].(string)
	providerID, _ := fbUser["id"].(string)

	avatar := ""
	if picture, ok := fbUser["picture"].(map[string]interface{}); ok {
		if data, ok := picture["data"].(map[string]interface{}); ok {
			avatar, _ = data["url"].(string)
		}
	}

	var user models.User
	db.DB.Where("email = ?", email).First(&user)

	if user.ID == 0 {
		user = models.User{
			Name:       name,
			Email:      email,
			Password:   "",
			Provider:   "facebook",
			ProviderID: providerID,
			Avatar:     avatar,
			Birthdate:  nil,
		}
		db.DB.Create(&user)
	}

	db.DB.Exec(`
		INSERT INTO login_logs (user_id, device, ip_address, location, logged_in_at)
		VALUES (?, ?, ?, ?, ?)
	`, user.ID, c.GetHeader("User-Agent"), c.ClientIP(), "Unknown", time.Now())

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Facebook Success",
		"user":    user,
	})
}


// _________________________________________________________________________________________________
// update
type UpdateUserInput struct {
    Name      string `json:"name"`
    Username  string `json:"username"`
    Gender    string `json:"gender"`
    Birthdate string `json:"birthdate"`
    Phone     string `json:"phone"`
    Bio       string `json:"bio"`
    Country   string `json:"country"`
    Address   string `json:"address"`
    Password  string `json:"password"`
}

func UpdateUser(c *gin.Context) {
    var input UpdateUserInput
    userID := c.Param("id")

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := db.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    user.Name = input.Name
    user.Username = input.Username
    user.Gender = input.Gender
    // user.Birthdate = input.Birthdate
    if input.Birthdate != "" {
        parsedDate, err := time.Parse("2006-01-02", input.Birthdate)
        if err != nil {
            c.JSON(400, gin.H{"error": "Invalid birthdate format (YYYY-MM-DD)"})
            return
        }
        user.Birthdate = &parsedDate
    } else {
        user.Birthdate = nil
    }
    user.Phone = input.Phone
    user.Bio = input.Bio
    user.Country = input.Country
    user.Address = input.Address

    if input.Password != "" {
        hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
        user.Password = string(hashed)
    }

    if err := db.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "User updated successfully",
        "user":    user,
    })
}


// _________________________________________________________________________________________________
// upload photo
// func UploadAvatar(c *gin.Context) {
//     userID := c.PostForm("user_id")
//     if userID == "" {
//         c.JSON(400, gin.H{"error": "user_id is required"})
//         return
//     }

//     file, err := c.FormFile("photo")
//     if err != nil {
//         c.JSON(400, gin.H{"error": "photo file is required"})
//         return
//     }

//     if !strings.Contains(file.Header.Get("Content-Type"), "image") {
//         c.JSON(400, gin.H{"error": "Only image files are allowed"})
//         return
//     }

//     filename := fmt.Sprintf("%s-%d%s",
//         userID,
//         time.Now().Unix(),
//         filepath.Ext(file.Filename),
//     )

//     uploadPath := "uploads/avatars"
//     os.MkdirAll(uploadPath, os.ModePerm)

//     fullPath := uploadPath + "/" + filename
//     baseURL := "https://testtestdomaingweh.com/"
//     avatarURL := baseURL + fullPath

//     if err := c.SaveUploadedFile(file, fullPath); err != nil {
//         c.JSON(500, gin.H{"error": "Failed to save file"})
//         return
//     }

//     db.DB.Model(&models.User{}).
//         Where("id = ?", userID).
//         Update("avatar", avatarURL)

//     c.JSON(200, gin.H{
//         "message": "Profile photo updated",
//         "avatar_url": avatarURL,
//     })
// }
func UploadAvatar(c *gin.Context) {
    userID := c.PostForm("user_id")
    if userID == "" {
        c.JSON(400, gin.H{"error": "user_id is required"})
        return
    }

    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(400, gin.H{"error": "Failed to read photos"})
        return
    }

    files := form.File["photos"]
    if len(files) == 0 {
        c.JSON(400, gin.H{"error": "No photos uploaded"})
        return
    }

    uploadPath := "uploads/users"
    baseURL := "https://testtestdomaingweh.com/"
    os.MkdirAll(uploadPath, os.ModePerm)

    uploadedURLs := []string{}
    var user models.User

    if err := db.DB.First(&user, userID).Error; err != nil {
        c.JSON(404, gin.H{"error": "User not found"})
        return
    }

    for _, file := range files {

        if !strings.Contains(file.Header.Get("Content-Type"), "image") {
            continue
        }

        filename := fmt.Sprintf("%s-%d-%s%s",
            userID,
            time.Now().Unix(),
            uuid.New().String(),
            filepath.Ext(file.Filename),
        )

        fullPath := uploadPath + "/" + filename
        photoURL := baseURL + fullPath

        if err := c.SaveUploadedFile(file, fullPath); err != nil {
            continue
        }

        db.DB.Create(&models.UserPhoto{
            UserID: cast.ToUint(userID),
            Photo:  photoURL,
        })

        uploadedURLs = append(uploadedURLs, photoURL)

        if user.Avatar == "" {
            db.DB.Model(&user).Update("avatar", photoURL)
        }
    }

    c.JSON(200, gin.H{
        "message": "Photos uploaded successfully",
        "photos":  uploadedURLs,
        "avatar":  user.Avatar,
    })
}

func GetUserPhotos(c *gin.Context) {
    userID := c.Param("id")

    var photos []models.UserPhoto
    db.DB.Where("user_id = ?", userID).Find(&photos)

    c.JSON(200, gin.H{
        "photos": photos,
    })
}

func SetProfilePhoto(c *gin.Context) {
    userID := c.PostForm("user_id")
    photoID := c.PostForm("photo_id")

    if userID == "" || photoID == "" {
        c.JSON(400, gin.H{"error": "user_id & photo_id required"})
        return
    }

    var photo models.UserPhoto
    if err := db.DB.Where("id = ? AND user_id = ?", photoID, userID).First(&photo).Error; err != nil {
        c.JSON(404, gin.H{"error": "Photo not found"})
        return
    }

    db.DB.Model(&models.User{}).
        Where("id = ?", userID).
        Update("avatar", photo.Photo)

    c.JSON(200, gin.H{
        "message": "Avatar updated",
        "avatar":  photo.Photo,
    })
}

func DeleteUserPhoto(c *gin.Context) {
    photoID := c.Param("id")

    var photo models.UserPhoto
    if err := db.DB.First(&photo, photoID).Error; err != nil {
        c.JSON(404, gin.H{"error": "Photo not found"})
        return
    }

    var user models.User
    db.DB.First(&user, photo.UserID)

    parts := strings.Split(photo.Photo, "/")
    filename := parts[len(parts)-1]
    os.Remove("uploads/users/" + filename)

    db.DB.Delete(&photo)

    if user.Avatar == photo.Photo {
        defaultPhoto := "https://testtestdomaingweh.com/default-avatar.png"
        db.DB.Model(&user).Update("avatar", defaultPhoto)
    }

    c.JSON(200, gin.H{
        "message": "Photo deleted",
    })
}

// _________________________________________________________________________________________________
// user detail
func GetUserDetail(c *gin.Context) {
    userID := c.Param("id")
	var users models.User

	if err := db.DB.Where(&userID).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User data found",
		"users":   users,
	})
}


// _________________________________________________________________________________________________
// forgot password
func ForgotPasswordRequest(c *gin.Context) {
    var input struct {
        Email string `json:"email" binding:"required,email"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.JSON(404, gin.H{"error": "Email not registered"})
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
        <h3>Reset Password</h3>
        <p>OTP kamu adalah:</p>
        <h2>%s</h2>
        <p>Berlaku selama 10 menit</p>
    `, otp)

    helpers.SendEmail(input.Email, "Reset Password OTP", emailBody)

    fmt.Println("FORGOT OTP:", otp)

    c.JSON(200, gin.H{
        "message": "OTP sent to email",
    })
}

func ForgotPasswordVerify(c *gin.Context) {
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

    c.JSON(200, gin.H{
        "message": "OTP verified",
    })
}

func ResetPassword(c *gin.Context) {
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
        c.JSON(400, gin.H{"error": "Invalid OTP"})
        return
    }

    if time.Now().After(otpData.ExpiredAt) {
        c.JSON(400, gin.H{"error": "OTP expired"})
        return
    }

    var user models.User
    if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.JSON(404, gin.H{"error": "User not found"})
        return
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to hash password"})
        return
    }

    user.Password = string(hashed)
    db.DB.Save(&user)

    // hapus OTP setelah sukses
    db.DB.Delete(&otpData)

    helpers.SendEmail(
        user.Email,
        "Password Updated Successfully",
        "<p>Password kamu berhasil diubah</p>",
    )

    c.JSON(200, gin.H{
        "message": "Password reset successfully",
    })
}


// _________________________________________________________________________________________________
// report
func GetReportReasons(c *gin.Context) {
	var reasons []models.ReportReason

	if err := db.DB.
		Where("is_active = true").
		Order("id ASC").
		Find(&reasons).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch report reasons",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": reasons,
	})
}

func ReportUser(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input struct {
		ReportedUserID uint   `json:"reported_user_id" binding:"required"`
		ReasonID       uint   `json:"reason_id" binding:"required"`
		Description    string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if userID == input.ReportedUserID {
		c.JSON(400, gin.H{"error": "cannot report yourself"})
		return
	}

	// Cegah spam report (1 user hanya boleh report 1x user yg sama)
	var exist models.UserReport
	if err := db.DB.Where(
		"reporter_id = ? AND reported_user_id = ?",
		userID, input.ReportedUserID,
	).First(&exist).Error; err == nil {
		c.JSON(400, gin.H{"error": "you already reported this user"})
		return
	}

	report := models.UserReport{
		ReporterID:     userID,
		ReportedUserID: input.ReportedUserID,
		ReasonID:       input.ReasonID,
		Description:    input.Description,
	}

	db.DB.Create(&report)

	c.JSON(201, gin.H{
		"message": "report submitted successfully",
	})
}



// _________________________________________________________________________________________________
// test
func GetAllUsers(c *gin.Context) {
	var users []models.User

	if err := db.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve users",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All user data retrieved",
		"data":    users,
	})
}

func GetPublicUsers(c *gin.Context) {
	type PublicUser struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}

	var users []PublicUser

	if err := db.DB.
		Model(&models.User{}).
		Select("id, name, username, email, phone").
		Find(&users).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve users",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Public user data retrieved",
		"data":    users,
	})
}


func Ceklog(c *gin.Context) {
	var login_logs []models.LoginLog

	if err := db.DB.Find(&login_logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User data found",
		"users":   login_logs,
	})
}
