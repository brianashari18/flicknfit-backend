package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetUserID(c *fiber.Ctx) (uint64, error) {
	rawID := c.Locals("userID")
	if rawID == nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "missing user id in context")
	}

	switch v := rawID.(type) {
	case uint64:
		return v, nil
	case int:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case float64:
		return uint64(v), nil
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, fiber.NewError(fiber.StatusUnauthorized, "invalid user id format")
		}
		return parsed, nil
	default:
		return 0, fiber.NewError(fiber.StatusUnauthorized, "invalid user id type")
	}
}
