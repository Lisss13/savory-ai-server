package service

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"os"
	"path/filepath"
	"savory-ai-server/utils/config"
)

type qrCodeService struct {
	cfg *config.Config
}

type QRCodeService interface {
	GenerateRestaurantQRCode(restaurantID uint) (string, error)
	GenerateTableQRCode(restaurantID, tableID uint) (string, error)
}

func NewQRCodeService(cfg *config.Config) QRCodeService {
	return &qrCodeService{
		cfg: cfg,
	}
}

func (s *qrCodeService) GenerateRestaurantQRCode(restaurantID uint) (string, error) {

	targetURL := fmt.Sprintf("%s/restaurant/%d", s.cfg.App.ChatServiceIrl, restaurantID)

	// Generate QR code image filename
	filename := fmt.Sprintf("restaurant_qr_%d.png", restaurantID)

	// Create the QR code directory if it doesn't exist
	qrDir := filepath.Join("storage", "public", "qr_codes")
	if err := os.MkdirAll(qrDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create QR code directory: %w", err)
	}

	filePath := filepath.Join(qrDir, filename)
	if err := qrcode.WriteFile(targetURL, qrcode.Medium, 256, filePath); err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Generate image URL
	imageURL := fmt.Sprintf("/qr_codes/%s", filename)

	return imageURL, nil
}

func (s *qrCodeService) GenerateTableQRCode(restaurantID, tableID uint) (string, error) {

	targetURL := fmt.Sprintf("%s/restaurant/%d/table/%d", s.cfg.App.ChatServiceIrl, restaurantID, tableID)

	// Generate QR code image filename
	filename := fmt.Sprintf("table_qr_%d_%d.png", restaurantID, tableID)

	// Create the QR code directory if it doesn't exist
	qrDir := filepath.Join("storage", "public", "qr_codes")
	if err := os.MkdirAll(qrDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create QR code directory: %w", err)
	}

	filePath := filepath.Join(qrDir, filename)
	if err := qrcode.WriteFile(targetURL, qrcode.Medium, 256, filePath); err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Generate image URL
	imageURL := fmt.Sprintf("/qr_codes/%s", filename)

	return imageURL, nil
}
