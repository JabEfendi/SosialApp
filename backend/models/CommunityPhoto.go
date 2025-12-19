package models

import "time"

type CommunityPhoto struct {
    ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    CommunityID uint   `json:"community_id"`
    Photo  string `json:"photo"`
    CreatedAt time.Time
}
