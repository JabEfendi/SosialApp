package main

import (
    "backend/db"
    "backend/routes"
    "backend/firebase"
    "os"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
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

    port := os.Getenv("APP_PORT")
    r.Run("0.0.0.0:" + port)
}

