package controllers

import (
	"backend/db"
	"backend/models"
	"net/http"
	"time"
	"encoding/json"

	"gorm.io/datatypes"
	"github.com/gin-gonic/gin"
)

func CreateRoom(c *gin.Context) {
    var input struct {
        UserID           uint   `json:"user_id"`
        Name             string `json:"name"`
        Description      string `json:"description"`
        StartTime        string `json:"start_time"`
        EndTime          string `json:"end_time"`
        Duration         string `json:"duration"`
        Location         string `json:"location"`
        Capacity         int    `json:"capacity"`
        FeePerPerson     float64 `json:"fee_per_person"`
        Gender           string `json:"gender"`
        AgeGroup         string `json:"age_group"`
        IsRegular        bool   `json:"is_regular"`
        AutoApprove      bool   `json:"auto_approve"`
        SendNotification bool   `json:"send_notification"`
        ImageURL         string `json:"image_url"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var kyc models.KycSession
    if err := db.DB.Where("user_id = ? AND used = false", input.UserID).First(&kyc).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "KYC belum submit atau sudah dipakai"})
        return
    }

		layout := "2006-01-02 15:04:05"

		startTime, err := time.ParseInLocation(layout, input.StartTime, time.Local)
		if err != nil {
				c.JSON(400, gin.H{"error": "Invalid start_time format. Use: YYYY-MM-DD HH:MM:SS"})
				return
		}

		endTime, err := time.ParseInLocation(layout, input.EndTime, time.Local)
		if err != nil {
				c.JSON(400, gin.H{"error": "Invalid end_time format. Use: YYYY-MM-DD HH:MM:SS"})
				return
		}

    room := models.Room{
        UserID:           input.UserID,
        Name:             input.Name,
        Description:      input.Description,
        StartTime:        startTime,
        EndTime:          endTime,
        Duration:         input.Duration,
        Location:         input.Location,
        Capacity:         input.Capacity,
        FeePerPerson:     input.FeePerPerson,
        Gender:           input.Gender,
        AgeGroup:         input.AgeGroup,
        IsRegular:        input.IsRegular,
        AutoApprove:      input.AutoApprove,
        SendNotification: input.SendNotification,
        ImageURL:         input.ImageURL,
        Status:           "active",
    }
    db.DB.Create(&room)

    kyc.Used = true
    db.DB.Save(&kyc)

		summary := map[string]interface{}{
				"room_name": room.Name,
				"kyc_data":  json.RawMessage(kyc.DataJSON),
		}
		summaryBytes, _ := json.Marshal(summary)

    log := models.RoomLog{
        RoomID:    room.ID,
        UserID:    input.UserID,
        KycStatus: "valid",
        Summary_json:   datatypes.JSON(summaryBytes),
    }
    db.DB.Create(&log)
		db.DB.Delete(&kyc)

    c.JSON(http.StatusOK, gin.H{"message": "Room created", "room": room})
}


func JoinRoom(c *gin.Context) {
    var input struct {
        RoomID uint `json:"room_id"`
        UserID uint `json:"user_id"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var room models.Room
    if err := db.DB.First(&room, input.RoomID).Error; err != nil {
        c.JSON(404, gin.H{"error": "Room not found"})
        return
    }

    if room.UserID == input.UserID {
        c.JSON(400, gin.H{"error": "Owner tidak boleh join ke room sendiri"})
        return
    }

    if room.Status != "active" {
        c.JSON(400, gin.H{"error": "Room tidak aktif"})
        return
    }

    var count int64
    db.DB.Model(&models.RoomParticipant{}).Where("room_id = ?", room.ID).Count(&count)

    if int(count) >= room.Capacity {
        c.JSON(400, gin.H{"error": "Room sudah penuh"})
        return
    }

    var existing models.RoomParticipant
    if err := db.DB.Where("room_id = ? AND user_id = ?", room.ID, input.UserID).First(&existing).Error; err == nil {
        c.JSON(400, gin.H{"error": "User sudah join room ini"})
        return
    }
    var checkuser models.User
    if err := db.DB.Where("id = ?", input.UserID).First(&checkuser).Error; err != nil {
        c.JSON(400, gin.H{"error": "Mohon Login Terlebih Dahulu"})
        return
    }

    participant := models.RoomParticipant{
        RoomID: input.RoomID,
        UserID: input.UserID,
    }

    db.DB.Create(&participant)

    c.JSON(200, gin.H{
        "message": "Berhasil join room",
        "data":    participant,
    })
}


func GetRoomParticipants(c *gin.Context) {
    roomID := c.Param("id")

    var participants []models.RoomParticipant
    if err := db.DB.Where("room_id = ?", roomID).Find(&participants).Error; err != nil {
        c.JSON(400, gin.H{"error": "Gagal mengambil peserta"})
        return
    }

    var result []map[string]interface{}

    for _, p := range participants {
        var user models.User
        db.DB.First(&user, p.UserID)

        result = append(result, map[string]interface{}{
            "user_id":    user.ID,
            "name":       user.Name,
            "avatar":     user.Avatar,
            "joined_at":  p.CreatedAt,
        })
    }

    c.JSON(200, gin.H{
        "room_id":     roomID,
        "participants": result,
    })
}

func ListRoom(c *gin.Context){
    // roomID := c.Param("id")

    var rooms []models.Room
    if err := db.DB.Find(&rooms).Error; err != nil {
        c.JSON(400, gin.H{"error": "Tidak ada room"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
		"message": "Data room ditemukan",
		"users":   rooms,
	})
}
