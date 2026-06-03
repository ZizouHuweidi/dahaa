package httpapi

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/zizouhuweidi/dahaa/internal/domain"
	"github.com/zizouhuweidi/dahaa/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *echo.Context) error {
	var req service.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.Username == "" || req.Email == "" || req.Password == "" || req.DisplayName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username, email, password, and display_name are required")
	}
	user, err := h.userService.Register(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, "username or email already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to register user")
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(c *echo.Context) error {
	var req service.LoginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.Username == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username and password are required")
	}
	token, err := h.userService.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid username or password")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to login")
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
