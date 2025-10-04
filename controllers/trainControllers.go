package controllers

import (
	"net/http"
	"strconv"
	"tiketsepur/dto"
	"tiketsepur/service"
	"tiketsepur/utils"

	"github.com/gin-gonic/gin"
)

// @title Train API
// @description API for managing trains
type TrainControllers struct {
	trainService service.TrainService
}

func NewTrainControllers(trainService service.TrainService) *TrainControllers {
	return &TrainControllers{trainService: trainService}
}

// Create godoc
// @Summary Buat kereta baru
// @Description Buat kereta baru
// @Tags trains
// @Accept json
// @Produce json
// @Param train body dto.CreateTrainRequest true "Detail kereta"
// @Success 201 {object} utils.Response{data=models.Train} "Kereta berhasil dibuat"
// @Failure 400 {object} utils.Response "Request tidak valid"
// @Router /admin/trains [post]
// @Security BearerAuth
func (h *TrainControllers) Create(c *gin.Context) {
	var req dto.CreateTrainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	train, err := h.trainService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "gagal membuat kereta", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "kereta berhasil dibuat", train)
}

// GetAll godoc
// @Summary Semua kereta
// @Description Semua daftar kereta
// @Tags trains
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Train} "Daftar semua kereta"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /trains [get]
func (h *TrainControllers) GetAll(c *gin.Context) {
	trains, err := h.trainService.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mendapatkan kereta", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "kereta berhasil didapatkan", trains)
}

// GetByID godoc
// @Summary Kereta by ID
// @Description Detail kereta by ID
// @Tags trains
// @Accept json
// @Produce json
// @Param id path int true "Train ID"
// @Success 200 {object} utils.Response{data=models.Train} "Detail kereta"
// @Failure 404 {object} utils.Response "Kereta tidak ditemukan"
// @Router /trains/{id} [get]
func (h *TrainControllers) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	train, err := h.trainService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "kereta tidak ditemukan", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "kereta berhasil didapatkan", train)
}

// Update godoc
// @Summary Update kereta
// @Description Update detail kereta
// @Tags trains
// @Accept json
// @Produce json
// @Param id path int true "Train ID"
// @Param train body dto.UpdateTrainRequest true "Detail kereta yang akan diupdate"
// @Success 200 {object} utils.Response{data=models.Train} "Kereta berhasil diupdate"
// @Failure 400 {object} utils.Response "Request tidak valid"
// @Router /admin/trains/{id} [put]
// @Security BearerAuth
func (h *TrainControllers) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req dto.UpdateTrainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	train, err := h.trainService.Update(id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "gagal update kereta", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "kereta berhasil diupdate", train)
}

func (h *TrainControllers) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.trainService.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "gagal delete train", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "kereta berhasil dihapus", nil)
}