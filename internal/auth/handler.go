package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

type ReqRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RespError struct {
	Error string `json:"error"`
}

type RespToken struct {
	Token string `json:"token"`
}

func Register(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ReqRegister
		if json.NewDecoder(r.Body).Decode(&req) != nil || req.Username == "" || req.Password == "" {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err := store.Create(&User{Username: req.Username, Password: hash}); err != nil {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func Login(store *MemoryStore, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ReqRegister
		if json.NewDecoder(r.Body).Decode(&req) != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		u, ok := store.Get(req.Username)
		if !ok || bcrypt.CompareHashAndPassword(u.Password, []byte(req.Password)) != nil {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(ErrWrongCreds)
			return
		}
		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
			"sub": req.Username,
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		})
		_ = json.NewEncoder(w).Encode(RespToken{Token: tokenString})
	}
}
