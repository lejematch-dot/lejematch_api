package database

import "Lejematch/internal/database/models"

func Migrate() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		println(err)
		return
	}
	
	err = DB.AutoMigrate(&models.Profile{})
	if err != nil {
		println(err)
		return
	}

	err = DB.AutoMigrate(&models.Listing{})
	if err != nil {
		println(err)
		return
	}

}
