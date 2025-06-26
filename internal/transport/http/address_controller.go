package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/your-org/geo-service-swagger/internal/address"
)

// AddressController отвечает за обработку запросов к адресному сервису (поиск и геокодирование)
type AddressController struct {
	svc      address.Service
	validate *validator.Validate
}

// NewAddressController создает новый контроллер адресов с заданным сервисом
func NewAddressController(svc address.Service) *AddressController {
	return &AddressController{
		svc:      svc,
		validate: validator.New(),
	}
}

// Search обрабатывает запрос поиска адресов по текстовому запросу (эндпоинт /api/address/search)
func (c *AddressController) Search(w http.ResponseWriter, r *http.Request) {
	responder := NewResponder(w)
	var req address.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder.Error(http.StatusBadRequest, "invalid json")
		return
	}
	// валидация входных данных (поле Query обязательно и мин. длина 2 символа)
	if err := c.validate.Struct(&req); err != nil {
		responder.Error(http.StatusBadRequest, "query is required")
		return
	}
	// вызов бизнес-логики поиска адресов
	results, err := c.svc.Search(r.Context(), req.Query)
	if err != nil {
		responder.Error(http.StatusInternalServerError, "upstream error")
		return
	}
	// успешный ответ с найденными адресами
	responder.JSON(http.StatusOK, address.Response{Addresses: results})
}

// Geocode обрабатывает запрос реверс-геокодирования (поиск адреса по координатам, эндпоинт /api/address/geocode)
func (c *AddressController) Geocode(w http.ResponseWriter, r *http.Request) {
	responder := NewResponder(w)
	var req address.GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder.Error(http.StatusBadRequest, "invalid json")
		return
	}
	// валидация входных данных (оба поля Lat и Lng обязательны)
	if err := c.validate.Struct(&req); err != nil {
		responder.Error(http.StatusBadRequest, "lat and lng required")
		return
	}
	// вызов бизнес-логики геокодирования
	results, err := c.svc.Geocode(r.Context(), req.Lat, req.Lng)
	if err != nil {
		responder.Error(http.StatusInternalServerError, "upstream error")
		return
	}
	// успешный ответ с результатом (адреса ближайшие к координатам)
	responder.JSON(http.StatusOK, address.Response{Addresses: results})
}
