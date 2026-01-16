package users

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
		Log(SeverityError, "FindAll: failed to fetch users", map[string]any{"error": err.Error()})
		return utils.Err(ctx, 500, "Failed to fetch entities", nil)
	}

	entitiesDto := make([]*models.UserDto, len(entities))
	for i, entity := range entities {
		entitiesDto[i] = entity.ToDto()
	}

	return utils.Ok(ctx, 200, "", entitiesDto)
}

func (self *Handler) GetContactInfo(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	if userId == "" {
		return utils.Err(ctx, 400, "User ID not provided", nil)
	}

	info, err := self.service.GetContactInfo(ctx.Context(), userId)
	if err != nil {
		Log(SeverityError, "GetContactInfo: failed to fetch", map[string]any{"userId": userId, "error": err.Error()})
		return utils.Err(ctx, 500, "Failed to fetch contact info", nil)
	}

	return utils.Ok(ctx, 200, "", info)
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

func (self *Handler) Login(ctx *fiber.Ctx) error {
	var payload struct {
		Email    string `json:"email"       validate:"required,email"`
		Password string `json:"password"    validate:"required"`
	}

	if err := ctx.BodyParser(&payload); err != nil {
		return utils.Err(ctx, 400, "Invalid payload body", err.Error())
	}

	if err := utils.ValidateStruct(payload); err != nil {
		return utils.Err(ctx, 400, "Validation failed", err.Error())
	}

	refreshToken, user, statusCode, err := self.service.Login(payload.Email, payload.Password)
	if err != nil {
		Log(SeverityWarn, "Login: failed attempt", map[string]any{"email": payload.Email})
		return utils.Err(ctx, statusCode, err.Error(), nil)
	}

	accessToken, err := self.service.GenerateAccessToken(user.ID, user.Email, user.Role.Name)
	if err != nil {
		Log(SeverityError, "Login: failed to generate access token", map[string]any{"userId": user.ID, "error": err.Error()})
		return utils.Err(ctx, 500, "Failed to generate access token", nil)
	}

	domain := utils.RequireEnv("DOMAIN")

	accessTokenExpiryStr := utils.RequireEnv("ACCESS_TOKEN_EXPIRY")
	accessTokenExpiry, _ := strconv.Atoi(accessTokenExpiryStr)

	refreshTokenExpiryStr := utils.RequireEnv("REFRESH_TOKEN_EXPIRY")
	refreshTokenExpiry, _ := strconv.Atoi(refreshTokenExpiryStr)

	sameSite := "Lax"
	if utils.RequireEnv("ENV") == "prod" {
		sameSite = "None"
	}

	ctx.Cookie(&fiber.Cookie{
		Name:        "refresh_token",
		Value:       refreshToken,
		MaxAge:      refreshTokenExpiry,
		Path:        "/",
		Domain:      domain,
		Secure:      utils.RequireEnv("ENV") == "prod",
		HTTPOnly:    true,
		SameSite:    sameSite,
		SessionOnly: false,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:        "access_token",
		Value:       accessToken,
		MaxAge:      accessTokenExpiry,
		Path:        "/",
		Domain:      domain,
		Secure:      utils.RequireEnv("ENV") == "prod",
		HTTPOnly:    true,
		SameSite:    sameSite,
		SessionOnly: false,
	})

	Log(SeverityInfo, "Login: successful", map[string]any{"userId": user.ID})

	return utils.Ok(ctx, statusCode, "Login successful", user.ToDto())
}

func (self *Handler) Me(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals("user").(models.JwtUser)
	if !ok {
		return utils.Err(ctx, 401, "Invalid user claims", nil)
	}

	filter := spec.Filter{
		Where: spec.WhereClause{
			And: []spec.WhereCondition{
				{
					Column:   "email",
					Operator: "=",
					Value:    user.Email,
				},
			},
		},
	}

	users, err := self.service.FindAll(ctx.Context(), nil, &filter)
	if err != nil {
		return utils.Err(ctx, 401, "Failed to fetch user from claims", nil)
	}

	if len(users) == 0 {
		return utils.Err(ctx, 401, "User not found from claims", nil)
	}

	return utils.Ok(ctx, 200, "", users[0].ToDto())
}

func (self *Handler) Register(ctx *fiber.Ctx) error {
	var payload struct {
		models.UserDto
		Password string `json:"password" validate:"required"`
	}

	if err := ctx.BodyParser(&payload); err != nil {
		return utils.Err(ctx, 400, "Invalid payload body", err.Error())
	}

	if err := utils.ValidateStruct(payload); err != nil {
		return utils.Err(ctx, 400, "Validation failed", err.Error())
	}

	user, err := self.service.Create(ctx.Context(), &payload.UserDto, payload.Password)
	if err != nil {
		Log(SeverityWarn, "Register: failed to create user", map[string]any{"email": payload.Email, "error": err.Error()})
		return utils.Err(ctx, 400, err.Error(), nil)
	}

	Log(SeverityInfo, "Register: user created", map[string]any{"userId": user.ID, "email": user.Email})

	refreshToken, _, statusCode, err := self.service.Login(payload.Email, payload.Password)
	if err != nil {
		Log(SeverityError, "Register: auto-login failed", map[string]any{"userId": user.ID, "error": err.Error()})
		return utils.Err(ctx, statusCode, err.Error(), nil)
	}

	accessToken, err := self.service.GenerateAccessToken(user.ID, user.Email, user.Role.Name)
	if err != nil {
		Log(SeverityError, "Register: failed to generate access token", map[string]any{"userId": user.ID, "error": err.Error()})
		return utils.Err(ctx, 500, "Failed to generate access token", nil)
	}

	domain := utils.RequireEnv("DOMAIN")

	accessTokenExpiryStr := utils.RequireEnv("ACCESS_TOKEN_EXPIRY")
	accessTokenExpiry, _ := strconv.Atoi(accessTokenExpiryStr)

	refreshTokenExpiryStr := utils.RequireEnv("REFRESH_TOKEN_EXPIRY")
	refreshTokenExpiry, _ := strconv.Atoi(refreshTokenExpiryStr)

	sameSite := "Lax"
	if utils.RequireEnv("ENV") == "prod" {
		sameSite = "None"
	}

	ctx.Cookie(&fiber.Cookie{
		Name:        "refresh_token",
		Value:       refreshToken,
		MaxAge:      refreshTokenExpiry,
		Path:        "/",
		Domain:      domain,
		Secure:      utils.RequireEnv("ENV") == "prod",
		HTTPOnly:    true,
		SameSite:    sameSite,
		SessionOnly: false,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:        "access_token",
		Value:       accessToken,
		MaxAge:      accessTokenExpiry,
		Path:        "/",
		Domain:      domain,
		Secure:      utils.RequireEnv("ENV") == "prod",
		HTTPOnly:    true,
		SameSite:    sameSite,
		SessionOnly: false,
	})

	return utils.Ok(ctx, 201, "User created successfully", user.ToDto())
}

func (self *Handler) RefreshToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return utils.Err(ctx, 401, "Refresh token not found", nil)
	}

	accessToken, statusCode, err := self.service.Refresh(refreshToken)
	if err != nil {
		Log(SeverityWarn, "RefreshToken: failed to refresh", map[string]any{"error": err.Error()})
		return utils.Err(ctx, statusCode, err.Error(), nil)
	}

	domain := utils.RequireEnv("DOMAIN")

	accessTokenExpiryStr := utils.RequireEnv("ACCESS_TOKEN_EXPIRY")
	accessTokenExpiry, _ := strconv.Atoi(accessTokenExpiryStr)

	sameSite := "Lax"
	if utils.RequireEnv("ENV") == "prod" {
		sameSite = "None"
	}

	ctx.Cookie(&fiber.Cookie{
		Name:        "access_token",
		Value:       accessToken,
		MaxAge:      accessTokenExpiry,
		Path:        "/",
		Domain:      domain,
		Secure:      utils.RequireEnv("ENV") == "prod",
		HTTPOnly:    true,
		SameSite:    sameSite,
		SessionOnly: false,
	})

	return utils.Ok(ctx, statusCode, "Token refreshed successfully", fiber.Map{})
}

