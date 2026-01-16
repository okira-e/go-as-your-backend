package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/okira-e/go-as-your-backend/app/models"
	"github.com/okira-e/go-as-your-backend/app/spec"
	"github.com/okira-e/go-as-your-backend/app/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository spec.Repository[models.User]
}

func NewService(repository spec.Repository[models.User]) *Service {
	return &Service{repository: repository}
}

func (self *Service) FindAll(ctx context.Context, queryOptions *spec.QueryOptions, filter *spec.Filter) ([]models.User, error) {
	entities, err := self.repository.FindAll(ctx, queryOptions, filter)
	if err != nil {
		return entities, err
	}

	return entities, nil
}

func (self *Service) GetContactInfo(ctx context.Context, userId string) (models.UserContact, error) {
	filter := spec.Filter{
		Where: spec.WhereClause{
			And: []spec.WhereCondition{
				{
					Column:   "id",
					Operator: "=",
					Value:    userId,
				},
			},
		},
	}
	entities, err := self.repository.FindAll(ctx, nil, &filter)
	if err != nil {
		return models.UserContact{}, err
	}
	
	if len(entities) == 0 {
		return models.UserContact{}, fmt.Errorf("User with ID %s not found.", userId)
	}
	
	info := models.UserContact{
		Phone: entities[0].Phone,
	}

	return info, nil
}

func (self *Service) GetCount(ctx context.Context, filter *spec.Filter) (int64, error) {
	count, err := self.repository.Count(ctx, filter)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (self *Service) Create(
	ctx context.Context,
	entityDto *models.UserDto,
	password string,
) (*models.User, error) {
	filter := spec.Filter{
		Where: spec.WhereClause{
			And: []spec.WhereCondition{
				{
					Column:   "email",
					Operator: "=",
					Value:    entityDto.Email,
				},
			},
		},
	}
	existingUsers, err := self.repository.FindAll(ctx, nil, &filter)
	if err != nil {
		return &models.User{}, fmt.Errorf("Couldn't fetch existing users for validation. %s", err)
	}

	if len(existingUsers) > 0 {
		return &models.User{}, fmt.Errorf("User with email %s already exists.", entityDto.Email)
	}

	if _, err = mail.ParseAddress(entityDto.Email); err != nil {
		return &models.User{}, fmt.Errorf("Invalid email. %s", err)
	}

	if err = ValidateUserPassword(password); err != nil {
		return &models.User{}, fmt.Errorf("Invalid password. %s", err)
	}

	hashedPassBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return &models.User{}, fmt.Errorf("Encountered an error while hashing the password. %s", err)
	}

	entity := entityDto.FromDto(true, string(hashedPassBytes))

	entity.IsActive = true

	entity, err = self.repository.Create(ctx, entity)
	if err != nil {
		return &models.User{}, err
	}

	return entity, nil
}

func ValidateUserPassword(password string) error {
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters long.")
	}

	if !strings.ContainsAny(password, "0123456789") {
		return errors.New("Password must contain at least one number.")
	}

	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return errors.New("Password must contain at least one uppercase letter.")
	}

	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return errors.New("Password must contain at least one lowercase letter.")
	}

	if !strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?") {
		return errors.New("Password must contain at least one special character.")
	}

	return nil
}

func (self *Service) Login(email string, pass string) (string, models.User, int, error) {
	user, err := self.ValidateSystemUserCredentials(email, pass)
	if err != nil {
		return "", models.User{}, 400, fmt.Errorf("User credentials were invalid.")
	}

	refreshToken, err := self.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", models.User{}, 500, fmt.Errorf("Error while generating a refresh token. %s\n", err.Error())
	}

	return refreshToken, user, 200, nil
}

