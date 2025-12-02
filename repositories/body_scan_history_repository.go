package repositories

import (
	"flicknfit_backend/models"
	"fmt"

	"gorm.io/gorm"
)

// BodyScanHistoryRepository handles database operations for body scan history
type BodyScanHistoryRepository interface {
	Create(history *models.BodyScanHistory) error
	FindByUserID(userID uint64) ([]models.BodyScanHistory, error)
	FindByID(id uint64) (*models.BodyScanHistory, error)
	Delete(id uint64) error
	CountByUserID(userID uint64) (int64, error)
}

type bodyScanHistoryRepository struct {
	db *gorm.DB
}

// NewBodyScanHistoryRepository creates a new body scan history repository
func NewBodyScanHistoryRepository(db *gorm.DB) BodyScanHistoryRepository {
	return &bodyScanHistoryRepository{db: db}
}

// Create creates a new body scan history record
func (r *bodyScanHistoryRepository) Create(history *models.BodyScanHistory) error {
	return r.db.Create(history).Error
}

// FindByUserID retrieves all body scan histories for a user
func (r *bodyScanHistoryRepository) FindByUserID(userID uint64) ([]models.BodyScanHistory, error) {
	var histories []models.BodyScanHistory
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// FindByID retrieves a body scan history by ID
func (r *bodyScanHistoryRepository) FindByID(id uint64) (*models.BodyScanHistory, error) {
	var history models.BodyScanHistory
	err := r.db.First(&history, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("body scan history not found")
		}
		return nil, err
	}
	return &history, nil
}

// Delete deletes a body scan history record
func (r *bodyScanHistoryRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.BodyScanHistory{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("body scan history not found")
	}
	return nil
}

// CountByUserID counts total body scans for a user
func (r *bodyScanHistoryRepository) CountByUserID(userID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&models.BodyScanHistory{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
