package service

import (
	"e-ticketing/config"
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{config: cfg}
}

func (s *EmailService) SendOTPEmail(to, otp, name string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Kode OTP Verifikasi - E-Ticketing")

	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; padding: 20px;">
			<h2>Halo %s!</h2>
			<p>Kode OTP Anda untuk verifikasi akun:</p>
			<div style="background-color: #f4f4f4; padding: 20px; text-align: center; margin: 20px 0;">
				<h1 style="color: #333; letter-spacing: 10px; margin: 0;">%s</h1>
			</div>
			<p>Kode ini berlaku selama <strong>%d menit</strong>.</p>
			<p>Jika Anda tidak meminta kode ini, abaikan email ini.</p>
			<br>
			<p>Salam,<br>Tim E-Ticketing</p>
		</body>
		</html>
	`, name, otp, s.config.OTPExpiryMinutes)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.config.SMTPHost, s.config.SMTPPort, s.config.SMTPUser, s.config.SMTPPassword)

	return d.DialAndSend(m)
}

func (s *EmailService) SendVerificationLinkEmail(to, link, name string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verifikasi Email - E-Ticketing")

	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; padding: 20px;">
			<h2>Halo %s!</h2>
			<p>Klik tombol di bawah untuk verifikasi email Anda:</p>
			<div style="text-align: center; margin: 30px 0;">
				<a href="%s" style="background-color: #4CAF50; color: white; padding: 15px 30px; text-decoration: none; border-radius: 5px; font-size: 16px;">
					Verifikasi Email
				</a>
			</div>
			<p>Atau copy link berikut ke browser:</p>
			<p style="word-break: break-all; color: #666;">%s</p>
			<p>Link ini berlaku selama <strong>%d menit</strong>.</p>
			<br>
			<p>Salam,<br>Tim E-Ticketing</p>
		</body>
		</html>
	`, name, link, link, s.config.OTPExpiryMinutes)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.config.SMTPHost, s.config.SMTPPort, s.config.SMTPUser, s.config.SMTPPassword)

	return d.DialAndSend(m)
}