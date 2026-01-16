package users

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/okira-e/go-as-your-backend/app/models"
	"github.com/okira-e/go-as-your-backend/app/utils"
)

// AuthMiddleware validates the user from the JWT token and adds user to context
func AuthMiddleware(usersService *Service) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		jwtSecret := []byte(utils.RequireEnv("JWT_SECRET"))

		// -------- Try access token --------
		if access := ctx.Cookies("access_token"); access != "" {
			claims, err := validateToken(access, jwtSecret)
			if err != nil {
				return utils.Err(ctx, 401, "Invalid access token.", nil)
			}

			jwtUser := models.JwtUser{
				UserID:   claims["userId"].(string),
				Email:    claims["email"].(string),
				RoleName: claims["roleName"].(string),
			}

			ctx.Locals("user", jwtUser)
			return ctx.Next()
		}

		// -------- Try refresh token --------
		refresh := ctx.Cookies("refresh_token")
		if refresh == "" {
			return utils.Err(ctx, 401, "Missing refresh token.", nil)
		}

		newAccessToken, status, err := usersService.Refresh(refresh)
		if err != nil || status != 200 {
			return utils.Err(ctx, 401, "Session expired. "+err.Error(), nil)
		}

		// -------- Set new access token cookie --------
		accessTokenExpiryStr := utils.RequireEnv("ACCESS_TOKEN_EXPIRY")
		accessTokenExpiry, _ := strconv.Atoi(accessTokenExpiryStr)

		ctx.Cookie(&fiber.Cookie{
			Name:   "access_token",
			Value:  newAccessToken,
			MaxAge: accessTokenExpiry,
			Path:   "/",
			// Domain:      domain,
			Secure:      utils.RequireEnv("ENV") == "prod",
			HTTPOnly:    true,
			SameSite:    "Lax",
			SessionOnly: false,
		})

		// -------- Decode refreshed token --------
		claims, err := validateToken(newAccessToken, jwtSecret)
		if err != nil {
			return utils.Err(ctx, 401, "Invalid refreshed token.", err)
		}

		// -------- Continue with the user claims --------
		jwtUser := models.JwtUser{
			UserID:   claims["userId"].(string),
			Email:    claims["email"].(string),
			RoleName: claims["roleName"].(string),
		}

		ctx.Locals("user", jwtUser)
		return ctx.Next()
	}
}

// RoleMiddleware validates that the user has the required role
// It takes the user from the Fiber Ctx. So call this after AuthMiddleware
func RoleMiddleware(requiredRole string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user := ctx.Locals("user")
		if user == nil {
			return utils.Err(ctx, 401, "User not found in context", nil)
		}

		jwtUser, ok := user.(models.JwtUser)
		if !ok {
			return utils.Err(ctx, 500, "Invalid user type in context", nil)
		}

		if jwtUser.RoleName != requiredRole {
			return utils.Err(ctx, 403, "Required role is missing", nil)
		}

		return ctx.Next()
	}
}
