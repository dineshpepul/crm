package repositories

import (
	"crm-app/backend/models"
	"errors"

	"gorm.io/gorm"
)

// GormContactRepository implements ContactRepository with GORM
type GormContactRepository struct {
	db *gorm.DB
}

// // NewContactRepository creates a new contact repository
// func NewContactRepository(db *gorm.DB) models.ContactRepository {
// 	return &GormContactRepository{db: db}
// }

// FindByID finds a contact by ID
func (r *GormContactRepository) FindByID(id int) (*models.Contact, error) {
	var contact models.Contact
	result := r.db.First(&contact, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &contact, nil
}

// List returns contacts with pagination
func (r *GormContactRepository) List(offset int, limit int) ([]models.Contact, error) {
	var contacts []models.Contact
	query := r.db

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

// FindByLead returns contacts for a specific lead
func (r *GormContactRepository) FindByLead(leadID int) ([]models.Contact, error) {
	var contacts []models.Contact
	if err := r.db.Where("lead_id = ?", leadID).Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

// Create creates a new contact
func (r *GormContactRepository) Create(contact *models.Contact) error {
	return r.db.Create(contact).Error
}

// Update updates an existing contact
func (r *GormContactRepository) Update(contact *models.Contact) error {
	return r.db.Save(contact).Error
}

// Delete deletes a contact
func (r *GormContactRepository) Delete(id int) error {
	return r.db.Delete(&models.Contact{}, id).Error
}

// Search searches for contacts
func (r *GormContactRepository) Search(query string) ([]models.Contact, error) {
	var contacts []models.Contact
	if err := r.db.Where("name LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%").Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}
