package utils

import (
	crypto "crypto/rand"
	"encoding/base32"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var randRead = crypto.Read

func AuthUser(c *fiber.Ctx) jwt.MapClaims {
	return c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CustomValidationMessages(err error) string {
	errMessages := make([]string, 0)

	for _, err := range err.(validator.ValidationErrors) {
		errMsg := fmt.Sprintf("Field %s ", err.Field())

		switch err.Tag() {
		case "required":
			errMsg += "is required"
		case "gte":
			errMsg += fmt.Sprintf("must be greater than or equal to %s", err.Param())
		case "lte":
			errMsg += fmt.Sprintf("must be less than or equal to %s", err.Param())
		case "email":
			errMsg += "must be a valid email address"
		case "min":
			errMsg += fmt.Sprintf("must be minimum of %s", err.Param())
		default:
			errMsg += "failed validation"
		}
		errMessages = append(errMessages, strings.ToLower(errMsg))
		fmt.Println(err.Field(), err.Tag())
	}

	return strings.Join(errMessages, ", ")
}

func RandomID(prefix string) (string, error) {
	buf, err := RandBytes(32)
	if err != nil {
		return "", err
	}
	str := base32.StdEncoding.EncodeToString(buf)
	str = strings.ReplaceAll(str, "=", "")
	str = prefix + str
	return str, nil
}

// RandBytes returns random bytes of length
func RandBytes(length int) ([]byte, error) {
	buf := make([]byte, length)
	if _, err := randRead(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func IsObjectID(v interface{}) bool {
	valType := reflect.TypeOf(v)
	objectIDType := reflect.TypeOf(primitive.ObjectID{})

	return valType == objectIDType
}
