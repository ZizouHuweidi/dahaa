package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zizouhuweidi/dahaa/internal/service"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body service.RegisterRequest true "User registration data"
// @Success 201 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /users/register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var req service.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
	}

	user, err := h.userService.Register(c.Request().Context(), req)
	if err != nil {
		switch err {
		case service.ErrUserAlreadyExists:
			return c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Username or email already exists",
			})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to register user",
			})
		}
	}

	return c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return session token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body service.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var req service.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
	}

	token, err := h.userService.Login(c.Request().Context(), req)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			return c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Invalid username or password",
			})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to login",
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

// SendGameInvite godoc
// @Summary Send game invitation
// @Description Send a game invitation to another user
// @Tags users
// @Accept json
// @Produce json
// @Param game_id path string true "Game ID"
// @Param to_user_id path string true "Recipient User ID"
// @Success 200 {object} domain.GameInvite
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/invites/{game_id}/{to_user_id} [post]
func (h *UserHandler) SendGameInvite(c echo.Context) error {
	gameID := c.Param("game_id")
	toUserID := c.Param("to_user_id")
	fromUserID := c.Get("user_id").(string)

	err := h.userService.SendGameInvite(c.Request().Context(), gameID, fromUserID, toUserID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "User not found",
			})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to send invitation",
			})
		}
	}

	return c.NoContent(http.StatusOK)
}

// AcceptGameInvite godoc
// @Summary Accept game invitation
// @Description Accept a pending game invitation
// @Tags users
// @Accept json
// @Produce json
// @Param invite_id path string true "Invitation ID"
// @Success 200 {object} domain.GameInvite
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/invites/{invite_id}/accept [post]
func (h *UserHandler) AcceptGameInvite(c echo.Context) error {
	inviteID := c.Param("invite_id")

	err := h.userService.AcceptGameInvite(c.Request().Context(), inviteID)
	if err != nil {
		switch err {
		case service.ErrInviteNotFound:
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Invitation not found",
			})
		case service.ErrInviteExpired:
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invitation has expired",
			})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to accept invitation",
			})
		}
	}

	return c.NoContent(http.StatusOK)
}

// DeclineGameInvite godoc
// @Summary Decline game invitation
// @Description Decline a pending game invitation
// @Tags users
// @Accept json
// @Produce json
// @Param invite_id path string true "Invitation ID"
// @Success 200 {object} domain.GameInvite
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/invites/{invite_id}/decline [post]
func (h *UserHandler) DeclineGameInvite(c echo.Context) error {
	inviteID := c.Param("invite_id")

	err := h.userService.DeclineGameInvite(c.Request().Context(), inviteID)
	if err != nil {
		switch err {
		case service.ErrInviteNotFound:
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Invitation not found",
			})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to decline invitation",
			})
		}
	}

	return c.NoContent(http.StatusOK)
}

// GetPendingInvites godoc
// @Summary Get pending invitations
// @Description Get all pending game invitations for the current user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} domain.GameInvite
// @Failure 500 {object} ErrorResponse
// @Router /users/invites [get]
func (h *UserHandler) GetPendingInvites(c echo.Context) error {
	userID := c.Get("user_id").(string)

	invites, err := h.userService.GetPendingInvites(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get pending invitations",
		})
	}

	return c.JSON(http.StatusOK, invites)
}
