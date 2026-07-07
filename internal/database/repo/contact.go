package repo

import (
	"Lejematch/internal/database"
	"Lejematch/internal/database/models"
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
