package auth

import (
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

var TokenAuth = jwtauth.New("HS256", []byte("super-secret-key"), nil)

func Middleware(tokenAuth *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return jwtauth.Verifier(tokenAuth)(jwtauth.Authenticator(next))
	}
}