func (self *Handler) Logout(ctx *fiber.Ctx) error {
	domain := utils.RequireEnv("DOMAIN")

	sameSite := "Lax"
	if utils.RequireEnv("ENV") == "prod" {
		sameSite = "None"
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HTTPOnly: true,
		Domain:   domain,
		Secure:   utils.RequireEnv("ENV") == "prod",
		SameSite: sameSite,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HTTPOnly: true,
		Domain:   domain,
		Secure:   utils.RequireEnv("ENV") == "prod",
		SameSite: sameSite,
	})

	return utils.Ok(ctx, 200, "Logout successful", nil)
}

func (self *Handler) ValidateToken(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	if token == "" {
		return utils.Err(ctx, 400, "Token is required", nil)
	}

	if err := utils.ValidateVar(token, "required"); err != nil {
		return utils.Err(ctx, 400, "Validation failed", err.Error())
	}

	jwtSecret := utils.RequireEnv("JWT_SECRET")
	if jwtSecret == "" {
		Log(SeverityError, "ValidateToken: JWT_SECRET not configured", nil)
		return utils.Err(ctx, 500, "JWT secret not configured", nil)
	}

	claims, err := validateToken(token, []byte(jwtSecret))
	if err != nil {
		return utils.Err(ctx, 401, "Invalid token", nil)
	}

	return utils.Ok(ctx, 200, "Token is valid", fiber.Map{
		"claims": claims,
	})
}
