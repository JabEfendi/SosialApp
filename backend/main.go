package main

import (
    "backend/db"
    "backend/routes"
    "backend/firebase"
    "backend/controllers"
    "os"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/robfig/cron/v3"
)

func main() {
    if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env not found, using system environment")
	}

	log.Println("APP NAME =", os.Getenv("APP_NAME"))

    r := gin.Default()

    db.ConnectDB()
    
    firebase.InitFirebase()

    routes.SystemRoutes(r)

    c := cron.New()

    _, err := c.AddFunc("0 0 * * *", func() {
        err := controllers.ReleaseRoomRevenue()
        if err != nil {
            log.Println("Error releasing room revenue:", err)
        } else {
            log.Println("Room revenue released successfully")
        }
    })
    if err != nil {
        log.Fatal("Failed to schedule cron job:", err)
    }

    c.Start()
    log.Println("Cron job for releasing room revenue started")

    port := os.Getenv("APP_PORT")
    r.Run("0.0.0.0:" + port)
}