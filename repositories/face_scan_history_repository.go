package repositories

import (
	"flicknfit_backend/models"
	"fmt"

	"gorm.io/gorm"
)

// FaceScanHistoryRepository handles database operations for face scan history
type FaceScanHistoryRepository interface {
	Create(history *models.FaceScanHistory) error
	FindByUserID(userID uint64) ([]models.FaceScanHistory, error)
	FindByID(id uint64) (*models.FaceScanHistory, error)
	Delete(id uint64) error
	CountByUserID(userID uint64) (int64, error)
}

type faceScanHistoryRepository struct {
	db *gorm.DB
}

// NewFaceScanHistoryRepository creates a new face scan history repository
func NewFaceScanHistoryRepository(db *gorm.DB) FaceScanHistoryRepository {
	return &faceScanHistoryRepository{db: db}
}

// Create creates a new face scan history record
func (r *faceScanHistoryRepository) Create(history *models.FaceScanHistory) error {
	return r.db.Create(history).Error
}

// FindByUserID retrieves all face scan histories for a user
func (r *faceScanHistoryRepository) FindByUserID(userID uint64) ([]models.FaceScanHistory, error) {
	var histories []models.FaceScanHistory
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// FindByID retrieves a face scan history by ID
func (r *faceScanHistoryRepository) FindByID(id uint64) (*models.FaceScanHistory, error) {
	var history models.FaceScanHistory
	err := r.db.First(&history, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("face scan history not found")
		}
		return nil, err
	}
	return &history, nil
}

// Delete deletes a face scan history record
func (r *faceScanHistoryRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.FaceScanHistory{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("face scan history not found")
	}
	return nil
}

// CountByUserID counts total face scans for a user
func (r *faceScanHistoryRepository) CountByUserID(userID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&models.FaceScanHistory{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
