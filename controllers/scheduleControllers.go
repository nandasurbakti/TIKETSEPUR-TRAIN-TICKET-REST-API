package controllers

import (
	"net/http"
	"strconv"
	"tiketsepur/dto"
	"tiketsepur/service"
	"tiketsepur/utils"

	"github.com/gin-gonic/gin"
)

// @title Schedule API
// @description API for managing train schedules
type ScheduleControllers struct {
	scheduleService service.ScheduleService
}

func NewScheduleControllers(scheduleService service.ScheduleService) *ScheduleControllers {
	return &ScheduleControllers{scheduleService: scheduleService}
}

// Create godoc
// @Summary Buat jadwal baru
// @Description Buat jadwal kereta baru
// @Tags schedules
// @Accept json
// @Produce json
// @Param schedule body dto.CreateScheduleRequest true "Detail jadwal"
// @Success 201 {object} utils.Response{data=models.Schedule} "Jadwal berhasil dibuat"
// @Failure 400 {object} utils.Response "Request tidak valid"
// @Router /schedules [post]
// @Security BearerAuth
func (h *ScheduleControllers) Create(c *gin.Context) {
	var req dto.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	schedule, err := h.scheduleService.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "jadwal gagal dibuat", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "jadwal berhasil dibuat", schedule)
}

// GetAll godoc
// @Summary Semua jadwal
// @Description Semua jadwal kereta
// @Tags schedules
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Schedule} "Daftar jadwal"
// @Failure 500 {object} utils.Response "Gagal mendapatkan jadwal"
// @Router /admin/schedules [get]
// @Security BearerAuth
func (h *ScheduleControllers) GetAll(c *gin.Context) {
	schedules, err := h.scheduleService.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mendapatkan jadwal", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "jadwal berhasil didapatkan", schedules)
}

// GetByID godoc
// @Summary Semua jadwal by ID
// @Description Semua jadwal kereta by its ID
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} utils.Response{data=models.Schedule} "Detail jadwal"
// @Failure 404 {object} utils.Response "Jadwal tidak ditemukan"
// @Router /schedules/{id} [get]
func (h *ScheduleControllers) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := h.scheduleService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "jadwal tidak ditemukan", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "jadwal berhasil didapatkan", user)
}

// Search godoc
// @Summary Search jadwal
// @Description Search jadwal kereta dengan kriteria
// @Tags schedules
// @Accept json
// @Produce json
// @Param origin query string false "Origin station"
// @Param destination query string false "Destination station"
// @Param date query string false "Travel date (YYYY-MM-DD)"
// @Success 200 {object} utils.Response{data=[]models.Schedule} "Daftar jadwal sesuai kriteria"
// @Failure 400 {object} utils.Response "Request tidak valid"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /schedules/search [get]
func (h *ScheduleControllers) Search(c *gin.Context) {
	var req dto.SearchScheduleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	schedules, err := h.scheduleService.Search(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mencari schedules", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "jadwal ditemukan", schedules)
}

func (h *ScheduleControllers) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req dto.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	schedule, err := h.scheduleService.Update(id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "gagal update schedule", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "jadwal berhasil diupdate", schedule)
}

func (h *ScheduleControllers) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.scheduleService.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "gagal menghapus schedule", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "jadwal berhasil dihapus", nil)
}