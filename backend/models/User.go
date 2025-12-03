package models

type User struct {
    ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Name       string `json:"name"`
    Username   string `json:"username"`
    Email      string `json:"email"`
    Password   string `json:"-"`
    Gender     string `json:"gender"`
    Birthdate  time.Time `json:"birthdate"`
    Phone      string `json:"phone"`
    Bio        string `json:"bio"`
    Country    string `json:"country"`
    Address    string `json:"address"`
    Provider   string `json:"provider"`
    ProviderID string `json:"provider_id"`
    Avatar     string `json:"avatar"`
    FCMToken string `gorm:"column:fcmtoken" json:"fcm_token"`
}
