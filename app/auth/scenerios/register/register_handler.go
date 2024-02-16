package register_scenerio

import (
	"errors"
	"fmt"
	auth_service "log-pulse/app/auth/services"
	"log-pulse/models"
	"log-pulse/pkg/constant"
	"log-pulse/pkg/repository"
	"log-pulse/pkg/utils"
	"log-pulse/platform/database"

	"github.com/gofiber/fiber/v2"
)

func RegisterHandler(c *fiber.Ctx) error {
	request := new(RegisterRequest)

	if err := utils.ValidateRequest(c, request); err != nil {
		return utils.ErrorResponse(c, fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := validateUserWithEmail(request.Email); err != nil {
		return utils.ErrorResponse(c, fiber.Error{
			Code:    fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	user, err := createAccount(request)

	if err != nil {
		return utils.ErrorResponse(c, fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	go func(user *models.User, organization string) {
		err := createOrganization(user, organization)
		if err != nil {
			fmt.Println("Error creating organization:", err)
		}

	}(user, request.Organization)

	token, err := auth_service.GenerateTokens(user)
	if err != nil {
		return constant.UnableToGenerateToken
	}

	organization, _ := repository.NewOrganizationRepository().FindByUser(user.ID)

	return utils.SuccessResponse(c, "Register successfully", fiber.Map{
		"user":    auth_service.AuthenticatedUser(user, organization),
		"token":   token.Access,
		"refresh": token.Refresh,
	})
}

func validateUserWithEmail(email string) error {
	userModel, err := repository.NewUserRepository().FindByEmail(email)
	if err != nil {
		if errors.Is(err, database.NOTFOUND) != true {
			return err
		}
	}

	if userModel != nil {
		return errors.New("Email address is already taken")
	}

	return nil
}

func createAccount(request *RegisterRequest) (*models.User, error) {
	result, err := repository.NewUserRepository().CreateOne(&models.User{
		FullName: request.FullName,
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return result, nil
}

func createOrganization(user *models.User, organization string) error {
	_, err := repository.NewOrganizationRepository().CreateOne(&models.Organization{
		UserID: user.ID,
		Name:   organization,
	})

	if err != nil {
		return err
	}

	return nil
}
