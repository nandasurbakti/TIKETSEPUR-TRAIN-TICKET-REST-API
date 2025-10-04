package controllers

import (
	"net/http"
	"strconv"
	"tiketsepur/dto"
	"tiketsepur/service"
	"tiketsepur/utils"

	"github.com/gin-gonic/gin"
)

// @title Ticket API
// @description API for managing train tickets
type TicketControllers struct {
	ticketService service.TicketService
}

func NewTicketControllers(ticketService service.TicketService) *TicketControllers {
	return &TicketControllers{ticketService: ticketService}
}

// Create godoc
// @Summary Pesan tiket baru
// @Description Pesan tiket kereta baru
// @Tags tickets
// @Accept json
// @Produce json
// @Param ticket body dto.CreateTicketRequest true "Rincian pemesanan tiket"
// @Success 201 {object} utils.Response{data=models.Ticket} "Tiket berhasil dipesan"
// @Failure 400 {object} utils.Response "Request tidak valid"
// @Router /tickets [post]
// @Security BearerAuth
func (h *TicketControllers) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req dto.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	ticket, err := h.ticketService.Create(c.Request.Context(), userID.(int), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "gagal membuat tiket", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "tiket berhasil dipesan", ticket)
}

// GetMyTickets godoc
// @Summary Tiket user
// @Description Semua tiket yang dipesan oleh pengguna saat ini
// @Tags tickets
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Ticket} "Daftar tiket pengguna"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /tickets/my-tickets [get]
// @Security BearerAuth
func (h *TicketControllers) GetMyTickets(c *gin.Context) {
	userID, _ := c.Get("user_id")

	tickets, err := h.ticketService.GetByUserID(userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mendapatkan daftar tiket", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "daftar tiket berhasil didapatkan", tickets)
}

// GetAll godoc
// @Summary Semua tiket
// @Description Semua tiket di dalam sistem (admin only)
// @Tags tickets
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Ticket} "Daftar semua tiket"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/tickets [get]
// @Security BearerAuth
func (h *TicketControllers) GetAll(c *gin.Context) {
	tickets, err := h.ticketService.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mendapatkan daftar tiket", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "daftar tiket berhasil didapatkan", tickets)
}

// GetByID godoc
// @Summary Tiket by ID
// @Description Detail tiket by ticket ID
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} utils.Response{data=models.Ticket} "Detail tiket"
// @Failure 404 {object} utils.Response "Tiket tidak ditemukan"
// @Router /tickets/{id} [get]
// @Security BearerAuth
func (h *TicketControllers) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	ticket, err := h.ticketService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "tiket tidak ditemukan", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "tiket berhasil didapatkan", ticket)
}

// Cancel godoc
// @Summary Cancel ticket
// @Description Cancel tiket yang dipesan
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} utils.Response "Tiket berhasil dibatalkan"
// @Failure 400 {object} utils.Response "Gagal membatalkan tiket"
// @Router /tickets/{id}/cancel [post]
// @Security BearerAuth
func (h *TicketControllers) Cancel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")

	if err := h.ticketService.Cancel(c.Request.Context(), id, userID.(int), role.(string)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "gagal membatalkan tiket", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "tiket berhasil dibatalkan", nil)
}