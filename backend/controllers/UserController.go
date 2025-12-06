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
}

func RegisterRequest(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var existing models.User
    if err := db.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
        c.JSON(400, gin.H{"error": "Email sudah digunakan"})
        return
    }

    db.DB.Where("email = ?", input.Email).Delete(&models.TempUser{})
    db.DB.Where("email = ?", input.Email).Delete(&models.OTPVerification{})

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

    temp := models.TempUser{
        Email:    input.Email,
        Name:     input.Name,
        Username: input.Username,
        Password: string(hashedPassword),
        Gender:   input.Gender,
		Birthdate: "0001-01-01",
        Phone:    input.Phone,
        Bio:      input.Bio,
        Country:  input.Country,
        Address:  input.Address,
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
        "message": "OTP telah dikirim ke email",
        "email":   input.Email,
    })
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
        c.JSON(400, gin.H{"error": "OTP salah"})
        return
    }

    if time.Now().After(otpData.ExpiredAt) {
        c.JSON(400, gin.H{"error": "OTP expired, silakan resend"})
        return
    }

    var temp models.TempUser
    if err := db.DB.Where("email = ?", input.Email).First(&temp).Error; err != nil {
        c.JSON(400, gin.H{"error": "Data tidak ditemukan"})
        return
    }
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
    }

    db.DB.Create(&user)
    db.DB.Delete(&temp)
    db.DB.Delete(&otpData)


		// Ambil token FCM user
		var token models.UserFCMToken
		db.DB.Where("user_id = ?", user.ID).Last(&token)

		// Buat greeting message
		message := fmt.Sprintf("Halo %s! Akun kamu berhasil dibuat pada %s",
				user.Name,
				time.Now().Format("02 Jan 2006"),
		)

		// Simpan ke tabel notification
		notif := models.Notification{
				UserID: user.ID,
				Title:  "Akun Berhasil Dibuat 🎉",
				Message: message,
		}
		db.DB.Create(&notif)

		// Kirim push notif (jika ada token)
		if token.FCMToken != "" {
				helpers.SendFCMToken(token.FCMToken,
						"Akun Berhasil Dibuat 🎉",
						message,
				)
		}


    c.JSON(200, gin.H{
        "message": "Akun berhasil dibuat",
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

    c.JSON(200, gin.H{"message": "OTP baru dikirim"})
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
        c.JSON(400, gin.H{"error": "FCM token user tidak ditemukan"})
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
        "message": "Notifikasi terkirim",
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
        c.JSON(400, gin.H{"error": "User tidak ditemukan"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
        c.JSON(400, gin.H{"error": "Password lama salah"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(500, gin.H{"error": "Gagal meng-hash password baru"})
        return
    }

    user.Password = string(hashedPassword)
    if err := db.DB.Save(&user).Error; err != nil {
        c.JSON(500, gin.H{"error": "Gagal update password"})
        return
    }

    c.JSON(200, gin.H{
        "message": "Password berhasil diperbarui!",
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

	// WAJIB pakai form-urlencoded
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

	// DEBUG RESPONSE GOOGLE
	fmt.Println("Google Token Response: ", string(body))

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":            "Failed to get access token",
			"google_response":  tokenResponse,
		})
		return
	}

	// GET USER INFO
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken
	res2, err := http.Get(userInfoURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed get Google user"})
		return
	}

	body2, _ := io.ReadAll(res2.Body)

	var googleUser map[string]interface{}
	json.Unmarshal(body2, &googleUser)

	// SAFE CASTING
	email, _ := googleUser["email"].(string)
	name, _ := googleUser["name"].(string)
	avatar, _ := googleUser["picture"].(string)
	providerID, _ := googleUser["id"].(string)

	// CEK USER
	var user models.User
	db.DB.Where("email = ?", email).First(&user)

	if user.ID == 0 {
		user = models.User{
			Name:       name,
			Email:      email,
			Password:   "",
			Provider:   "google",
			Birthdate:  "0001-01-01",
			ProviderID: providerID,
			Avatar:     avatar,
		}

		db.DB.Create(&user)
	}

	// LOG
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

	// DEBUG
	fmt.Println("Facebook User Response: ", string(body))

	// Parse hasil JSON
	var fbUser map[string]interface{}
	json.Unmarshal(body, &fbUser)

	// Jika token salah / expired
	if fbUser["error"] != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "Facebook token invalid",
			"facebook_error": fbUser["error"],
		})
		return
	}

	// SAFE CAST
	email, _ := fbUser["email"].(string)
	name, _ := fbUser["name"].(string)
	providerID, _ := fbUser["id"].(string)

	// Picture
	avatar := ""
	if picture, ok := fbUser["picture"].(map[string]interface{}); ok {
		if data, ok := picture["data"].(map[string]interface{}); ok {
			avatar, _ = data["url"].(string)
		}
	}

	// CEK USER EXIST
	var user models.User
	db.DB.Where("email = ?", email).First(&user)

	if user.ID == 0 {
		// AUTO REGISTER
		user = models.User{
			Name:       name,
			Email:      email,
			Password:   "",
			Provider:   "facebook",
			ProviderID: providerID,
			Avatar:     avatar,
			Birthdate:  "0001-01-01", // biar konsisten
		}
		db.DB.Create(&user)
	}

	// SIMPAN LOG LOGIN
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
    user.Birthdate = input.Birthdate
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
func UploadAvatar(c *gin.Context) {
    userID := c.PostForm("user_id")
    if userID == "" {
        c.JSON(400, gin.H{"error": "user_id is required"})
        return
    }

    file, err := c.FormFile("photo")
    if err != nil {
        c.JSON(400, gin.H{"error": "photo file is required"})
        return
    }

    if !strings.Contains(file.Header.Get("Content-Type"), "image") {
        c.JSON(400, gin.H{"error": "Only image files are allowed"})
        return
    }

    filename := fmt.Sprintf("%s-%d%s",
        userID,
        time.Now().Unix(),
        filepath.Ext(file.Filename),
    )

    uploadPath := "uploads/avatars"
    os.MkdirAll(uploadPath, os.ModePerm)

    fullPath := uploadPath + "/" + filename
    baseURL := "https://testtestdomaingweh.com/"
    avatarURL := baseURL + fullPath

    if err := c.SaveUploadedFile(file, fullPath); err != nil {
        c.JSON(500, gin.H{"error": "Failed to save file"})
        return
    }

    db.DB.Model(&models.User{}).
        Where("id = ?", userID).
        Update("avatar", avatarURL)

    c.JSON(200, gin.H{
        "message": "Profile photo updated",
        "avatar_url": avatarURL,
    })
}


// _________________________________________________________________________________________________
// test
func Test(c *gin.Context) {
	var users []models.User

	if err := db.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data user ditemukan",
		"users":   users,
	})
}


func Ceklog(c *gin.Context) {
	var login_logs []models.LoginLog

	if err := db.DB.Find(&login_logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data user ditemukan",
		"users":   login_logs,
	})
}
