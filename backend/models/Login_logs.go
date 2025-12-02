package models

type LoginLog struct {
    ID         uint   `json:"id"`
    UserID     uint   `json:"user_id"`
    IPAddress  string `json:"ip_address"`
    Device     string `json:"device"`
    Location   string `json:"location"`
    UserAgent  string `json:"user_agent"`
    LoggedInAt string `json:"logged_in_at"`
}
