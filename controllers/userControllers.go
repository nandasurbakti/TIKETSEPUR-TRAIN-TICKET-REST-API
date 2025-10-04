package controllers

import (
	"net/http"
	"strconv"
	"tiketsepur/dto"
	"tiketsepur/service"
	"tiketsepur/utils"

	"github.com/gin-gonic/gin"
)

// @title User API
// @description API for managing users
type UserControllers struct {
	userService service.UserService
}

func NewUserControllers(userService service.UserService) *UserControllers {
	return &UserControllers{userService: userService}
}

// Create godoc
// @Summary User admin
// @Description Buat akun admin baru
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "Detail admin yang akan dibuat"
// @Success 201 {object} utils.Response{data=models.User} "User admin berhasil dibuat"
// @Failure 400 {object} utils.Response "Invalid request"
// @Router /users [post]
func (h *UserControllers) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	user, err := h.userService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "gagal membuat user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "user admin berhasil dibuat", user)
}

// GetAll godoc
// @Summary Semua user
// @Description Daftar semua user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.User} "Daftar user"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/users [get]
// @Security BearerAuth
func (h *UserControllers) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mendapatkan user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user berhasil didapatkan", users)
}

// GetByID godoc
// @Summary User by ID
// @Description Detail user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response{data=models.User} "Detail user"
// @Failure 404 {object} utils.Response "User tidak ditemukan"
// @Router /users/{id} [get]
// @Security BearerAuth
func (h *UserControllers) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := h.userService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "user tidak ditemukan", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user berhasil didapatkan", user)
}

// Update godoc
// @Summary Update user
// @Description Update detail user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body dto.UpdateUserRequest true "Detail user yang akan diupdate"
// @Success 200 {object} utils.Response{data=models.User} "User berhasil diupdate"
// @Failure 400 {object} utils.Response "Request tidak valid"
// @Router /users/{id} [put]
// @Security BearerAuth
func (h *UserControllers) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	user, err := h.userService.Update(id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "gagal update user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user berhasil diupdate", user)
}

// Delete godoc
// @Summary Delete user
// @Description Delete akun user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response "User berhasil diupdate"
// @Failure 400 {object} utils.Response "Gagal hapus user"
// @Router /users/{id} [delete]
// @Security BearerAuth
func (h *UserControllers) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.userService.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "gagal hapus user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user berhasil diupdate", nil)
}