package repository

import (
	"database/sql"
	"e-ticketing/internal/model"
	"time"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (name, email, phone, password, is_verified)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, user.Name, user.Email, user.Phone, user.Password, false).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, name, email, phone, password, is_verified, verification_method, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password,
		&user.IsVerified, &user.VerificationMethod, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, name, email, phone, password, is_verified, COALESCE(verification_method, '') as verification_method, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password,
		&user.IsVerified, &user.VerificationMethod, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUserVerification(userID uuid.UUID, isVerified bool, method string) error {
	query := `UPDATE users SET is_verified = $1, verification_method = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, isVerified, method, time.Now(), userID)
	return err
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) PhoneExists(phone string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`
	err := r.db.QueryRow(query, phone).Scan(&exists)
	return exists, err
}

// OTP Methods
func (r *UserRepository) CreateOTP(otp *model.OTPVerification) error {
	query := `
		INSERT INTO otp_verifications (user_id, otp_code, token, method, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	return r.db.QueryRow(query, otp.UserID, otp.OTPCode, otp.Token, otp.Method, otp.ExpiresAt).
		Scan(&otp.ID, &otp.CreatedAt)
}

func (r *UserRepository) GetValidOTP(userID uuid.UUID, otpCode string) (*model.OTPVerification, error) {
	otp := &model.OTPVerification{}
	query := `
		SELECT id, user_id, otp_code, COALESCE(token, '') as token, method, expires_at, is_used, created_at
		FROM otp_verifications
		WHERE user_id = $1 AND otp_code = $2 AND is_used = false AND expires_at > $3
		ORDER BY created_at DESC
		LIMIT 1`

	err := r.db.QueryRow(query, userID, otpCode, time.Now()).Scan(
		&otp.ID, &otp.UserID, &otp.OTPCode, &otp.Token, &otp.Method, &otp.ExpiresAt, &otp.IsUsed, &otp.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return otp, nil
}

func (r *UserRepository) GetValidOTPByToken(token string) (*model.OTPVerification, error) {
	otp := &model.OTPVerification{}
	query := `
		SELECT id, user_id, otp_code, token, method, expires_at, is_used, created_at
		FROM otp_verifications
		WHERE token = $1 AND is_used = false AND expires_at > $2
		ORDER BY created_at DESC
		LIMIT 1`

	err := r.db.QueryRow(query, token, time.Now()).Scan(
		&otp.ID, &otp.UserID, &otp.OTPCode, &otp.Token, &otp.Method, &otp.ExpiresAt, &otp.IsUsed, &otp.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return otp, nil
}

func (r *UserRepository) MarkOTPAsUsed(otpID uuid.UUID) error {
	query := `UPDATE otp_verifications SET is_used = true WHERE id = $1`
	_, err := r.db.Exec(query, otpID)
	return err
}

func (r *UserRepository) InvalidateOldOTPs(userID uuid.UUID) error {
	query := `UPDATE otp_verifications SET is_used = true WHERE user_id = $1 AND is_used = false`
	_, err := r.db.Exec(query, userID)
	return err
}