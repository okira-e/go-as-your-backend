package routes

import (
	"encoding/json"
	"fmt"
	"github.com/okira-e/go-as-your-backend/app/tables"
	"github.com/okira-e/go-as-your-backend/app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

func systemUsersRouter(app *fiber.App) {
	routerPath := "/api/system-users"

	app.Get(routerPath, getSystemUsers)
	app.Post(routerPath, registerSystemUser)
}

type SystemUserDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Title     string `json:"title"`
}

// @Summary Get all system users DTO
// @Description Get all system users DTO
// @Tags system-users
// @Accept */*
// @Produce application/json
// @Success 200 "OK"
// @Router /api/system-users [GET]
func getSystemUsers(c *fiber.Ctx) error {
	var systemUsers []tables.SystemUsers
	result := db.Find(&systemUsers)
	if result.Error != nil {
		return utils.Err(c, 500, result.Error.Error(), nil)
	}

	// Copy the values of systemUsers structs to systemUsersDTO structs.
	systemUsersDTO := []SystemUserDTO{}
	err := copier.Copy(&systemUsersDTO, &systemUsers)
	if err != nil {
		return utils.Err(c, 500, err.Error(), nil)
	}

	return utils.Ok(c, 200, "", systemUsersDTO)
}

// @Summary Register a system user
// @Description Register a system user for the first time
// @Tags system-users
// @Accept */*
// @Produce application/json
// @Param firstName body string true "First name of the system user"
func registerSystemUser(c *fiber.Ctx) error {
	body := struct {
		tables.SystemUsers

		Roles []uint `validate:"min=1" json:"roles"` // min takes the place of required as well.
	}{}
	err := json.Unmarshal(c.Body(), &body)
	if err != nil {
		return utils.Err(c, 400, err.Error(), nil)
	}

	err = utils.ValidateFields(&body)
	if err != nil {
		return utils.Err(c, 400, err.Error(), nil)
	}

	var roles []*tables.Roles
	results := db.Where("id IN ?", body.Roles).Find(&roles)
	if results.Error != nil {
		return utils.Err(c, 500, results.Error.Error(), nil)
	}

	// Check if the roles returned are the same as the roles requested. If not, return an error.
	inValidRoles := []uint{}
	for _, requestedRole := range body.Roles {
		found := false
		for _, returnedRole := range roles {
			if requestedRole == returnedRole.ID {
				found = true
				break
			}
		}
		if !found {
			inValidRoles = append(inValidRoles, requestedRole)
		}
	}
	if len(inValidRoles) > 0 {
		return utils.Err(c, 400, fmt.Sprintf("Invalid roles: %d", inValidRoles), nil)
	}

	var systemUser tables.SystemUsers
	err = copier.Copy(&systemUser, &body)
	if err != nil {
		return utils.Err(c, 500, err.Error(), nil)
	}

	systemUser.Roles = roles

	result := db.Create(&systemUser)
	if result.Error != nil {
		return utils.Err(c, 500, result.Error.Error(), nil)
	}

	userId := systemUser.ID // Get the ID of the system user that was just created. Feels like magic.

	return utils.Ok(c, 201, "", userId)
}

