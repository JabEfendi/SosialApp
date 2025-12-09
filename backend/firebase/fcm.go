package firebase

import (
    "context"
    "log"

    "firebase.google.com/go/messaging"
    firebase "firebase.google.com/go"
    "google.golang.org/api/option"
)

var App *firebase.App

func InitFirebase() {
    opt := option.WithCredentialsFile("serviceAccountKey.json")

    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        log.Fatalf("Firebase init failed: %v", err)
    }

    App = app
    log.Println("Firebase Firestore initialized!")
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
