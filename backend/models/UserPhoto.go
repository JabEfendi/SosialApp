package models

type UserPhoto struct {
    ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID uint   `json:"user_id"`
    Photo  string `json:"photo"`
}
