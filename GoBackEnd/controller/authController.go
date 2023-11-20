package controller

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alireza-frj4/BlogBackEnd/database"
	"github.com/alireza-frj4/BlogBackEnd/models"
	"github.com/alireza-frj4/BlogBackEnd/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func validateEmail(email string) bool {
	Re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return Re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
		return err
	}

	// Check if password is less than 6 characters
	if password, ok := data["password"].(string); ok && len(password) <= 6 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Password must be greater than 6 character",
		})
	}

	// Check if email is valid
	if email, ok := data["email"].(string); ok && !validateEmail(strings.TrimSpace(email)) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid email address",
		})
	}

	// Check if email already exists in the database
	email, ok := data["email"].(string)
	if ok {
		database.DB.Where("email=?", strings.TrimSpace(email)).First(&userData)
		if userData.Id != 0 {
			c.Status(400)
			return c.JSON(fiber.Map{
				"message": "Email Already exists",
			})
		}
	}

	// Create and set password for the user
	user := models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Email:     strings.TrimSpace(data["email"].(string)),
		Phone:     data["phone"].(string),
	}

	if password, ok := data["password"].(string); ok {
		user.SetPassword(password)
	}

	// Create the user in the database
	err := database.DB.Create(&user)
	if err != nil {
		log.Println(err)
	}

	// Send success response with the created user
	c.Status(200)
	return c.JSON(fiber.Map{
		"user":    user,
		"message": "Account created successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
		return err
	}
	//Check is email is in our database
	var user models.User
	database.DB.Where("email=?", data["email"]).First(&user)
	if user.Id == 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Email address doesn't excit, kindly create an account",
		})
	}
	if err := user.ComparePassword(data["password"].(string)); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}
	token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "you have successfully login",
		"user":    user,
	})
}

type Claims struct {
	jwt.StandardClaims
}
