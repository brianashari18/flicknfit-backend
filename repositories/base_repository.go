package repositories

import (
	"gorm.io/gorm"
)

// BaseRepository provides a common structure for all repositories.
// It holds the GORM database instance.
type BaseRepository struct {
	DB *gorm.DB
}