func (self *Service) GenerateAccessToken(userId string, email string, roleName string) (string, error) {
	accessTokenExpiryStr := utils.RequireEnv("ACCESS_TOKEN_EXPIRY")
	accessTokenExpiry, err := strconv.Atoi(accessTokenExpiryStr)
	if err != nil {
		log.Fatalf("Invalid ACCESS_TOKEN_EXPIRY value: %s\n", err.Error())
	}

	claims := jwt.MapClaims{
		"userId":   userId,
		"email":    email,
		"roleName": roleName,
		"sub":      email,                   // Subject claim (typically user ID)
		"iss":      "go-as-your-backend",                 // Issuer claim
		"aud":      "https://api.go-as-your-backend.com", // Audience claim
		"exp":      time.Now().Add(time.Duration(accessTokenExpiry) * time.Second).Unix(),
		"iat":      time.Now().Unix(), // Issued At
		"nbf":      time.Now().Unix(), // Not Before
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := utils.RequireEnv("JWT_SECRET")

	return token.SignedString([]byte(jwtSecret))
}

func (self *Service) GenerateRefreshToken(userId string) (string, error) {
	refreshTokenExpiryStr := utils.RequireEnv("REFRESH_TOKEN_EXPIRY")
	refreshTokenExpiry, err := strconv.Atoi(refreshTokenExpiryStr)
	if err != nil {
		log.Fatalf("Invalid REFRESH_TOKEN_EXPIRY value: %s\n", err.Error())
	}

	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Duration(refreshTokenExpiry) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := utils.RequireEnv("JWT_SECRET")

	return token.SignedString([]byte(jwtSecret))
}

// validateToken validates a token and returns claims.
func validateToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	prefix := "Bearer "
	tokenString = strings.TrimPrefix(tokenString, prefix)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method.")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token.")
}

// Refresh generates a new access token using the refresh token.
func (self *Service) Refresh(refreshToken string) (string, int, error) {
	jwtSecret := utils.RequireEnv("JWT_SECRET")

	claims, err := validateToken(refreshToken, []byte(jwtSecret))
	if err != nil {
		return "", 401, fmt.Errorf("Invalid token. %s\n", err.Error())
	}

	userId, _ := claims["sub"].(string)
	if userId == "" {
		return "", 401, fmt.Errorf("Invalid refresh token: missing subject")
	}

	// Fetch the latest user info from DB
	filter := spec.Filter{
		Where: spec.WhereClause{
			And: []spec.WhereCondition{
				{
					Column:   "id",
					Operator: "=",
					Value:    userId,
				},
			},
		},
	}
	results, err := self.FindAll(context.Background(), nil, &filter)
	if err != nil || len(results) == 0 {
		return "", 500, fmt.Errorf("User not found.")
	}
	user := results[0]

	accessToken, err := self.GenerateAccessToken(user.ID, user.Email, user.Role.Name)
	if err != nil {
		return "", 500, fmt.Errorf("Error while generating an access token. %s\n", err.Error())
	}

	return accessToken, 200, nil
}

// ValidateSystemUserCredentials validates the username and password of a user.
func (self *Service) ValidateSystemUserCredentials(email string, password string) (models.User, error) {
	// @Todo: Simplify
	filter := spec.Filter{
		Where: spec.WhereClause{
			And: []spec.WhereCondition{
				{
					Column:   "email",
					Operator: "=",
					Value:    email,
				},
			},
		},
	}
	users, err := self.FindAll(context.Background(), nil, &filter)
	if err != nil {
		return models.User{}, fmt.Errorf("Error fetching user. %s", err)
	}

	if len(users) == 0 {
		return models.User{}, fmt.Errorf("User not found.")
	}

	user := users[0]

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, fmt.Errorf("Invalid credentials. %s", err)
	}

	return user, nil
}

func (self *Service) ValidateUserFromRequest(requestHeaders map[string][]string) (models.User, int, error) {
	jwtSecret := utils.RequireEnv("JWT_SECRET")

	authorizationValues := requestHeaders["Authorization"]
	if len(authorizationValues) == 0 {
		return models.User{}, 401, fmt.Errorf("No token was found in the headers.")
	}

	token := authorizationValues[0]

	claims, err := validateToken(token, []byte(jwtSecret))
	if err != nil {
		return models.User{}, 401, fmt.Errorf("Invalid token. %s\n", err.Error())
	}

	email := claims["email"].(string)
	if email == "" {
		return models.User{}, 401, fmt.Errorf("\"email\" claim is missing in the token.")
	}

	filter := spec.Filter{
		Where: spec.WhereClause{
			And: []spec.WhereCondition{
				{
					Column:   "email",
					Operator: "=",
					Value:    email,
				},
			},
		},
	}
	users, err := self.FindAll(context.Background(), nil, &filter)
	if err != nil {
		return models.User{}, 500, fmt.Errorf("Error fetching user. %s", err)
	}

	if len(users) < 1 {
		return models.User{}, 401, fmt.Errorf("User doesn't exist.")
	}

	return users[0], 0, nil
}
