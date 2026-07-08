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
	err := u.db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error
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

// HardDelete fjerner brugerens række permanent i stedet for GORMs
// almindelige soft delete (som blot sætter deleted_at). En rigtig SQL
// DELETE udløser databasens ON DELETE CASCADE-regler, så profil, opslag,
// favoritter og beskeder også fjernes — i overensstemmelse med
// privatlivspolitikkens løfte om permanent sletning ved kontosletning.
func (u *UsersRepo) HardDelete(user *models.User) error {
	return u.db.Unscoped().Delete(user).Error
}
