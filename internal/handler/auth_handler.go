package handler

import (
	"e-ticketing/internal/model"
	"e-ticketing/internal/service"
	"e-ticketing/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", err.Error())
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Registrasi gagal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, response.Message, gin.H{
		"user_id": response.UserID,
	})
}

func (h *AuthHandler) SelectVerificationMethod(c *gin.Context) {
	var req model.SelectVerificationMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", err.Error())
		return
	}

	baseURL := c.Request.Host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	fullBaseURL := scheme + "://" + baseURL

	err := h.authService.SelectVerificationMethod(&req, fullBaseURL)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal mengirim OTP", err.Error())
		return
	}

	message := "OTP berhasil dikirim via " + req.Method
	if req.Method == "email" {
		message = "Link verifikasi berhasil dikirim ke email"
	}

	utils.SuccessResponse(c, http.StatusOK, message, nil)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req model.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", err.Error())
		return
	}

	response, err := h.authService.VerifyOTP(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Verifikasi gagal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, response.Message, gin.H{
		"verified": response.Success,
	})
}

func (h *AuthHandler) VerifyEmailToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Token tidak ditemukan", "Token is required")
		return
	}

	response, err := h.authService.VerifyEmailToken(token)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Verifikasi gagal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, response.Message, gin.H{
		"verified": response.Success,
	})
}

func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var req model.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", err.Error())
		return
	}

	baseURL := c.Request.Host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	fullBaseURL := scheme + "://" + baseURL

	err := h.authService.ResendOTP(&req, fullBaseURL)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal mengirim ulang OTP", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OTP berhasil dikirim ulang via "+req.Method, nil)
}