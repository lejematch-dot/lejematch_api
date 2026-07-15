package repo

import (
	"Lejematch/internal/database"
	"Lejematch/internal/database/models"
	"time"
)

type ContactsRepo struct {
	*GenericRepo[models.Contact]
}

func NewContactsRepo() *ContactsRepo {
	return &ContactsRepo{NewGenericRepo[models.Contact](database.DB)}
}

func (r *ContactsRepo) FindByRecipient(userID uint) ([]*models.Contact, error) {
	var contacts []*models.Contact
	err := r.db.Where("recipient_id = ?", userID).Order("created_at desc").Find(&contacts).Error
	return contacts, err
}

// CountBetween tæller antal kontakter oprettet i perioden [from, to).
func (r *ContactsRepo) CountBetween(from, to time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.Contact{}).Where("created_at >= ? AND created_at < ?", from, to).Count(&count).Error
	return count, err
}
