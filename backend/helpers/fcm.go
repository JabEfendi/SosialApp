package helpers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

var FCMServerKey = "SERVER_KEY_MU"

type FCMMessage struct {
    To           string                 `json:"to"`
    Notification map[string]string      `json:"notification"`
    Data         map[string]interface{} `json:"data"`
}

func SendFCMToken(token, title, msg string) error {

    body := FCMMessage{
        To: token,
        Notification: map[string]string{
            "title": title,
            "body":  msg,
        },
        Data: map[string]interface{}{
            "click_action": "FLUTTER_NOTIFICATION_CLICK",
            "type":         "register_success",
        },
    }

    jsonBytes, _ := json.Marshal(body)

    req, _ := http.NewRequest("POST",
        "https://fcm.googleapis.com/fcm/send",
        bytes.NewBuffer(jsonBytes))

    req.Header.Set("Authorization", "key="+FCMServerKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)

    if err != nil {
        return err
    }
    defer resp.Body.Close()

    fmt.Println("FCM Response:", resp.Status)
    return nil
}
