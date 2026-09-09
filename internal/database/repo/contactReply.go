package repo

import (
	"Lejematch/internal/database"
	"Lejematch/internal/database/models"
)

type ContactRepliesRepo struct {
	*GenericRepo[models.ContactReply]
}

func NewContactRepliesRepo() *ContactRepliesRepo {
	return &ContactRepliesRepo{NewGenericRepo[models.ContactReply](database.DB)}
}

func (r *ContactRepliesRepo) FindByContactID(contactID uint) ([]*models.ContactReply, error) {
	var replies []*models.ContactReply
	err := r.db.Where("contact_id = ?", contactID).Order("created_at asc").Find(&replies).Error
	return replies, err
}
