package qrcode

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/skip2/go-qrcode"
)

// QRCodeGenerator handles QR code generation and file management
type QRCodeGenerator struct {
	staticDir string
	baseURL   string
}

// NewQRCodeGenerator creates a new QR code generator
func NewQRCodeGenerator(staticDir, baseURL string) *QRCodeGenerator {
	return &QRCodeGenerator{
		staticDir: staticDir,
		baseURL:   baseURL,
	}
}

// GenerateQRCode generates a QR code image and returns the URL to access it
func (g *QRCodeGenerator) GenerateQRCode(content string) (string, error) {
	// Ensure static directory exists
	if err := os.MkdirAll(g.staticDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create static directory: %w", err)
	}

	// Generate unique filename
	filename, err := g.generateUniqueFilename()
	if err != nil {
		return "", fmt.Errorf("failed to generate unique filename: %w", err)
	}

	// Full path to the image file
	filePath := filepath.Join(g.staticDir, filename)

	// Generate QR code image
	err = qrcode.WriteFile(content, qrcode.Medium, 256, filePath)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Return URL to access the image
	// baseURL already includes the full path (e.g., http://localhost:6679/static)
	imageURL := fmt.Sprintf("%s/%s", g.baseURL, filename)

	// Schedule cleanup after 5 minutes
	go g.scheduleCleanup(filePath, 5*time.Minute)

	return imageURL, nil
}

// generateUniqueFilename creates a unique filename for the QR code image
func (g *QRCodeGenerator) generateUniqueFilename() (string, error) {
	// Generate random bytes
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Create filename with timestamp and random hex
	timestamp := time.Now().Unix()
	randomHex := hex.EncodeToString(bytes)
	filename := fmt.Sprintf("qr_%d_%s.png", timestamp, randomHex[:8])

	return filename, nil
}

// scheduleCleanup removes the QR code file after the specified duration
func (g *QRCodeGenerator) scheduleCleanup(filePath string, delay time.Duration) {
	time.Sleep(delay)

	if err := os.Remove(filePath); err != nil {
		// Log error but don't fail - file might already be deleted
		fmt.Printf("Warning: failed to cleanup QR code file %s: %v\n", filePath, err)
	}
}

// CleanupExpiredFiles removes QR code files older than the specified duration
func (g *QRCodeGenerator) CleanupExpiredFiles(maxAge time.Duration) error {
	entries, err := os.ReadDir(g.staticDir)
	if err != nil {
		return fmt.Errorf("failed to read static directory: %w", err)
	}

	now := time.Now()
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if it's a QR code file
		if !g.isQRCodeFile(entry.Name()) {
			continue
		}

		filePath := filepath.Join(g.staticDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Remove if older than maxAge
		if now.Sub(info.ModTime()) > maxAge {
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("Warning: failed to cleanup expired QR code file %s: %v\n", filePath, err)
			}
		}
	}

	return nil
}

// isQRCodeFile checks if the filename matches QR code file pattern
func (g *QRCodeGenerator) isQRCodeFile(filename string) bool {
	return len(filename) > 3 && filename[:3] == "qr_" && filepath.Ext(filename) == ".png"
}

// StartPeriodicCleanup starts a background goroutine to periodically clean up expired files
func (g *QRCodeGenerator) StartPeriodicCleanup(interval, maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := g.CleanupExpiredFiles(maxAge); err != nil {
				fmt.Printf("Warning: periodic cleanup failed: %v\n", err)
			}
		}
	}()
}
