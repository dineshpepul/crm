package repositories

import (
	"crm-app/backend/models"

	"gorm.io/gorm"
)

// type gormUserRepository struct {
// 	db *gorm.DB
// }

// // NewUserRepository creates a new user repository
// func NewUserRepository(db *gorm.DB) models.UserRepository {
// 	return &gormUserRepository{db: db}
// }

// FindByID finds a user by ID
func (r *gormUserRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *gormUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (r *gormUserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update updates a user
func (r *gormUserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user
func (r *gormUserRepository) Delete(id int) error {
	return r.db.Delete(&models.User{}, id).Error
}

// List returns all users
func (r *gormUserRepository) List() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}
