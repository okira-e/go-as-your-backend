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

	// Check for fields in the struct that are required and not provided.
	err = utils.ValidateFields(&body)
	if err != nil {
		return utils.Err(c, 400, err.Error(), nil)
	}

	// Retrieve the roles provided in the request body from the database from the roles table.
	// `Find` automatically knows to look in the `roles` table because the type of `roles` is `[]*tables.Roles` matches
	// the type of the `roles` table in `models.go`
	var roles []*tables.Roles
	results := db.Where("id IN ?", body.Roles).Find(&roles)
	if results.Error != nil {
		return utils.Err(c, 500, results.Error.Error(), nil)
	}

	// Check if the roles returned are the same as the roles requested. If not, return an error.
	// This ensures all the roles are valid.
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

	// Remove extra fields on the request body and move the data to a SystemUsers struct object.
	var systemUser tables.SystemUsers
	err = copier.Copy(&systemUser, &body) // Copies the values of body to systemUser.
	if err != nil {
		return utils.Err(c, 500, err.Error(), nil)
	}

	// Set the roles of the system user to the roles retrieved from the database.
	// Gorm will automatically create the relationship between the system user and the roles by inserting the roles
	// IDs into the `users_roles` table.
	systemUser.Roles = roles

	// Create the user in the database. Again, Gorm knows to insert in the `system_users` table because the type of
	// `systemUser` is `tables.SystemUsers` which matches the type of the `system_users` table in `models.go`.
	result := db.Create(&systemUser)
	if result.Error != nil {
		return utils.Err(c, 500, result.Error.Error(), nil)
	}

	userId := systemUser.ID // Get the ID of the system user that was just created. Feels like magic.

	return utils.Ok(c, 201, "", userId)
}

