package posts

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	. "github.com/okira-e/go-as-your-backend/app/logging"
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
		Log(SeverityError, "FindAll: failed to fetch posts", map[string]any{"error": err.Error()})
		return utils.Err(ctx, 500, "Failed to fetch entities", nil)
	}

	entitiesDto := make([]*models.PostDto, len(entities))
	for i, entity := range entities {
		entitiesDto[i] = entity.ToDto()
	}

	return utils.Ok(ctx, 200, "", entitiesDto)
}

func (self *Handler) GetPublished(ctx *fiber.Ctx) error {
	entities, err := self.service.GetPublished(ctx.Context())
	if err != nil {
		Log(SeverityError, "GetPublished: failed to fetch posts", map[string]any{"error": err.Error()})
		return utils.Err(ctx, 500, "Failed to fetch entities", nil)
	}

	entitiesDto := make([]*models.PostDto, len(entities))
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
		Log(SeverityError, "GetCount: failed to fetch count", map[string]any{"error": err.Error()})
		return utils.Err(ctx, 500, "Failed to fetch entities count", nil)
	}

	return utils.Ok(ctx, 200, "", count)
}

func (self *Handler) Create(ctx *fiber.Ctx) error {
	entityDto := models.CreatePostDto{}

	if err := ctx.BodyParser(&entityDto); err != nil {
		return utils.Err(ctx, 400, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(entityDto); err != nil {
		Log(SeverityWarn, "Create: validation failed", map[string]any{"error": err.Error()})
		return utils.Err(ctx, 400, "Validation failed", err.Error())
	}

	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		return utils.Err(ctx, 401, "Missing user in headers", nil)
	}

	post, err := self.service.Create(ctx.Context(), &entityDto, user.UserID)
	if err != nil {
		Log(SeverityError, "Create: failed to create post", map[string]any{"userId": user.UserID, "error": err.Error()})
		return utils.Err(ctx, 500, "Failed to create entity", err.Error())
	}

	Log(SeverityInfo, "Create: post created", map[string]any{"postId": post.ID, "userId": user.UserID})

	return utils.Ok(ctx, 201, "", post)
}
