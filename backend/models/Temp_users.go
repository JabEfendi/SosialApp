package models

import "time"

type TempUser struct {
    ID        uint      `gorm:"primaryKey"`
    Email     string    `gorm:"unique"`
    Name      string
    Username  string
    Password  string
    Gender    string
    Birthdate string
    Phone     string
    Bio       string
    Country   string
    Address   string
    CreatedAt time.Time
}
