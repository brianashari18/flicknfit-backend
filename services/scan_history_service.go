package services

import (
	"flicknfit_backend/dtos"
	"flicknfit_backend/models"
	"flicknfit_backend/repositories"
	"fmt"
	"log"
	"mime/multipart"
)

// ScanHistoryService handles business logic for scan history
type ScanHistoryService interface {
	// Face scan history
	SaveFaceScanHistory(userID uint64, file multipart.File, filename string, result *dtos.SkinColorTonePredictionResponseDTO) error
	GetFaceScanHistories(userID uint64) (*dtos.FaceScanHistoryListDTO, error)
	DeleteFaceScanHistory(userID, historyID uint64) error

	// Body scan history
	SaveBodyScanHistory(userID uint64, file multipart.File, filename, gender string, result interface{}) error
	GetBodyScanHistories(userID uint64) (*dtos.BodyScanHistoryListDTO, error)
	DeleteBodyScanHistory(userID, historyID uint64) error
}

type scanHistoryService struct {
	faceScanRepo   repositories.FaceScanHistoryRepository
	bodyScanRepo   repositories.BodyScanHistoryRepository
	storageService SupabaseStorageService
}

// NewScanHistoryService creates a new scan history service
func NewScanHistoryService(
	faceScanRepo repositories.FaceScanHistoryRepository,
	bodyScanRepo repositories.BodyScanHistoryRepository,
	storageService SupabaseStorageService,
) ScanHistoryService {
	return &scanHistoryService{
		faceScanRepo:   faceScanRepo,
		bodyScanRepo:   bodyScanRepo,
		storageService: storageService,
	}
}

// SaveFaceScanHistory saves face scan result to database and uploads image to Supabase
func (s *scanHistoryService) SaveFaceScanHistory(
	userID uint64,
	file multipart.File,
	filename string,
	result *dtos.SkinColorTonePredictionResponseDTO,
) error {
	// Upload image to Supabase
	imagePath, err := s.storageService.UploadScanImage(userID, file, filename, "face")
	if err != nil {
		log.Printf("[ScanHistoryService] Failed to upload face scan image: %v", err)
		return fmt.Errorf("failed to upload image: %w", err)
	}

	// Count existing scans to generate name
	count, err := s.faceScanRepo.CountByUserID(userID)
	if err != nil {
		log.Printf("[ScanHistoryService] Failed to count face scans: %v", err)
		return fmt.Errorf("failed to count scans: %w", err)
	}

	scanName := fmt.Sprintf("ScanWajah_%d", count+1)

	// Create history record
	history := &models.FaceScanHistory{
		UserID:               userID,
		ScanName:             scanName,
		SkinTone:             result.SkinTone,
		ImagePath:            imagePath,
		ColorRecommendations: result.ColorRecommendations,
		Confidence:           nil, // Optional field
	}

	if err := s.faceScanRepo.Create(history); err != nil {
		log.Printf("[ScanHistoryService] Failed to save face scan history: %v", err)
		// Try to delete uploaded image on failure
		_ = s.storageService.DeleteFile(imagePath)
		return fmt.Errorf("failed to save history: %w", err)
	}

	log.Printf("[ScanHistoryService] Face scan history saved successfully: %s", scanName)
	return nil
}

// GetFaceScanHistories retrieves all face scan histories for a user with signed URLs
func (s *scanHistoryService) GetFaceScanHistories(userID uint64) (*dtos.FaceScanHistoryListDTO, error) {
	histories, err := s.faceScanRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve histories: %w", err)
	}

	historyDTOs := make([]dtos.FaceScanHistoryDTO, 0, len(histories))
	for _, h := range histories {
		historyDTOs = append(historyDTOs, dtos.FaceScanHistoryDTO{
			ID:                   h.ID,
			ScanName:             h.ScanName,
			SkinTone:             h.SkinTone,
			ImageURL:             h.ImagePath, // Return path, client will request via /history/image endpoint
			ColorRecommendations: h.ColorRecommendations,
			Confidence:           h.Confidence,
			CreatedAt:            h.CreatedAt,
		})
	}

	return &dtos.FaceScanHistoryListDTO{
		Histories: historyDTOs,
		Total:     int64(len(historyDTOs)),
	}, nil
}

