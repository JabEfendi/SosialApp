package controllers

import (
	"backend/db"
	"backend/models"
	"net/http"
  "bytes"
	"encoding/json"
	"io/ioutil"
	"time"
  "backend/firebase"
	"fmt"

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
    // Birthdate string `json:"birthdate"`
    Phone     string `json:"phone"`
    Bio       string `json:"bio"`
    Country   string `json:"country"`
    Address   string `json:"address"`
    FCMToken  string `gorm:"column:fcmtoken" json:"fcm_token"`
}

func Register(c *gin.Context) {
    var input RegisterInput
		var birthdate time.Time
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

		// if input.Birthdate == "" {
		// 		birthdate = "0001-01-01"
		// } else {
		// 	birthdate = input.Birthdate
		// }

    user := models.User{
        Name:      input.Name,
        Username:  input.Username,
        Email:     input.Email,
        Password:  string(hashedPassword),
        Gender:    input.Gender,
        Birthdate: birthdate,
        Phone:     input.Phone,
        Bio:       input.Bio,
        Country:   input.Country,
        Address:   input.Address,
				FCMToken:  input.FCMToken,
        Provider:  "local",
    }

    // db.DB.Create(&user)
		result := db.DB.Create(&user)
		fmt.Println("Error:", result.Error)
		fmt.Println("RowsAffected:", result.RowsAffected)


    if user.FCMToken != "" {
        messageID, err := firebase.SendNotification(user.FCMToken, "Akun Terbuat", "Selamat "+user.Name+", akun Anda berhasil dibuat!")
        if err != nil {
            c.JSON(http.StatusOK, gin.H{
                "message":   "Register sukses, tapi gagal kirim notifikasi",
                "user":      user,
                "fcm_error": err.Error(),
            })
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "message":        "Register sukses dan notifikasi berhasil dikirim",
            "user":           user,
            "fcm_message_id": messageID,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Register sukses (tanpa notifikasi karena token kosong)",
        "user":    user,
    })
}



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


// BELOM KELAR DARI SINI
// MASIH GAGAL DI TOKEN GOOGLE & FACEBOOK NYA TAPI CONTROLLER NYA UDAH BENER
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

	//Tukar CODE → ACCESS TOKEN
	tokenURL := "https://oauth2.googleapis.com/token"

	payload := map[string]string{
		"code":          input.Code,
		"client_id":     "34179357585-ojeb7n28gu1doapa3drn6db8hsjdhfpk.apps.googleusercontent.com",
		"client_secret": "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",
		"redirect_uri":  "https://developers.google.com/oauthplayground",
		"grant_type":    "authorization_code",
	}

	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", tokenURL, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed exchange code"})
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var tokenResponse map[string]interface{}
	json.Unmarshal(body, &tokenResponse)

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed get access token"})
		return
	}

	//GET USER INFO DARI GOOGLE
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken

	res2, err := http.Get(userInfoURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed get Google user"})
		return
	}

	body2, _ := ioutil.ReadAll(res2.Body)

	var googleUser map[string]interface{}
	json.Unmarshal(body2, &googleUser)

	email := googleUser["email"].(string)
	name := googleUser["name"].(string)
	avatar := googleUser["picture"].(string)
	providerID := googleUser["id"].(string)

	//CEK USER SUDAH ADA BELUM
	var user models.User
	db.DB.Where("email = ?", email).First(&user)

	if user.ID == 0 {
		//Register baru otomatis
		user = models.User{
			Name:       name,
			Email:      email,
			Password:   "",

			Provider:   "google",
			ProviderID: providerID,
			Avatar:     avatar,
		}

		fmt.Println("Google User Body:", string(userBody))
		db.DB.Create(&user)
	}

	//UPDATE LAST LOGIN LOG
	db.DB.Exec(`
		INSERT INTO last_login_logs (user_id, device, ip_address, location, login_at)
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

	//Ambil user info dari Facebook Graph API
	userInfoURL := "https://graph.facebook.com/me?fields=id,name,email,picture&access_token=" + input.AccessToken

	res, err := http.Get(userInfoURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed get Facebook user"})
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var fbUser map[string]interface{}
	json.Unmarshal(body, &fbUser)

	email, _ := fbUser["email"].(string)
	name, _ := fbUser["name"].(string)
	avatar := ""
	if picture, ok := fbUser["picture"].(map[string]interface{}); ok {
		if data, ok := picture["data"].(map[string]interface{}); ok {
			avatar, _ = data["url"].(string)
		}
	}
	providerID, _ := fbUser["id"].(string)

	//Cek user sudah ada
	var user models.User
	db.DB.Where("email = ?", email).First(&user)

	if user.ID == 0 {
		//Register baru otomatis
		user = models.User{
			Name:       name,
			Email:      email,
			Password:   "",
			Provider:   "facebook",
			ProviderID: providerID,
			Avatar:     avatar,
		}
		db.DB.Create(&user)
	}

	//Update last login
	db.DB.Exec(`
		INSERT INTO last_login_logs (user_id, device, ip_address, location, login_at)
		VALUES (?, ?, ?, ?, ?)
	`, user.ID, c.GetHeader("User-Agent"), c.ClientIP(), "Unknown", time.Now())

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Facebook Success",
		"user":    user,
	})
}









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
