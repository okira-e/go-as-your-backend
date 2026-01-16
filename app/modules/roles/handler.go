package roles

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/okira-e/go-as-your-backend/app/models"
	"github.com/okira-e/go-as-your-backend/app/spec"
	"github.com/okira-e/go-as-your-backend/app/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(Service *Service) *Handler {
	return &Handler{service: Service}
}

func (self *Handler) FindAll(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit", "100"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))

	queryOptions := spec.QueryOptions{
		Limit:  limit,
		Offset: offset,
	}

	filter, err := spec.ParseFilter(ctx.Query("filter", ""))
	if err != nil {
		return utils.Err(ctx, 400, "Invalid filter parameter", err)
	}

	entities, err := self.service.FindAll(ctx.Context(), &queryOptions, filter)
	if err != nil {
		return utils.Err(ctx, 500, "Failed to fetch entities", nil)
	}

	entitiesDto := make([]*models.RoleDto, len(entities))
	for i, entity := range entities {
		entitiesDto[i] = entity.ToDto()
	}

	return utils.Ok(ctx, 200, "", entitiesDto)
}

func (self *Handler) GetCount(ctx *fiber.Ctx) error {
	filter, err := spec.ParseFilter(ctx.Query("filter", ""))
	if err != nil {
		return utils.Err(ctx, 400, "Invalid filter parameter", err)
	}

	count, err := self.service.GetCount(ctx.Context(), filter)
	if err != nil {
		return utils.Err(ctx, 500, "Failed to fetch entities count", nil)
	}

	return utils.Ok(ctx, 200, "", count)
}

func (self *Handler) Create(ctx *fiber.Ctx) error {
	var entityDto models.CreateRoleDto

	if err := ctx.BodyParser(&entityDto); err != nil {
		return utils.Err(ctx, 400, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(entityDto); err != nil {
		return utils.Err(ctx, 400, "Validation failed", err.Error())
	}

	id, err := self.service.Create(ctx.Context(), &entityDto)
	if err != nil {
		return utils.Err(ctx, 500, "Failed to create entity", err)
	}

	return utils.Ok(ctx, 201, "", id)
}
