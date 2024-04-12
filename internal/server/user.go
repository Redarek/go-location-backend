package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"location-backend/internal/config"
	"time"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func (s *Fiber) Login(c *fiber.Ctx) error {
	loginInput := new(LoginInput)
	err := c.BodyParser(loginInput)
	if err != nil {
		return err
	}

	user, err := s.db.GetUserByUsername(loginInput.Username)

	log.Debug().Msgf("usernames: %v %v", loginInput.Username, user.Username)
	// Throws Unauthorized error
	if loginInput.Username != user.Username || !s.CheckPasswordHash(loginInput.Password, user.Password) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func (s *Fiber) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	log.Debug().Msgf("Password: %v", password)
	log.Debug().Msgf("HashPassword: %v", bytes)
	return string(bytes), err
}

func (s *Fiber) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	log.Debug().Msgf("Password: %v", password)
	log.Debug().Msgf("HashPassword: %v", hash)
	return err == nil
}

func (s *Fiber) Register(c *fiber.Ctx) error {
	registerInput := new(RegisterInput)
	err := c.BodyParser(registerInput)
	if err != nil {
		return err
	}
	hash, err := s.HashPassword(registerInput.Password)
	userID, err := s.db.CreateUser(registerInput.Username, hash)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"id": userID})
}
