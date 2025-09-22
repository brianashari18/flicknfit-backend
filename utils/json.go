package utils

import (
	"bytes"
	"encoding/json"
	"flicknfit_backend/dtos"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func StrictBodyParser(c *fiber.Ctx, out interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(c.Body()))
	dec.DisallowUnknownFields()
	return dec.Decode(out)
}

func SendResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(dtos.ResponseDTO{
		Status:  http.StatusText(status),
		Message: message,
		Data:    data,
	})
}
