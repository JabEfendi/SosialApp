package controllers

import (
    "context"
    "strconv"
    "time"
    "backend/db"
    "backend/firebase"
    "backend/models"

    "github.com/gin-gonic/gin"
)

type SendMessageInput struct {
    RoomID  int   `json:"room_id"`
    UserID  int   `json:"user_id"`
    Message string `json:"message"`
}

func SendMessage(c *gin.Context) {
    var input SendMessageInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // VALIDASI: user harus sudah join room
    var participant models.RoomParticipant
    if err := db.DB.Where("room_id = ? AND user_id = ?", input.RoomID, input.UserID).First(&participant).Error; err != nil {
        c.JSON(403, gin.H{"error": "User not part of this room"})
        return
    }

    // SIMPAN KE POSTGRES
    msg := models.ChatMessage{
        RoomID:  input.RoomID,
        UserID:  input.UserID,
        Message: input.Message,
    }

    if err := db.DB.Create(&msg).Error; err != nil {
        c.JSON(500, gin.H{"error": "Database insert failed"})
        return
    }

    // FIRESTORE
    ctx := context.Background()
    fs, err := firebase.App.Firestore(ctx)
    if err != nil {
        c.JSON(500, gin.H{"error": "Firestore init failed"})
        return
    }
    defer fs.Close()

    data := map[string]interface{}{
        "id":         msg.ID,
        "room_id":    msg.RoomID,
        "user_id":    msg.UserID,
        "message":    msg.Message,
        "created_at": time.Now(),
    }

		_, err = fs.Collection("rooms").
			Doc(strconv.Itoa(int(input.RoomID))).
			Collection("messages").
			NewDoc().
			Set(ctx, data)

    if err != nil {
        c.JSON(500, gin.H{"error": "Firestore push failed", "detail": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "status":  "success",
        "message": "Message sent",
        "data":    msg,
    })
}

func GetMessages(c *gin.Context) {
    roomID := c.Param("roomID")

    var messages []models.ChatMessage
    db.DB.Where("room_id = ?", roomID).
        Order("created_at ASC").
        Find(&messages)

    c.JSON(200, messages)
}

func GetRealtimeStream(c *gin.Context) {
    roomID := c.Param("roomID")

    url := "rooms/" + roomID + "/messages"

    c.JSON(200, gin.H{
        "firestore_path": url,
    })
}

// ___________________________________________________________________________________________
// CONTROLLER UNTUK DIRECT
type DirectMessageInput struct {
    SenderID   uint   `json:"sender_id"`
    ReceiverID uint   `json:"receiver_id"`
    Message    string `json:"message"`
}

func getOrCreateThread(user1, user2 uint) (uint, error) {
    var thread models.ChatDirectThread

    u1 := user1
    u2 := user2
    if u1 > u2 {
        u1, u2 = u2, u1
    }

    err := db.DB.Where("user1_id = ? AND user2_id = ?", u1, u2).First(&thread).Error
    if err == nil {
        return thread.ID, nil
    }

    newThread := models.ChatDirectThread{
        User1ID: u1,
        User2ID: u2,
    }

    if err := db.DB.Create(&newThread).Error; err != nil {
        return 0, err
    }

    return newThread.ID, nil
}

func SendDirectMessage(c *gin.Context) {
    var input DirectMessageInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    threadID, err := getOrCreateThread(input.SenderID, input.ReceiverID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed creating thread"})
        return
    }

    msg := models.ChatDirectMessage{
        ThreadID:   threadID,
        SenderID:   input.SenderID,
        ReceiverID: input.ReceiverID,
        Message:    input.Message,
        Status:     "sent",
    }

    if err := db.DB.Create(&msg).Error; err != nil {
        c.JSON(500, gin.H{"error": "Insert DB failed"})
        return
    }

    c.JSON(200, gin.H{
        "status":  "success",
        "message": "Direct message sent",
        "data":    msg,
    })
}

func GetDirectMessages(c *gin.Context) {
    threadID := c.Param("threadID")

    var messages []models.ChatDirectMessage

    db.DB.Where("thread_id = ?", threadID).
        Order("created_at ASC").
        Find(&messages)

    c.JSON(200, messages)
}

func MarkDirectDelivered(c *gin.Context) {
    threadID := c.Param("threadID")
    userID := c.Param("userID")

    db.DB.Model(&models.ChatDirectMessage{}).
        Where("thread_id = ? AND receiver_id = ?", threadID, userID).
        Where("status = ?", "sent").
        Update("status", "delivered")

    c.JSON(200, gin.H{
        "status": "success",
        "message": "Messages marked as delivered",
    })
}
