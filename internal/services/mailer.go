package services

import (
	"Lejematch/config"
	"context"
	"errors"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// SendEmail sender en HTML-mail via Mailgun. Bruges af kontakt-mails og
// auth-mails (e-mail-bekræftelse, nulstil adgangskode).
func SendEmail(to, subject, html string) error {
	apiKey := config.AppConfigInstance.MailgunAPIKey
	domain := config.AppConfigInstance.MailgunDomain

	if apiKey == "" {
		return errors.New("MAILGUN_API_KEY not set in environment")
	}
	if domain == "" {
		return errors.New("MAILGUN_DOMAIN not set in environment")
	}

	mg := mailgun.NewMailgun(domain, apiKey)
	// Domænet er oprettet i Mailguns EU-region (mxa/mxb.eu.mailgun.org) —
	// uden dette rammer requests US-endpointet og fejler stille.
	mg.SetAPIBase(mailgun.APIBaseEU)

	sender := "LejeMatch <noreply@" + domain + ">"
	message := mg.NewMessage(sender, subject, "", to)
	message.SetHTML(html)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	return err
}

const emailLogoURL = "https://lejematch.dk/email-logo.png"

// EmailHeader er logoet der vises øverst i vores kundevendte mails
// (kontakt-notifikationer, velkomstmail m.fl.) — centreret, begrænset
// bredde så det ikke dominerer mailen.
func EmailHeader() string {
	return `<div style="text-align: center; margin-bottom: 24px;">
		<img src="` + emailLogoURL + `" alt="LejeMatch" style="max-width: 150px; height: auto;" />
	</div>`
}

// EmailSignature er den faste afslutning på vores kundevendte mails.
func EmailSignature() string {
	return `<p style="margin-top: 32px; color: #333;">
		Mvh,<br>
		LejeMatch<br>
		<a href="https://lejematch.dk" style="color: #006644;">lejematch.dk</a>
	</p>`
}
