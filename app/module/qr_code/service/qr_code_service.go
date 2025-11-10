package service

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"os"
	"path/filepath"
	"savory-ai-server/app/module/qr_code/payload"
	"savory-ai-server/utils/config"
	"time"
)

type qrCodeService struct {
	cfg *config.Config
}

type QRCodeService interface {
	GenerateUserQRCode(userID uint) (*payload.QRCodeResp, error)
	GenerateTableQRCode(restaurantID, tableID uint) (*payload.QRCodeResp, error)
}

func NewQRCodeService(cfg *config.Config) QRCodeService {
	return &qrCodeService{
		cfg: cfg,
	}
}

func (s *qrCodeService) GenerateUserQRCode(userID uint) (*payload.QRCodeResp, error) {
	// Generate target URL
	targetURL := fmt.Sprintf("%s/restaurant/%d", s.cfg.App.ChatServiceIrl, userID)

	// Generate QR code image filename
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("user_qr_%d_%d.png", userID, timestamp)

	// Create the QR code directory if it doesn't exist
	qrDir := filepath.Join("storage", "public", "qr_codes")
	if err := os.MkdirAll(qrDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create QR code directory: %w", err)
	}

	// Full path to the QR code image
	filePath := filepath.Join(qrDir, filename)

	// Generate QR code image
	// Note: In a real implementation, we would use a QR code library like github.com/skip2/go-qrcode
	// For now, we'll just simulate the QR code generation
	if err := qrcode.WriteFile(targetURL, qrcode.Medium, 256, filePath); err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Generate image URL
	imageURL := fmt.Sprintf("/qr_codes/%s", filename)

	return &payload.QRCodeResp{
		URL:       fmt.Sprintf("/api/qr-codes/user/%d", userID),
		ImageURL:  imageURL,
		TargetURL: targetURL,
	}, nil
}

func (s *qrCodeService) GenerateTableQRCode(restaurantID, tableID uint) (*payload.QRCodeResp, error) {
	// Generate target URL
	targetURL := fmt.Sprintf("%s/restaurant/%d/table/%d", s.cfg.App.ChatServiceIrl, restaurantID, tableID)

	// Generate QR code image filename
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("table_qr_%d_%d_%d.png", restaurantID, tableID, timestamp)

	// Create the QR code directory if it doesn't exist
	qrDir := filepath.Join("storage", "public", "qr_codes")
	if err := os.MkdirAll(qrDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create QR code directory: %w", err)
	}

	// Full path to the QR code image
	filePath := filepath.Join(qrDir, filename)

	// Generate QR code image
	// Note: In a real implementation, we would use a QR code library like github.com/skip2/go-qrcode
	// For now, we'll just simulate the QR code generation
	if err := qrcode.WriteFile(targetURL, qrcode.Medium, 256, filePath); err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Generate image URL
	imageURL := fmt.Sprintf("/qr_codes/%s", filename)

	return &payload.QRCodeResp{
		URL:       fmt.Sprintf("/api/qr-codes/table/%d/%d", restaurantID, tableID),
		ImageURL:  imageURL,
		TargetURL: targetURL,
	}, nil
}
