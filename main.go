package main

import (
	"sensor-server/db"
	"sensor-server/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	db.Init()

	e := echo.New()

	e.POST("/data", handlers.InsertSensor)
	e.GET("/data", handlers.GetSensor)
	e.GET("/data/latest", handlers.GetLatest)
	e.GET("/data/stats", handlers.GetStats)
	e.GET("/data/group", handlers.GetGrouped)
	e.PATCH("/data", handlers.UpdateSensor)
	e.DELETE("/data", handlers.DelSensor)

	e.Logger.Fatal(e.Start(":8080"))
}
