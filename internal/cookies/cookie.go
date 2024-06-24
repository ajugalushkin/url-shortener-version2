package cookies

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/logger"
)

type Claims struct {
	jwt.RegisteredClaims
	User dto.User
}

const TokenExp = time.Hour * 3

func buildJWTString(ctx context.Context) (string, error) {
	flags := config.FlagsFromContext(ctx)

	rawUser, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return "", err
	}

	user := dto.User{
		UserID: int(rawUser.Int64())}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		User: user,
	})

	tokenString, err := token.SignedString([]byte(flags.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//func GetUserID(ctx context.Context, tokenString string) int {
//	flags := config.FlagsFromContext(ctx)
//	claims := &Claims{}
//	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
//		return []byte(flags.SecretKey), nil
//	})
//	if err != nil {
//		return 0
//	}
//
//	return claims.User
//}

//func createCookie(ctx context.Context, nameCookie string) *http.Cookie {
//	cookie := new(http.Cookie)
//	cookie.Name = nameCookie
//	cookie.Value, _ = buildJWTString(ctx)
//	cookie.Expires = time.Now().Add(TokenExp)
//	return cookie
//}
//
//func Write(ctx context.Context, echoCtx echo.Context, nameCookie string) string {
//	cookie := createCookie(ctx, nameCookie)
//	echoCtx.SetCookie(cookie)
//	return cookie.Value
//}

//func Read(echoCtx echo.Context, name string) (string, error) {
//	cookie, err := echoCtx.Cookie(name)
//	if err != nil {
//		return "", err
//	}
//	return cookie.Value, nil
//}

func Create(ctx context.Context, nameCookie string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = nameCookie
	cookie.Value, _ = buildJWTString(ctx)
	cookie.Expires = time.Now().Add(TokenExp)
	return cookie
}

func GetUser(ctx context.Context, tokenString string) *dto.User {
	cfg := config.FlagsFromContext(ctx)
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.SecretKey), nil
	})
	if err != nil {
		logger.GetSingleton(ctx).GetLogger().Debug("cookies.GetLogin: unable to parse token")
		return &dto.User{}
	}

	return &claims.User
}
