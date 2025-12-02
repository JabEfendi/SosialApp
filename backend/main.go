package main

import (
    "backend/db"
    "backend/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    db.ConnectDB()

    routes.UserRoutes(r)

    r.Run(":8080")
}
