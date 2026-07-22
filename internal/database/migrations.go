package database

import (
	"Lejematch/internal/citynorm"
	"Lejematch/internal/database/models"
)

func Migrate() {
	// Kør før AutoMigrate — hvis kolonnen allerede findes med NULL-rækker
	// (fra før den fik en not-null-default), kan AutoMigrate ellers ikke
	// stramme constraintet. Fejler harmløst hvis kolonnen slet ikke findes
	// endnu (frisk database), AutoMigrate opretter den så korrekt fra bunden.
	backfillNewsletterOptIn()

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

	err = DB.AutoMigrate(&models.SeekerListing{})
	if err != nil {
		println(err)
		return
	}

	err = DB.AutoMigrate(&models.Favorite{})
	if err != nil {
		println(err)
		return
	}

	err = DB.AutoMigrate(&models.Contact{})
	if err != nil {
		println(err)
		return
	}

	err = DB.AutoMigrate(&models.Report{})
	if err != nil {
		println(err)
		return
	}

	backfillSeekerCityDisplay()
	normalizeCities()
}

// backfillNewsletterOptIn sætter newsletter_opt_in til false for rækker
// hvor den står som NULL (fx fordi kolonnen blev tilføjet uden en
// not-null-default før dette blev rettet).
func backfillNewsletterOptIn() {
	if err := DB.Exec(`UPDATE users SET newsletter_opt_in = false WHERE newsletter_opt_in IS NULL`).Error; err != nil {
		println(err.Error())
	}
}

// backfillSeekerCityDisplay kopierer den oprindelige (endnu ikke
// normaliserede) City-tekst over i CityDisplay for rækker der mangler den —
// kun relevant lige efter CityDisplay-kolonnen er tilføjet. Kører før
// normalizeCities(), som ellers ville rense City først.
func backfillSeekerCityDisplay() {
	if err := DB.Exec(`UPDATE seeker_listings SET city_display = city WHERE city_display = '' OR city_display IS NULL`).Error; err != nil {
		println(err.Error())
	}
}

// normalizeCities ensretter allerede gemte bynavne (f.eks. "Århus",
// "Aarhus C" -> "Aarhus"), så by-filteret ikke viser samme by flere
// gange. Idempotent — kører sikkert ved hver opstart.
func normalizeCities() {
	normalizeTableCities("listings")
	normalizeTableCities("seeker_listings")
}

func normalizeTableCities(table string) {
	var cities []string
	if err := DB.Table(table).Distinct("city").Pluck("city", &cities).Error; err != nil {
		println(err.Error())
		return
	}

	for _, city := range cities {
		normalized := citynorm.Normalize(city)
		if normalized == "" || normalized == city {
			continue
		}
		if err := DB.Table(table).Where("city = ?", city).Update("city", normalized).Error; err != nil {
			println(err.Error())
		}
	}
}
