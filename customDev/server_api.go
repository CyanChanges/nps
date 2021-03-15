package customDev

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func ApiWebServer() {
	app := fiber.New()

	setupRoutes(app)

	_ = app.Listen(":" + ApiPort)
}

// Set Routes
func setupRoutes(app *fiber.App) {
	// set handler for index page
	app.Get("/api/freePort", GetFreePort) // set handler for index page
}

// 返回一个空闲可用端口, 注意防火墙开启端口
func GetFreePort(c *fiber.Ctx) error {
	availablePort := strconv.Itoa(FindFreePort())

	return c.SendString(availablePort)
}
