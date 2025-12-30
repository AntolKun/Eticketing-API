package service

import (
	"database/sql"
	"e-ticketing/config"
	"e-ticketing/internal/model"
	"e-ticketing/internal/repository"
	"e-ticketing/pkg/utils"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	emailSvc    *EmailService
	whatsappSvc *WhatsAppService
	config      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, emailSvc *EmailService, whatsappSvc *WhatsAppService, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		emailSvc:    emailSvc,
		whatsappSvc: whatsappSvc,
		config:      cfg,
	}
}

func (s *AuthService) Register(req *model.RegisterRequest) (*model.RegisterResponse, error) {
	// Check if email exists
	emailExists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, errors.New("email sudah terdaftar")
	}

	// Check if phone exists
	phoneExists, err := s.userRepo.PhoneExists(req.Phone)
	if err != nil {
		return nil, err
	}
	if phoneExists {
		return nil, errors.New("nomor telepon sudah terdaftar")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return &model.RegisterResponse{
		UserID:  user.ID.String(),
		Message: "Registrasi berhasil. Silakan pilih metode verifikasi",
	}, nil
}

func (s *AuthService) SelectVerificationMethod(req *model.SelectVerificationMethodRequest, baseURL string) error {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return errors.New("user ID tidak valid")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user tidak ditemukan")
		}
		return err
	}

	if user.IsVerified {
		return errors.New("user sudah terverifikasi")
	}

	// Invalidate old OTPs
	s.userRepo.InvalidateOldOTPs(userID)

	// Generate OTP
	otpCode := utils.GenerateOTP(6)
	token := utils.GenerateToken()

	otp := &model.OTPVerification{
		UserID:    userID,
		OTPCode:   otpCode,
		Token:     token,
		Method:    req.Method,
		ExpiresAt: time.Now().Add(time.Duration(s.config.OTPExpiryMinutes) * time.Minute),
	}

	if err := s.userRepo.CreateOTP(otp); err != nil {
		return err
	}

	// Send OTP based on method
	switch req.Method {
	case "email":
		link := utils.GenerateVerificationLink(baseURL, token)
		if err := s.emailSvc.SendVerificationLinkEmail(user.Email, link, user.Name); err != nil {
			return errors.New("gagal mengirim email: " + err.Error())
		}
	case "whatsapp":
		if err := s.whatsappSvc.SendOTP(user.Phone, otpCode, user.Name); err != nil {
			return errors.New("gagal mengirim WhatsApp: " + err.Error())
		}
	}

	return nil
}

func (s *AuthService) VerifyOTP(req *model.VerifyOTPRequest) (*model.VerificationResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.New("user ID tidak valid")
	}

	otp, err := s.userRepo.GetValidOTP(userID, req.OTP)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("OTP tidak valid atau sudah expired")
		}
		return nil, err
	}

	// Mark OTP as used
	if err := s.userRepo.MarkOTPAsUsed(otp.ID); err != nil {
		return nil, err
	}

	// Update user verification status
	if err := s.userRepo.UpdateUserVerification(userID, true, otp.Method); err != nil {
		return nil, err
	}

	return &model.VerificationResponse{
		Success: true,
		Message: "Verifikasi berhasil! Akun Anda sudah aktif",
	}, nil
}

func (s *AuthService) VerifyEmailToken(token string) (*model.VerificationResponse, error) {
	otp, err := s.userRepo.GetValidOTPByToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("token tidak valid atau sudah expired")
		}
		return nil, err
	}

	// Mark OTP as used
	if err := s.userRepo.MarkOTPAsUsed(otp.ID); err != nil {
		return nil, err
	}

	// Update user verification status
	if err := s.userRepo.UpdateUserVerification(otp.UserID, true, otp.Method); err != nil {
		return nil, err
	}

	return &model.VerificationResponse{
		Success: true,
		Message: "Email berhasil diverifikasi! Akun Anda sudah aktif",
	}, nil
}

func (s *AuthService) ResendOTP(req *model.ResendOTPRequest, baseURL string) error {
	selectReq := &model.SelectVerificationMethodRequest{
		UserID: req.UserID,
		Method: req.Method,
	}
	return s.SelectVerificationMethod(selectReq, baseURL)
}