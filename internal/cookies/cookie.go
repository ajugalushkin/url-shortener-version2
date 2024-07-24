package cookies

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
)

// Claims структура для генерации токена.
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// TokenExp храним срок истечения токена.
const TokenExp = time.Hour * 3

// buildJWTString функция генерации токена
func buildJWTString(ctx context.Context) (string, error) {
	flags := config.GetConfig()

	rawUser, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return "", err
	}

	userID := int(rawUser.Int64())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(flags.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUser функция для получения пользователя из токена
func GetUser(ctx context.Context, tokenString string) *dto.User {
	flags := config.GetConfig()
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(flags.SecretKey), nil
	})
	if err != nil {
		return &dto.User{}
	}

	return &dto.User{ID: claims.UserID}
}

// createCookie функция для создания куки
func createCookie(ctx context.Context, nameCookie string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = nameCookie
	cookie.Value, _ = buildJWTString(ctx)
	cookie.Expires = time.Now().Add(TokenExp)
	return cookie
}

// Write функция записывает куки в контекст
func Write(ctx context.Context, echoCtx echo.Context, nameCookie string) string {
	cookie := createCookie(ctx, nameCookie)
	echoCtx.SetCookie(cookie)
	return cookie.Value
}

// Read функция читает куки из контекста
func Read(echoCtx echo.Context, name string) (string, error) {
	cookie, err := echoCtx.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
