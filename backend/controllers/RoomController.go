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

type InputRoom struct {
    UserID           uint    `json:"user_id" binding:"required"`
    MapLocID         uint    `json:"map_loc_id" binding:"required"`
    Name             string  `json:"name" binding:"required"`
    Description      string  `json:"description" binding:"required"`
    StartTime        string  `json:"start_time" binding:"required,datetime=2006-01-02 15:04:05"`
    EndTime          string  `json:"end_time"`
    Duration         int  `json:"duration" binding:"required"`
    Capacity         int     `json:"capacity"`
    FeePerPerson     float64 `json:"fee_per_person"`
    Gender           string  `json:"gender"`
    AgeGroup         string  `json:"age_group"`
    IsRegular        bool    `json:"is_regular"`
    AutoApprove      bool    `json:"auto_approve"`
    SendNotification bool    `json:"send_notification"`
    ImageURL         string  `json:"image_url"`
}

func CreateRoom(c *gin.Context) {
    var input InputRoom
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Validation failed",
            "details": err.Error(),
        })
        return
    }

    var kyc models.KycSession
    if err := db.DB.Where("user_id = ? AND used = false", input.UserID).First(&kyc).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "KYC has not been submitted or has been used"})
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
    
    var mapLoc models.MapLoc
    if err := db.DB.
        Where("id = ? AND is_active = true", input.MapLocID).
        First(&mapLoc).Error; err != nil {

        c.JSON(400, gin.H{"error": "Invalid or inactive map location"})
        return
    }

    room := models.Room{
        UserID:           input.UserID,
        MapLocID:         &input.MapLocID,
        Name:             input.Name,
        Description:      input.Description,
        StartTime:        startTime,
        EndTime:          endTime,
        Duration:         time.Duration(input.Duration) * time.Minute,
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
        c.JSON(400, gin.H{"error": "Rooms are not active"})
        return
    }

    var count int64
    db.DB.Model(&models.RoomParticipant{}).Where("room_id = ?", room.ID).Count(&count)

    if room.Capacity > 0 && int(count) >= room.Capacity {
        c.JSON(400, gin.H{"error": "The room is full"})
        return
    }

    var existing models.RoomParticipant
    if err := db.DB.Where("room_id = ? AND user_id = ?", room.ID, input.UserID).First(&existing).Error; err == nil {
        c.JSON(400, gin.H{"error": "User has joined this room"})
        return
    }
    var checkuser models.User
    if err := db.DB.Where("id = ?", input.UserID).First(&checkuser).Error; err != nil {
        c.JSON(400, gin.H{"error": "Please Login First"})
        return
    }

    participant := models.RoomParticipant{
        RoomID: input.RoomID,
        UserID: input.UserID,
    }

    db.DB.Create(&participant)

    c.JSON(200, gin.H{
        "message": "Successfully joined the room",
        "data":    participant,
    })
}


func GetRoomParticipants(c *gin.Context) {
    roomID := c.Param("id")

    var participants []models.RoomParticipant
    if err := db.DB.Where("room_id = ?", roomID).Find(&participants).Error; err != nil {
        c.JSON(400, gin.H{"error": "Failed to take participants"})
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

func ListRoom(c *gin.Context) {
    lat := c.Query("lat")
    lng := c.Query("lng")

    if lat == "" || lng == "" {
        c.JSON(400, gin.H{
            "error": "lat and lng are required",
        })
        return
    }

    radius := 10.0 // KM

    var rooms []models.Room

    query := `
        SELECT rooms.* FROM rooms
        JOIN map_locs ON map_locs.id = rooms.map_loc_id
        WHERE map_locs.is_active = true
        AND (
            6371 * acos(
                cos(radians(?)) * cos(radians(map_locs.latitude)) *
                cos(radians(map_locs.longitude) - radians(?)) +
                sin(radians(?)) * sin(radians(map_locs.latitude))
            )
        ) <= ?
    `

    if err := db.DB.
        Raw(query, lat, lng, lat, radius).
        Scan(&rooms).Error; err != nil {

        c.JSON(500, gin.H{"error": "Failed to fetch rooms"})
        return
    }

    for i := range rooms {
        db.DB.Preload("MapLoc").First(&rooms[i], rooms[i].ID)
    }

    c.JSON(200, gin.H{
        "message": "Rooms found nearby",
        "radius_km": radius,
        "total": len(rooms),
        "rooms": rooms,
    })
}

func DetailRoom(c *gin.Context) {
    roomID := c.Param("id")
    var room models.Room

    if err := db.DB.
        Preload("MapLoc").
        First(&room, roomID).Error; err != nil {

        c.JSON(404, gin.H{"error": "Room not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Data room found",
        "room":    room,
    })
}

func UpdateRoom(c *gin.Context){
    var input map[string]interface{}
    roomID := c.Param("id")

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var room models.Room
    if err := db.DB.First(&room, roomID).Error; err != nil {
        c.JSON(400, gin.H{"error": "There are no rooms"})
        return
    }

    // room.Name              = input.Name
    // room.Description       = input.Description
    // room.StartTime         = input.StartTime
    // room.EndTime           = input.EndTime
    // room.Duration          = input.Duration
    // room.Location          = input.Location
    // room.Capacity          = input.Capacity
    // room.FeePerPerson      = input.FeePerPerson 
    // room.Gender            = input.Gender
    // room.AgeGroup          = input.AgeGroup
    // room.IsRegular         = input.IsRegular
    // room.AutoApprove       = input.AutoApprove
    // room.SendNotification  = input.SendNotification

    // if err := db.DB.Save(&room).Error; err != nil {
    //     c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
    //     return
    // }

    if mapLocID, ok := input["map_loc_id"]; ok {
        var mapLoc models.MapLoc

        if err := db.DB.
            Where("id = ? AND is_active = true", mapLocID).
            First(&mapLoc).Error; err != nil {

            c.JSON(400, gin.H{
                "error": "Location not found or inactive",
            })
            return
        }
    }

    if err := db.DB.Model(&room).Updates(input).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
        return
    }

    db.DB.Preload("MapLoc").First(&room, roomID)
    c.JSON(http.StatusOK, gin.H{
		"message": "Update room successfully",
		"Room":   room,
	})
}