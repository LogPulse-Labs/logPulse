package auth_service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	auth_models "log-pulse/app/auth/models"
	"log-pulse/models"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	Access  string
	Refresh string
}

func IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GenerateTokens(user *models.User) (*Tokens, error) {
	accessToken, err := generateAccessToken(user)
	if err != nil {
		fmt.Println("Error generating access token:", err)
		return nil, err
	}

	refreshToken, err := generateNewRefreshToken()
	if err != nil {
		return nil, err
	}

	return &Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func generateAccessToken(user *models.User) (string, error) {
	jwtExpiresIn, _ := strconv.Atoi(os.Getenv("JWT_SECRET_KEY_EXPIRES_IN"))

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Minute * time.Duration(jwtExpiresIn)).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return t, err
}

func generateNewRefreshToken() (string, error) {
	hash := sha256.New()

	refresh := os.Getenv("JWT_REFRESH_KEY") + time.Now().String()

	// See: https://pkg.go.dev/io#Writer.Write
	_, err := hash.Write([]byte(refresh))
	if err != nil {
		return "", err
	}

	hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))
	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(hoursCount)).Unix())

	// Create a new refresh token (sha256 string with salt + expire time).
	t := hex.EncodeToString(hash.Sum(nil)) + "." + expireTime

	return t, nil
}

func ParseRefreshToken(refreshToken string) (int64, error) {
	return strconv.ParseInt(strings.Split(refreshToken, ".")[1], 0, 64)
}

func ValidateJWTToken(tokenString string) (jwt.MapClaims, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return os.Getenv("JWT_SECRET"), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func AuthenticatedUser(user *models.User, organization *models.Organization) interface{} {
	return &auth_models.AuthResponsePayload{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Organization: struct {
			ID   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		}{
			ID:   organization.ID,
			Name: organization.Name,
		},
	}
}
