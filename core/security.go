package core

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/BimaAdi/fiberGormBoilerplate/models"
	"github.com/BimaAdi/fiberGormBoilerplate/settings"
	"github.com/gofiber/fiber/v2"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWTToken(user_id string, user_email string) (string, error) {
	// Generate Payload
	expiredAt := time.Now().Add(time.Minute * time.Duration(settings.ACCESS_TOKEN_EXPIRE_MINUTES))
	tok, err := jwt.NewBuilder().
		IssuedAt(time.Now()).
		Expiration(expiredAt).
		Build()
	tok.Set("id", user_id)
	tok.Set("email", user_email)
	if err != nil {
		return "", err
	}

	// Sign a JWT!
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS256, []byte(settings.JWT_SECRET)))
	if err != nil {
		return "", err
	}
	return string(signed[:]), nil

}

func GetPayloadFromJWTToken(jwtToken string) (string, string, error) {
	tok, err := jwt.Parse([]byte(jwtToken), jwt.WithKey(jwa.HS256, []byte(settings.JWT_SECRET)), jwt.WithValidate(false))
	if err != nil {
		return "", "", err
	}

	// Validate token
	err = jwt.Validate(tok)
	if err != nil {
		return "", "", err
	}

	// Get Payload
	id, isIdFound := tok.Get("id")
	if !isIdFound {
		return "", "", errors.New("id not found on token payload")
	}
	email, isEmailFound := tok.Get("email")
	if !isEmailFound {
		return "", "", errors.New("email not found on token payload")
	}

	return fmt.Sprint(id), fmt.Sprint(email), nil
}

func GenerateJWTTokenFromUser(tx *gorm.DB, user models.User) (string, error) {
	tok, err := GenerateJWTToken(user.ID, user.Email)
	return tok, err
}

func GetUserFromJWTToken(tx *gorm.DB, jwtToken string) (models.User, error) {
	user := models.User{}
	userId, _, err := GetPayloadFromJWTToken(jwtToken)
	if err != nil {
		return user, err
	}

	if err := models.DBConn.Where("id = ? AND deleted_at IS NULL", userId).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

type Header struct {
	Authorization string `regHeader:"authorization"`
}

func GetUserFromAuthorizationHeader(tx *gorm.DB, c *fiber.Ctx) (models.User, error) {
	header := new(Header)

	if err := c.ReqHeaderParser(header); err != nil {
		return models.User{}, err
	}
	authHeader := header.Authorization

	arrayHeader := strings.Fields(authHeader)
	if len(arrayHeader) != 2 {
		return models.User{}, errors.New("invalid token key lenght no 2")
	}

	key := arrayHeader[0]
	token := arrayHeader[1]
	if key != "Bearer" {
		return models.User{}, errors.New("invalid token key not Bearer")
	}

	user, err := GetUserFromJWTToken(tx, token)
	if err != nil {
		return models.User{}, errors.New("invalid token")
	}

	return user, nil
}
