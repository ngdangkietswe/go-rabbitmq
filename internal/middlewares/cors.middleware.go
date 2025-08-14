/**
 * Author : ngdangkietswe
 * Since  : 8/14/2025
 */

package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:  "*", // or specific origins: "http://localhost:3000, https://yourdomain.com"
		AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Content-Length, Content-Type",
	})
}
