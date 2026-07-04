package repo

import (
	"Lejematch/internal/database"
	"Lejematch/internal/database/models"
)

type UsersRepo struct {
	*GenericRepo[models.User]
}

func NewUsersRepo() *UsersRepo {
	return &UsersRepo{NewGenericRepo[models.User](database.DB)}
}

func (u *UsersRepo) GetByEmailWithPassword(email string) (*models.User, error) {
	var user models.User
	err := u.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (u *UsersRepo) FindByID(id int) (*models.User, error) {
	var user models.User
	err := u.db.Omit("password").Where("id = ?", id).First(&user).Error
	return &user, err
}

func (u *UsersRepo) FindByIDWithPassword(id int) (*models.User, error) {
	var user models.User
	err := u.db.Where("id = ?", id).First(&user).Error
	return &user, err
}
