/**
 * Author : ngdangkietswe
 * Since  : 8/14/2025
 */

package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
	"time"
)

func NewLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "[${time}] ${status} | ${latency} | ${ip} | ${method} ${path} ${queryParams} | ${bytesSent}B ${error}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "Asia/Ho_Chi_Minh",
		Output:     os.Stdout,
	})
}
