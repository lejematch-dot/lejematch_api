package services

import (
	"Lejematch/config"
	"Lejematch/internal/database/repo"
	"log"
	"time"
)

const newsletterUnsubscribeTTL = 365 * 24 * time.Hour

// UnsubscribeNewsletter slår NewsletterOptIn fra for den bruger tokenet
// tilhører. Bruges af det offentlige afmeldingslink i hver nyhedsbrev-mail.
func UnsubscribeNewsletter(userRepo *repo.UsersRepo, token string) error {
	userID, err := ParseActionToken(token, "unsubscribe_newsletter")
	if err != nil {
		return err
	}
	return userRepo.UpdateFields(int(userID), map[string]interface{}{"newsletter_opt_in": false})
}

// SubscribeNewsletter slår NewsletterOptIn til for den bruger tokenet
// tilhører. Bruges af invitationslinket der sendes til eksisterende
// brugere, som ikke fik mulighed for at tilmelde sig ved oprettelse.
func SubscribeNewsletter(userRepo *repo.UsersRepo, token string) error {
	userID, err := ParseActionToken(token, "subscribe_newsletter")
	if err != nil {
		return err
	}
	return userRepo.UpdateFields(int(userID), map[string]interface{}{"newsletter_opt_in": true})
}

// SendNewsletterInvite sender en engangsmail til alle aktive brugere der
// endnu ikke er tilmeldt nyhedsbreve, og spørger om de vil tilmelde sig.
// Idempotent i praksis — kør den igen for kun at ramme dem der stadig
// mangler at tage stilling.
func SendNewsletterInvite(userRepo *repo.UsersRepo) (sent int, failed int) {
	targets, err := userRepo.FindNewsletterInviteTargets()
	if err != nil {
		log.Printf("newsletter invite: failed to load targets: %v", err)
		return 0, 0
	}

	subject := "Vil du modtage nyhedsbreve fra LejeMatch?"

	for _, user := range targets {
		token, err := GenerateActionToken(user.ID, "subscribe_newsletter", newsletterUnsubscribeTTL)
		if err != nil {
			log.Printf("newsletter invite: failed to generate token for %s: %v", user.Email, err)
			failed++
			continue
		}

		subscribeLink := config.AppConfigInstance.FrontendURL + "/tilmeld-nyhedsbrev/" + token
		html := `
		<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<h2>Vil du høre fra os, når der sker noget nyt?</h2>
				<p>Hej ` + user.FirstName + `,</p>
				<p>
					Vi er begyndt at sende nyhedsbreve ud, fx når der kommer nye boliger op på LejeMatch.
					Da du allerede har en profil hos os, vil vi gerne spørge om du har lyst til at modtage dem.
				</p>
				<p><a href="` + subscribeLink + `">Ja tak, tilmeld mig nyhedsbreve</a></p>
				<p style="color: #666; font-size: 12px;">
					Gør du ingenting, sker der ingenting — vi tilmelder ingen uden aktivt samtykke.
				</p>
			</body>
		</html>
		`

		if err := SendEmail(user.Email, subject, html); err != nil {
			log.Printf("newsletter invite: failed to send to %s: %v", user.Email, err)
			failed++
			continue
		}
		sent++

		time.Sleep(300 * time.Millisecond)
	}

	return sent, failed
}

// SendNewsletter sender en mail til alle brugere der har sagt ja tak til
// nyhedsbreve. Hver mail får sin egen afmeldingslink tilføjet automatisk,
// så indholdet (subject/html) ikke selv behøver håndtere det. Kører
// synkront i kalderens goroutine — kald denne fra en baggrunds-goroutine
// hvis den bruges fra en HTTP-handler, så requesten ikke venter på alle sends.
func SendNewsletter(userRepo *repo.UsersRepo, subject, html string) (sent int, failed int) {
	recipients, err := userRepo.FindNewsletterRecipients()
	if err != nil {
		log.Printf("newsletter: failed to load recipients: %v", err)
		return 0, 0
	}

	for _, user := range recipients {
		token, err := GenerateActionToken(user.ID, "unsubscribe_newsletter", newsletterUnsubscribeTTL)
		if err != nil {
			log.Printf("newsletter: failed to generate unsubscribe token for %s: %v", user.Email, err)
			failed++
			continue
		}

		unsubscribeLink := config.AppConfigInstance.FrontendURL + "/afmeld-nyhedsbrev/" + token
		fullHTML := html + `
		<hr style="margin-top: 32px; border: none; border-top: 1px solid #e5e5e5;">
		<p style="color: #999; font-size: 11px; margin-top: 16px;">
			Du modtager denne mail fordi du har tilmeldt dig nyhedsbreve fra LejeMatch.
			<a href="` + unsubscribeLink + `" style="color: #999;">Afmeld nyhedsbreve</a>.
		</p>
		`

		if err := SendEmail(user.Email, subject, fullHTML); err != nil {
			log.Printf("newsletter: failed to send to %s: %v", user.Email, err)
			failed++
			continue
		}
		sent++

		time.Sleep(300 * time.Millisecond)
	}

	return sent, failed
}
