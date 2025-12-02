package services

import (
	"bytes"
	"flicknfit_backend/config"
	"flicknfit_backend/utils"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	storage_go "github.com/supabase-community/storage-go"
)

// SupabaseStorageService handles file uploads to Supabase Storage
type SupabaseStorageService interface {
	UploadScanImage(userID uint64, file multipart.File, filename, scanType string) (string, error)
	GetSignedURL(filePath string, expiresIn int) (string, error)
	DeleteFile(filePath string) error
	DownloadAndDecryptFile(filePath string) ([]byte, error)
}

type supabaseStorageService struct {
	client        *storage_go.Client
	bucket        string
	supabaseURL   string
	apiKey        string
	encryptionKey string
}

// NewSupabaseStorageService creates a new Supabase storage service
func NewSupabaseStorageService(cfg *config.Config) SupabaseStorageService {
	if cfg.SupabaseURL == "" || cfg.SupabaseKey == "" {
		panic("Supabase credentials not configured")
	}

	if cfg.EncryptionKey == "" {
		panic("Encryption key not configured")
	}

	client := storage_go.NewClient(cfg.SupabaseURL+"/storage/v1", cfg.SupabaseKey, nil)

	return &supabaseStorageService{
		client:        client,
		bucket:        cfg.SupabaseBucket,
		supabaseURL:   cfg.SupabaseURL,
		apiKey:        cfg.SupabaseKey,
		encryptionKey: cfg.EncryptionKey,
	}
}

// UploadScanImage uploads a scan image to Supabase Storage with encryption
// scanType: "face" or "body"
func (s *supabaseStorageService) UploadScanImage(userID uint64, file multipart.File, filename, scanType string) (string, error) {
	// Reset file pointer
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Encrypt file content
	encryptedBytes, err := utils.EncryptFile(fileBytes, s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt file: %w", err)
	}

	// Generate unique filename with .enc extension
	ext := filepath.Ext(filename)
	uniqueFilename := fmt.Sprintf("%s_%d_%s%s.enc", uuid.New().String(), time.Now().Unix(), scanType, ext)

	// Path format: scans/{scanType}/{userID}/{uniqueFilename}
	filePath := fmt.Sprintf("scans/%s/%d/%s", scanType, userID, uniqueFilename)

	// Upload encrypted file to Supabase
	_, err = s.client.UploadFile(s.bucket, filePath, bytes.NewReader(encryptedBytes))
	if err != nil {
		return "", fmt.Errorf("failed to upload to Supabase: %w", err)
	}

	return filePath, nil
}

// GetSignedURL generates a signed URL for temporary access
func (s *supabaseStorageService) GetSignedURL(filePath string, expiresIn int) (string, error) {
	if expiresIn == 0 {
		expiresIn = 600 // Default 10 minutes
	}

	resp, err := s.client.CreateSignedUrl(s.bucket, filePath, expiresIn)
	if err != nil {
		return "", fmt.Errorf("failed to create signed URL: %w", err)
	}

	return resp.SignedURL, nil
}

// DeleteFile deletes a file from Supabase Storage
func (s *supabaseStorageService) DeleteFile(filePath string) error {
	filePaths := []string{filePath}

	_, err := s.client.RemoveFile(s.bucket, filePaths)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// DownloadAndDecryptFile downloads an encrypted file from Supabase and decrypts it
func (s *supabaseStorageService) DownloadAndDecryptFile(filePath string) ([]byte, error) {
	// Download encrypted file
	fileData, err := s.client.DownloadFile(s.bucket, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	// Decrypt file content
	decryptedBytes, err := utils.DecryptFile(fileData, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt file: %w", err)
	}

	return decryptedBytes, nil
}
