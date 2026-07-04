package repo

import (
	"Lejematch/internal/database"
	"Lejematch/internal/database/models"
)

type ProfilesRepo struct {
	*GenericRepo[models.Profile]
}

func NewProfilesRepo() *ProfilesRepo {
	return &ProfilesRepo{NewGenericRepo[models.Profile](database.DB)}
}

func (p *ProfilesRepo) FindByUserID(userID uint) (*models.Profile, error) {
	var profile models.Profile
	err := p.db.Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}

func (p *ProfilesRepo) UpdateByUserID(userID uint, fields map[string]interface{}) error {
	return p.db.Model(&models.Profile{}).Where("user_id = ?", userID).Updates(fields).Error
}

