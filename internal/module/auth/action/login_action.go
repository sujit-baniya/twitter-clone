package action

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/twitter-clone/internal/common/config"
	"github.com/sujit-baniya/twitter-clone/internal/module"
	"github.com/sujit-baniya/twitter-clone/internal/module/auth/service"
	"github.com/sujit-baniya/twitter-clone/internal/module/user/entity"
)

type loginAction struct {
	service service.LoginService
}

func NewLoginAction(service service.LoginService) module.Action {
	return loginAction{service: service}
}

func (a loginAction) Execute(c *fiber.Ctx) error {
	var input service.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if errors := input.Validate(); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	accessToken, refreshToken, err := a.service.Execute(input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserDoesNotExist) || errors.Is(err, service.ErrInvalidPassword):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid email/password",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "There was a problem on our side",
			})
		}
	}
	fmt.Println(refreshToken, accessToken)
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.String(),
		Expires:  refreshToken.ExpiresAt(),
		HTTPOnly: true,
		Secure:   config.GetString("APP_ENV", "development") == "production",
		Path:     "/",
		Domain:   config.GetString("APP_DOMAIN", ""),
		SameSite: "None",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken.String(),
	})
}
