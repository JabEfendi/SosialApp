package firebase

import (
    "context"
    "log"

    firebase "firebase.google.com/go"
    "firebase.google.com/go/messaging"
    "google.golang.org/api/option"
)

var App *firebase.App

func InitFirebase() {
    opt := option.WithCredentialsFile("serviceAccountKey.json")

    firebaseConfig := &firebase.Config{
        ProjectID: "sosialapp-reaction",
    }

    app, err := firebase.NewApp(context.Background(), firebaseConfig, opt)
    if err != nil {
        log.Fatalf("Error initializing Firebase: %v", err)
    }

    App = app
    log.Println("Firebase berhasil diinisialisasi!")
}

func SendNotification(token string, title string, body string) (string, error) {
    ctx := context.Background()
    client, err := App.Messaging(ctx)
    if err != nil {
        return "", err
    }

    message := &messaging.Message{
        Token: token,
        Notification: &messaging.Notification{
            Title: title,
            Body:  body,
        },
    }

    response, err := client.Send(ctx, message)
    if err != nil {
        return "", err
    }

    return response, nil
}
