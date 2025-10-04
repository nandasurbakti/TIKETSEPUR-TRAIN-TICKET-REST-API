package controllers

import (
	"net/http"
	"tiketsepur/dto"
	"tiketsepur/service"
	"tiketsepur/utils"

	"github.com/gin-gonic/gin"
)

// @title Authentication API
// @description API for user authentication and authorization
type AuthControllers struct {
	authService service.AuthService
	userService service.UserService
}

func NewAuthControllers(authService service.AuthService, userService service.UserService) *AuthControllers {
    return &AuthControllers{
        authService: authService,
        userService: userService,
    }
}

// Register godoc
// @Summary Register user baru
// @Description Register akun user baru
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "detail registrasi user"
// @Success 201 {object} utils.Response{data=models.User} "registrasi berhasil"
// @Failure 400 {object} utils.Response "request tidak valid"
// @Router /auth/register [post]
func (h *AuthControllers) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "registrasi gagal", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "registrasi berhasil", user)
}

// RegisterAdmin godoc
// @Summary Register admin user
// @Description Register akun admin baru
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "admin registrasi sukses"
// @Success 201 {object} utils.Response{data=models.User} "registrasi berhasil"
// @Failure 400 {object} utils.Response "request tidak valid"
// @Router /auth/register-admin [post]
// @Security BearerAuth
func (h *AuthControllers) RegisterAdmin(c *gin.Context) {
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "request tidak valid", err)
        return
    }

    user, err := h.userService.Create(req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "registrasi gagal", err)
        return
    }

    utils.SuccessResponse(c, http.StatusCreated, "admin registrasi berhasil", user)
}

// Login godoc
// @Summary User login
// @Description Authenticate user dan memberi JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} utils.Response{data=dto.LoginResponse} "Login sukses"
// @Failure 400 {object} utils.Response "request tidak valid"
// @Failure 401 {object} utils.Response "Authentication gagal"
// @Router /auth/login [post]
func (h *AuthControllers) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "login gagal", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "login sukses", result)
}

// @Summary User logout
// @Description Logout user and invalidate token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "Logout sukses"
// @Failure 500 {object} utils.Response "Logout gagal"
// @Router /auth/logout [post]
// @Security BearerAuth
func (h *AuthControllers) Logout(c *gin.Context) {
	token, _ := c.Get("token")

	if err := h.authService.Logout(c.Request.Context(), token.(string)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Logout failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "logout sukses", nil)
}

// @Summary Get current user info
// @Description Dapatkan informasi tentang user yang saat ini masuk
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=map[string]interface{}} "User info didapatkan"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /auth/me [get]
// @Security BearerAuth
func (h *AuthControllers) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("user_email")
	role, _ := c.Get("user_role")

	user := map[string]interface{}{
		"id":    userID,
		"email": email,
		"role":  role,
	}

	utils.SuccessResponse(c, http.StatusOK, "user info didapatkan", user)
}