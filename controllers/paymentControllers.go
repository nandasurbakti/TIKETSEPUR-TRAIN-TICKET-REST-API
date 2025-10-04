package controllers

import (
	"net/http"
	"tiketsepur/service"
	"tiketsepur/utils"

	"github.com/gin-gonic/gin"
)

// @title Payment API
// @description API for managing payments
type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// ConfirmPayment godoc
// @Summary Confirm payment
// @Description Konfirmasi pembayaran untuk pemesanan tiket
// @Tags payments
// @Accept json
// @Produce json
// @Param paymentCode path string true "Payment Code"
// @Success 200 {object} utils.Response{data=map[string]string} "payment confirmed sukses"
// @Failure 400 {object} utils.Response "Payment code tidak valid"
// @Router /payments/confirm/{paymentCode} [post]
// @Security BearerAuth
func (h *PaymentHandler) ConfirmPayment(c *gin.Context) {
	paymentCode := c.Param("paymentCode")
	if paymentCode == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "payment code diperlukan", nil)
		return
	}

	if err := h.paymentService.ConfirmPayment(c.Request.Context(), paymentCode); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "payment confirmation gagal", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "payment confirmed sukses", gin.H{
		"message": "tiket kamu sudah bisa digunakan. cek email kamu untuk detail pemesanan.",
		"payment_code": paymentCode,
	})
}

// GetPaymentStatus godoc
// @Summary Get payment status
// @Description Status terkini dari pembayaran
// @Tags payments
// @Accept json
// @Produce json
// @Param paymentCode path string true "Payment Code"
// @Success 200 {object} utils.Response{data=map[string]string} "Payment status"
// @Failure 400 {object} utils.Response "Payment code tidak valid"
// @Router /payments/status/{paymentCode} [get]
// @Security BearerAuth
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	paymentCode := c.Param("paymentCode")
	if paymentCode == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Payment code diperlukan", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "payment status didapatkan", gin.H{
		"payment_code": paymentCode,
		"status": "pending",
	})
}