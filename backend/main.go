package main

import (
    "backend/db"
    "backend/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // Connect to DB
    db.ConnectDB()

    // Register routes
    routes.UserRoutes(r)

    r.Run(":8080")
}
