package controllers

import (
	"authx/initializers"
	"authx/model"
	"fmt"

	// "os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint64 `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Photo     string `json:"photo"`
}

func convertToJSON(user model.User) User {
	return User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Photo:     user.Photo,
	}
}

func Register(c *fiber.Ctx) error {
	c.Accepts("applicatin/json")

	var user model.User
	err := c.BodyParser(&user)
	password, _ := HashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "error parsing JSON " + err.Error(),
		})
	}
	user.Password = password
	err = model.CreateUser(&user)

	resp := convertToJSON(user)
	// resp.Token = "123dks"

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error creating user " + err.Error(),
		})
	}

	return c.Status(200).JSON(resp)
}

func Login(c *fiber.Ctx) error {
	c.Accepts("applicatin/json")

	var user model.User
	err := c.BodyParser(&user)
	password, _ := HashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "error reading your data " + err.Error(),
		})
	}

	userResp, err := model.FindUserWithPassword(user)

	match := CheckPasswordHash(userResp.Password, password)
	
	if !match || err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": " invalid login provided",
		})
	}

	resp := convertToJSON(userResp)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "user not found " + err.Error(),
		})
	}

	config, _ := initializers.LoadConfig(".")

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["userId"] = userResp.ID
	claims["exp"] = now.Add(config.JwtExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(config.JwtSecret))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("generating JWT Token failed: %v", err),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   config.JwtMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"token":  tokenString,
		"data":   resp,
	})

	// return c.Status(200).JSON(resp)

}

func LogoutUser(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// func validateJWT(tokenString string) (*jwt.Token, error) {
// 	secret := os.Getenv("JWT_SECRET")

// 	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 		}

// 		return []byte(secret), nil
// 	})
// }
