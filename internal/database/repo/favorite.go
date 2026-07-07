package repo

import (
	"Lejematch/internal/database"
	"Lejematch/internal/database/models"
)

type FavoritesRepo struct {
	*GenericRepo[models.Favorite]
}

func NewFavoritesRepo() *FavoritesRepo {
	return &FavoritesRepo{NewGenericRepo[models.Favorite](database.DB)}
}

func (r *FavoritesRepo) FindByUser(userID uint) ([]*models.Favorite, error) {
	var favorites []*models.Favorite
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&favorites).Error
	return favorites, err
}

func (r *FavoritesRepo) FindOne(userID uint, favoriteType string, favoriteID uint) (*models.Favorite, error) {
	var favorite models.Favorite
	err := r.db.Where("user_id = ? AND favorite_type = ? AND favorite_id = ?", userID, favoriteType, favoriteID).First(&favorite).Error
	return &favorite, err
}

func (r *FavoritesRepo) DeleteByUserAndTarget(userID uint, favoriteType string, favoriteID uint) error {
	return r.db.Where("user_id = ? AND favorite_type = ? AND favorite_id = ?", userID, favoriteType, favoriteID).Delete(&models.Favorite{}).Error
}
