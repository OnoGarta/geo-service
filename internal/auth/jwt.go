package auth

import (
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

// Middleware возвращает HTTP-мидлвар для защиты маршрутов с помощью JWT-аутентификации
func Middleware(tokenAuth *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return jwtauth.Verifier(tokenAuth)(jwtauth.Authenticator(next))
	}
}
