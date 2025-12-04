package main

import (
    "backend/db"
    "backend/routes"
    "backend/firebase"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    db.ConnectDB()
    firebase.InitFirebase()

    routes.UserRoutes(r)
    routes.AuthRoutes(r)
    routes.SystemRoutes(r)

    r.Run(":8080")
}

