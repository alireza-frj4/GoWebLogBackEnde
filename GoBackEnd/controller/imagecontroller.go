package controller

import (
	"fmt"
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

var letters = []rune("abcdefghijklmnopqrsuvwxyz")

func randLetter(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Upload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["image"]
	fileName := ""

	for _, file := range files {

		fileName = randLetter(5) + "-" + file.Filename
		if err := c.SaveFile(file, fmt.Sprintf("./%s", file.Filename)); err != nil {
			return err
		}
	}
	return c.JSON(fiber.Map{
		"url": "http://localhost:8080/api/uploads/" + fileName,
	})

}
