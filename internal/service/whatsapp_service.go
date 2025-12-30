package service

import (
	"bytes"
	"e-ticketing/config"
	"encoding/json"
	"fmt"
	"net/http"
)

type WhatsAppService struct {
	config *config.Config
}

func NewWhatsAppService(cfg *config.Config) *WhatsAppService {
	return &WhatsAppService{config: cfg}
}

func (s *WhatsAppService) SendOTP(phone, otp, name string) error {
	message := fmt.Sprintf(
		"Halo %s!\n\nKode OTP Anda: *%s*\n\nKode berlaku %d menit.\n\nJangan bagikan kode ini kepada siapapun.\n\n- Tim E-Ticketing",
		name, otp, s.config.OTPExpiryMinutes,
	)

	payload := map[string]interface{}{
		"target":  phone,
		"message": message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.config.WAAPIUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.config.WAAPIToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("WhatsApp API error: status %d", resp.StatusCode)
	}

	return nil
}