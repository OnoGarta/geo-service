package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/your-org/geo-service-swagger/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

// AuthController отвечает за регистрацию и авторизацию пользователей
type AuthController struct {
	store     *auth.MemoryStore
	tokenAuth *jwtauth.JWTAuth
}

// новый тип для успешного ответа с JWT-токеном
type tokenResponse struct {
	Token string `json:"token"`
}

// NewAuthController создает контроллер аутентификации с заданным хранилищем пользователей и JWT-менеджером
func NewAuthController(store *auth.MemoryStore, tokenAuth *jwtauth.JWTAuth) *AuthController {
	return &AuthController{
		store:     store,
		tokenAuth: tokenAuth,
	}
}

// Register обрабатывает регистрацию нового пользователя (эндпоинт /api/register)
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	responder := NewResponder(w)
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// декодируем JSON-запрос
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder.Error(http.StatusBadRequest, "invalid json")
		return
	}
	// проверяем, что имя пользователя и пароль не пустые
	if req.Username == "" || req.Password == "" {
		responder.Error(http.StatusBadRequest, "username and password required")
		return
	}
	// хешируем пароль перед сохранением
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		responder.Error(http.StatusInternalServerError, "internal error")
		return
	}
	// пытаемся создать пользователя в хранилище
	if err := c.store.Create(&auth.User{Username: req.Username, Password: hash}); err != nil {
		if err == auth.ErrExists {
			// пользователь уже существует - вернём ошибку (200 OK с сообщением об ошибке)
			responder.Error(http.StatusOK, err.Error())
		} else {
			// непредвиденная ошибка сохранения
			responder.Error(http.StatusInternalServerError, "internal error")
		}
		return
	}
	// успешная регистрация - возвращаем код 201 (Created) без тела ответа
	responder.JSON(http.StatusCreated, nil)
}

// Login обрабатывает авторизацию пользователя и выдачу JWT (эндпоинт /api/login)
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	responder := NewResponder(w)
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder.Error(http.StatusBadRequest, "invalid json")
		return
	}
	// пытаемся получить пользователя из хранилища и проверить пароль
	u, ok := c.store.Get(req.Username)
	if !ok || bcrypt.CompareHashAndPassword(u.Password, []byte(req.Password)) != nil {
		// неверные учётные данные - возвращаем 200 OK с сообщением об ошибке
		responder.Error(http.StatusOK, auth.ErrWrongCreds.Error())
		return
	}
	// учётные данные верны - генерируем JWT токен
	_, tokenString, err := c.tokenAuth.Encode(map[string]interface{}{
		"sub": req.Username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	if err != nil {
		responder.Error(http.StatusInternalServerError, "internal error")
		return
	}
	// успешный вход - возвращаем токен JWT
	responder.JSON(http.StatusOK, tokenResponse{Token: tokenString})
}
