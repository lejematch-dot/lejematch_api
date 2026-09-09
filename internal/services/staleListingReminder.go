package services

import (
	"Lejematch/config"
	"Lejematch/internal/database/repo"
	"log"
	"strconv"
	"time"
)

// staleListingThreshold — hvor længe et opslag skal have stået aktivt
// (eller siden sidste påmindelse), før det udløser en ny "stadig
// aktuelt?"-mail. 30 dage giver i praksis en månedlig kadence pr. opslag.
const staleListingThreshold = 30 * 24 * time.Hour

// StartStaleListingReminders starter en baggrundsrutine der hver dag kl. 09
// (dansk tid) tjekker for forældede opslag og sender påmindelser. Kører
// dagligt (ikke månedligt) så kadencen pr. opslag styres af
// staleListingThreshold, ikke af hvornår rutinen selv starter.
func StartStaleListingReminders() {
	go func() {
		for {
			now := time.Now().In(copenhagen)
			next := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, copenhagen)
			if !next.After(now) {
				next = next.Add(24 * time.Hour)
			}
			time.Sleep(time.Until(next))
			sendStaleListingReminders()
		}
	}()
}

// TriggerStaleListingRemindersNow sender påmindelser med det samme, uden at
// vente på kl. 09. Bruges af admin-endpointet til at teste at afsendelsen
// virker.
func TriggerStaleListingRemindersNow() {
	sendStaleListingReminders()
}

func sendStaleListingReminders() {
	threshold := time.Now().Add(-staleListingThreshold)
	usersRepo := repo.NewUsersRepo()

	listingsRepo := repo.NewListingsRepo()
	staleListings, err := listingsRepo.FindStaleActive(threshold)
	if err != nil {
		log.Printf("stale listing reminder: failed to load stale listings: %v", err)
	}
	for _, listing := range staleListings {
		owner, err := usersRepo.FindByID(int(listing.UserID))
		if err != nil {
			log.Printf("stale listing reminder: failed to load owner %d: %v", listing.UserID, err)
			continue
		}
		editLink := config.AppConfigInstance.FrontendURL + "/dashboard/listings/" + strconv.Itoa(int(listing.ID)) + "/edit"
		if err := sendStaleReminderEmail(owner.Email, owner.FirstName, listing.Title, editLink); err != nil {
			log.Printf("stale listing reminder: failed to send to %s: %v", owner.Email, err)
			continue
		}
		now := time.Now()
		if err := listingsRepo.UpdateFields(int(listing.ID), map[string]interface{}{"last_reminder_sent_at": now}); err != nil {
			log.Printf("stale listing reminder: failed to update last_reminder_sent_at for listing %d: %v", listing.ID, err)
		}
		time.Sleep(300 * time.Millisecond)
	}

	seekersRepo := repo.NewSeekersRepo()
	staleSeekers, err := seekersRepo.FindStaleActive(threshold)
	if err != nil {
		log.Printf("stale listing reminder: failed to load stale seeker posts: %v", err)
	}
	for _, seeker := range staleSeekers {
		owner, err := usersRepo.FindByID(int(seeker.UserID))
		if err != nil {
			log.Printf("stale listing reminder: failed to load owner %d: %v", seeker.UserID, err)
			continue
		}
		editLink := config.AppConfigInstance.FrontendURL + "/dashboard/seekers/" + strconv.Itoa(int(seeker.ID)) + "/edit"
		if err := sendStaleReminderEmail(owner.Email, owner.FirstName, seeker.Title, editLink); err != nil {
			log.Printf("stale listing reminder: failed to send to %s: %v", owner.Email, err)
			continue
		}
		now := time.Now()
		if err := seekersRepo.UpdateFields(int(seeker.ID), map[string]interface{}{"last_reminder_sent_at": now}); err != nil {
			log.Printf("stale listing reminder: failed to update last_reminder_sent_at for seeker %d: %v", seeker.ID, err)
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func sendStaleReminderEmail(email, firstName, title, editLink string) error {
	subject := "Er dit opslag stadig aktuelt? — " + title

	html := `
	<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			` + EmailHeader() + `
			<h2>Er dit opslag stadig aktuelt?</h2>
			<p>Hej ` + firstName + `,</p>
			<p>Dit opslag <strong>` + title + `</strong> har været aktivt på LejeMatch i et stykke tid.</p>
			<p>Hvis du stadig leder, er der ikke noget at gøre — opslaget bliver stående. Men hvis du allerede har fundet noget, vil vi meget gerne have at du opdaterer eller sletter opslaget, så det ikke optager plads for andre og fylder unødvendigt i søgeresultaterne.</p>
			<p><a href="` + editLink + `" style="background: #006644; color: #fff; padding: 10px 20px; text-decoration: none; font-weight: bold; display: inline-block;">Opdatér eller slet dit opslag</a></p>
			<p style="color: #666; font-size: 12px; margin-top: 24px;">
				Du kan altid administrere alle dine opslag under "Mine opslag" på dashboardet.
			</p>
			` + EmailSignature() + `
		</body>
	</html>
	`

	return SendEmail(email, subject, html)
}