// DeleteFaceScanHistory deletes a face scan history and its image
func (s *scanHistoryService) DeleteFaceScanHistory(userID, historyID uint64) error {
	// Get history to verify ownership and get image path
	history, err := s.faceScanRepo.FindByID(historyID)
	if err != nil {
		return fmt.Errorf("history not found: %w", err)
	}

	// Verify ownership
	if history.UserID != userID {
		return fmt.Errorf("unauthorized: scan does not belong to user")
	}

	// Delete from database
	if err := s.faceScanRepo.Delete(historyID); err != nil {
		return fmt.Errorf("failed to delete history: %w", err)
	}

	// Delete image from Supabase (best effort, don't fail if image deletion fails)
	if err := s.storageService.DeleteFile(history.ImagePath); err != nil {
		log.Printf("[ScanHistoryService] Warning: Failed to delete image %s: %v", history.ImagePath, err)
	}

	return nil
}

// SaveBodyScanHistory saves body scan result to database and uploads image to Supabase
func (s *scanHistoryService) SaveBodyScanHistory(
	userID uint64,
	file multipart.File,
	filename string,
	gender string,
	result interface{},
) error {
	// Upload image to Supabase
	imagePath, err := s.storageService.UploadScanImage(userID, file, filename, "body")
	if err != nil {
		log.Printf("[ScanHistoryService] Failed to upload body scan image: %v", err)
		return fmt.Errorf("failed to upload image: %w", err)
	}

	// Count existing scans to generate name
	count, err := s.bodyScanRepo.CountByUserID(userID)
	if err != nil {
		log.Printf("[ScanHistoryService] Failed to count body scans: %v", err)
		return fmt.Errorf("failed to count scans: %w", err)
	}

	scanName := fmt.Sprintf("ScanTubuh_%d", count+1)

	var bodyType string
	var styleRecommendations []string
	var confidence float64

	// Type assertion based on gender
	switch gender {
	case "woman":
		if r, ok := result.(*dtos.WomanBodyScanPredictionResponseDTO); ok {
			bodyType = r.PredictedClass
			styleRecommendations = r.StyleRecommendations
			confidence = r.Confidence
		}
	case "man":
		if r, ok := result.(*dtos.MenBodyScanPredictionResponseDTO); ok {
			bodyType = r.PredictedClass
			styleRecommendations = r.StyleRecommendations
			confidence = r.Confidence
		}
	default:
		return fmt.Errorf("invalid gender: %s", gender)
	}

	// Create history record
	history := &models.BodyScanHistory{
		UserID:               userID,
		ScanName:             scanName,
		BodyType:             bodyType,
		Gender:               gender,
		ImagePath:            imagePath,
		StyleRecommendations: styleRecommendations,
		Confidence:           confidence,
	}

	if err := s.bodyScanRepo.Create(history); err != nil {
		log.Printf("[ScanHistoryService] Failed to save body scan history: %v", err)
		// Try to delete uploaded image on failure
		_ = s.storageService.DeleteFile(imagePath)
		return fmt.Errorf("failed to save history: %w", err)
	}

	log.Printf("[ScanHistoryService] Body scan history saved successfully: %s", scanName)
	return nil
}

// GetBodyScanHistories retrieves all body scan histories for a user
func (s *scanHistoryService) GetBodyScanHistories(userID uint64) (*dtos.BodyScanHistoryListDTO, error) {
	histories, err := s.bodyScanRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve histories: %w", err)
	}

	historyDTOs := make([]dtos.BodyScanHistoryDTO, 0, len(histories))
	for _, h := range histories {
		historyDTOs = append(historyDTOs, dtos.BodyScanHistoryDTO{
			ID:                   h.ID,
			ScanName:             h.ScanName,
			BodyType:             h.BodyType,
			Gender:               h.Gender,
			ImageURL:             h.ImagePath, // Return path, client will request via /history/image endpoint
			StyleRecommendations: h.StyleRecommendations,
			Confidence:           h.Confidence,
			CreatedAt:            h.CreatedAt,
		})
	}

	return &dtos.BodyScanHistoryListDTO{
		Histories: historyDTOs,
		Total:     int64(len(historyDTOs)),
	}, nil
}

// DeleteBodyScanHistory deletes a body scan history and its image
func (s *scanHistoryService) DeleteBodyScanHistory(userID, historyID uint64) error {
	// Get history to verify ownership and get image path
	history, err := s.bodyScanRepo.FindByID(historyID)
	if err != nil {
		return fmt.Errorf("history not found: %w", err)
	}

	// Verify ownership
	if history.UserID != userID {
		return fmt.Errorf("unauthorized: scan does not belong to user")
	}

	// Delete from database
	if err := s.bodyScanRepo.Delete(historyID); err != nil {
		return fmt.Errorf("failed to delete history: %w", err)
	}

	// Delete image from Supabase (best effort, don't fail if image deletion fails)
	if err := s.storageService.DeleteFile(history.ImagePath); err != nil {
		log.Printf("[ScanHistoryService] Warning: Failed to delete image %s: %v", history.ImagePath, err)
	}

	return nil
}
