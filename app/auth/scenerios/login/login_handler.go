package login_scenerio

import (
	"context"
	"errors"
	"fmt"
	auth_service "log-pulse/app/auth/services"
	"log-pulse/models"
	"log-pulse/pkg/constant"
	"log-pulse/pkg/repository"
	"log-pulse/pkg/utils"
	"log-pulse/platform/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MAX_LOGIN_ATTEMPTS uint = 6
	BLOCKED_FOR        uint = 3
	loginAttempt       models.LoginAttempts
)

const (
	InvalidLoginCredentials string = "Invalid login credentials."
)

func LoginHandler(c *fiber.Ctx) error {
	payload := new(LoginRequest)

	if err := utils.ValidateRequest(c, payload); err != nil {
		return utils.ErrorResponse(c, fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	user, err := repository.NewUserRepository().FindByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, database.NOTFOUND) {
			return utils.ErrorResponse(c, fiber.Error{
				Code:    fiber.StatusUnauthorized,
				Message: InvalidLoginCredentials,
			})
		}
	}

	if err := loadLoginAttemptModel(user); err != nil {
		return err
	}

	if hasTooManyLoginAttempts(user) == true {
		return sendLockoutResponse(c)
	}

	if utils.CheckPasswordHash(payload.Password, user.Password) == false {
		return sendFailedLoginResponse(c, user, incrementLoginAttempts(user))
	}

	clearLoginAttempt(user)

	token, _ := auth_service.GenerateTokens(user)
	organization, _ := repository.NewOrganizationRepository().FindByUser(user.ID)

	return utils.SuccessResponse(c, "Login successful", fiber.Map{
		"user":    auth_service.AuthenticatedUser(user, organization),
		"token":   token.Access,
		"refresh": token.Refresh,
	})
}

func loadLoginAttemptModel(user *models.User) error {
	err := database.DB.Collection("login_attempts").
		FindOne(context.Background(), bson.D{{Key: "user_id", Value: user.ID}}).
		Decode(&loginAttempt)

	if err != nil {
		if errors.Is(err, database.NOTFOUND) == false {
			fmt.Println("Failed to load login attempt for ", user.ID, err)
			return constant.ErrSomethingWentWrong
		}
	}

	return nil
}

func clearLoginAttempt(user *models.User) {
	attempts := loginAttempt.Attempts

	if attempts != nil {
		_, err := database.DB.Collection("login_attempts").
			UpdateByID(context.Background(), loginAttempt.ID, bson.D{{Key: "attempts", Value: 0}})

		if err != nil {
			fmt.Println("Failed to update login attempt", err.Error())
		}
	}
}

func hasTooManyLoginAttempts(user *models.User) bool {
	if loginAttempt.LastAttemptTime == nil {
		return false
	}

	diff := getTimeDiffForAttempt(getLastLoginAttemptTime(loginAttempt))

	return (*loginAttempt.Attempts >= int(MAX_LOGIN_ATTEMPTS)) && (diff < int(BLOCKED_FOR))
}

func sendLockoutResponse(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.Error{
		Message: fmt.Sprintf(
			"Your account has been blocked, Please try again after %v minutes",
			getBlockedMinutesLeft(getLastLoginAttemptTime(loginAttempt)),
		),
		Code: fiber.StatusUnauthorized,
	})
}

func getLastLoginAttemptTime(loginAttempt models.LoginAttempts) time.Time {
	lastAttemptTime := time.Time{}
	if loginAttempt.LastAttemptTime != nil {
		lastAttemptTime = *loginAttempt.LastAttemptTime
	}

	return lastAttemptTime
}

func sendFailedLoginResponse(c *fiber.Ctx, user *models.User, attempts int) error {
	remainingAttempts := int(MAX_LOGIN_ATTEMPTS) - attempts
	if remainingAttempts == 0 && loginAttempt.Attempts != nil {
		return utils.ErrorResponse(c, fiber.Error{
			Message: fmt.Sprintf(
				"Your Account has been blocked, Please try again after %v minutes",
				getBlockedMinutesLeft(getLastLoginAttemptTime(loginAttempt)),
			),
			Code: fiber.StatusUnauthorized,
		})
	}

	if remainingAttempts < 3 {
		return utils.ErrorResponse(c, fiber.Error{
			Message: fmt.Sprintf(InvalidLoginCredentials+" %v Attempts left", remainingAttempts),
			Code:    fiber.StatusUnauthorized,
		})
	}

	return utils.ErrorResponse(c, fiber.Error{
		Message: InvalidLoginCredentials,
		Code:    fiber.StatusUnauthorized,
	})
}

func incrementLoginAttempts(user *models.User) int {
	var attempts *int = loginAttempt.Attempts
	if attempts == nil {
		attempts = new(int)
		*attempts = 1
	}

	if loginAttempt.LastAttemptTime != nil {
		diff := getTimeDiffForAttempt(getLastLoginAttemptTime(loginAttempt))
		if diff < int(BLOCKED_FOR) {
			*attempts = *attempts + 1
		}
	}

	_, err := database.DB.Collection("login_attempts").UpdateOne(
		context.Background(),
		bson.D{{Key: "user_id", Value: user.ID}},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "user_id", Value: user.ID},
				{Key: "attempts", Value: attempts},
				{Key: "last_attempt_time", Value: time.Now()},
			}}},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		fmt.Printf("Error updating login attempt", err.Error())
	}

	return *attempts
}

func getTimeDiffForAttempt(lastAttemptTime time.Time) int {
	currentTime := time.Now()
	diffMinutes := int(currentTime.Sub(lastAttemptTime).Minutes())

	return diffMinutes
}

func getBlockedMinutesLeft(lastFailedTime time.Time) int {
	diff := getTimeDiffForAttempt(lastFailedTime)

	return int(BLOCKED_FOR) - diff
}
