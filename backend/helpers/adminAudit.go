package helpers

import (
	"encoding/json"
	"backend/db"
	"backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type AuditPayload struct {
	AdminID    uint
	Action     string
	TargetType string
	TargetID   *uint
	Before     interface{}
	After      interface{}
	Context    *gin.Context
}

func CreateAdminAuditLog(payload AuditPayload) error {
	var beforeJSON datatypes.JSON
	var afterJSON datatypes.JSON

	if payload.Before != nil {
		b, _ := json.Marshal(payload.Before)
		beforeJSON = b
	}

	if payload.After != nil {
		a, _ := json.Marshal(payload.After)
		afterJSON = a
	}

	ip := ""
	ua := ""

	if payload.Context != nil {
		ip = payload.Context.ClientIP()
		ua = payload.Context.GetHeader("User-Agent")
	}

	log := models.AdminAuditLog{
		AdminID:    payload.AdminID,
		Action:     payload.Action,
		TargetType: payload.TargetType,
		TargetID:   payload.TargetID,
		BeforeData: beforeJSON,
		AfterData:  afterJSON,
		IPAddress:  ip,
		UserAgent:  ua,
	}

	return db.DB.Create(&log).Error
}


func UintPtr(v uint) *uint {
	return &v
}