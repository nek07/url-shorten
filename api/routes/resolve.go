package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/nek07/url-shorten/api/database"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	// ctx := context.Background() // Используем контекст для Redis операций

	// Создаем Redis-клиент для основной базы данных
	r := database.CreateClient(0)
	defer r.Close()

	// Получаем значение из базы данных
	value, err := r.Get(url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short not found in the database"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to DB"})
	}

	// Инкремент счетчика использования
	rInr := database.CreateClient(1)
	defer rInr.Close()

	if incrErr := rInr.Incr("counter").Err(); incrErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot update usage counter"})
	}

	// Перенаправление на оригинальный URL
	return c.Redirect(value, fiber.StatusTemporaryRedirect)
}
