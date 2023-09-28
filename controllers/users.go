package controllers

import (
	"authx/initializers"
	"authx/middleware"
	"authx/model"
	"context"
	"fmt"
	"log"

	// "net/http"
	// "os"
	// "path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
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
	fmt.Println(userResp.Password)
	fmt.Println(password)

	match := CheckPasswordHash(user.Password, userResp.Password)

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
		MaxAge:   config.JwtMaxAge * 6000,
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

func UploadPhoto(c *fiber.Ctx) error {
	file, err := c.FormFile("image")

	var user model.User
	// err := c.BodyParser(&user)

	if err != nil {
		log.Println("image upload error --> ", err)
		return c.JSON(fiber.Map{"status": 500, "message": "Server error", "data": nil})
	}

	// generate new uuid for image name
	uniqueId := uuid.New()

	// remove "- from imageName"
	filename := strings.Replace(uniqueId.String(), "-", "", -1)

	// extract image extension from original file filename
	fileExt := strings.Split(file.Filename, ".")[1]

	// generate image from filename and extension
	image := fmt.Sprintf("%s.%s", filename, fileExt)

	// Open file for operations
	f, err := file.Open()

	if err != nil {
		return err
	}

	// Upload S3
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("gox-uploads"),
		Key:    aws.String(image),
		Body:   f,
		ACL:    "public-read",
	})

	if err != nil {
		log.Println("image upload error --> ", err)
		return c.JSON(fiber.Map{"status": 500, "message": "Server error", "data": nil})
	}

	// Update Table
	url := result.Location
	loggedInUser := middleware.GetLoggedInUser(c)
	// fmt.Println(url)
	user.ID = loggedInUser.ID
	// user.Photo = url
	err = model.UploadUserPhoto(user, url)
	if err != nil {
		log.Println("image upload error --> ", err)
		return c.JSON(fiber.Map{"status": 500, "message": "Server error", "data": nil})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   result.Location,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	c.Accepts("applicatin/json")

	var user model.User
	err := c.BodyParser(&user)
	// password, _ := HashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "error parsing JSON " + err.Error(),
		})
	}
	// user.Password = password
	loggedInUser := middleware.GetLoggedInUser(c)
	// fmt.Println(url)
	user.ID = loggedInUser.ID
	err = model.UpdateUser(&user)

	resp := convertToJSON(user)
	// resp.Token = "123dks"

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error creating user " + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   resp,
	})
}

func Logout(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func FetchAllUsers(c *fiber.Ctx) error {
	userResp, err := model.GetAllUsers()
	fmt.Println(userResp)
	users := []User{}

	// resp := convertToJSON(userResp)
	for _, user := range userResp {
		conv := convertToJSON(user)
		users = append(users, conv)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error fetching users " + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
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

// func uploadImageToS3(bucketName string, imagePath string) (string, error) {
// 	// Should we upload to the filesystem or temp memory before pushing to S3?

// 	// Create a new session with default session credentials
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String("us-east-1")},
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	// Open the image file
// 	file, err := os.Open(imagePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Get the file size and content type
// 	fileInfo, err := file.Stat()
// 	if err != nil {
// 		return err
// 	}
// 	fileSize := fileInfo.Size()
// 	fileType := http.DetectContentType(make([]byte, 512), fileInfo.Name())

// 	// Create a new S3 service client
// 	svc := s3.New(sess)

// 	// Configure the S3 object input parameters
// 	input := &s3.PutObjectInput{
// 		Body:          file,
// 		Bucket:        aws.String(bucketName),
// 		Key:           aws.String(filepath.Base(imagePath)),
// 		ContentType:   aws.String(fileType),
// 		ContentLength: aws.Int64(fileSize),
// 	}

// 	fmt.Println(input)

// 	// Upload the image to S3
// 	_, err = svc.PutObject(input)
// 	if err != nil {
// 		return err
// 	}

// 	return input, nil
// }
