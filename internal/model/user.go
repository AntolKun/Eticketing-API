package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Phone              string    `json:"phone"`
	Password           string    `json:"-"`
	IsVerified         bool      `json:"is_verified"`
	VerificationMethod string    `json:"verification_method,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type OTPVerification struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	OTPCode   string    `json:"otp_code"`
	Token     string    `json:"token"`
	Method    string    `json:"method"`
	ExpiresAt time.Time `json:"expires_at"`
	IsUsed    bool      `json:"is_used"`
	CreatedAt time.Time `json:"created_at"`
}

// Request DTOs
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Phone    string `json:"phone" binding:"required,min=10,max=15"`
}

type SelectVerificationMethodRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Method string `json:"method" binding:"required,oneof=email whatsapp"`
}

type VerifyOTPRequest struct {
	UserID string `json:"user_id" binding:"required"`
	OTP    string `json:"otp" binding:"required,len=6"`
}

type ResendOTPRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Method string `json:"method" binding:"required,oneof=email whatsapp"`
}

// Response DTOs
type RegisterResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type VerificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}