package services

import (
	"Lejematch/config"
	"Lejematch/internal/database/repo"
	"Lejematch/internal/security"
	"errors"
	"log"
	"time"
)

const passwordResetTTL = 1 * time.Hour

var ErrAlreadyVerified = errors.New("email already verified")

// ResendVerification sender et nyt bekræftelseslink. Kaldes stille (ingen fejl
// til klienten) hvis brugeren ikke findes, for ikke at afsløre hvilke e-mails
// er registreret.
func ResendVerification(userRepo *repo.UsersRepo, email string) error {
	user, err := userRepo.GetByEmailWithPassword(email)
	if err != nil {
		return nil
	}
	if user.IsActive {
		return ErrAlreadyVerified
	}
	return sendVerificationEmail(user.ID, user.Email, user.FirstName)
}

// VerifyEmail aktiverer kontoen tilhørende et gyldigt bekræftelses-token, og
// sender derefter et velkomstbrev. Fejl i selve velkomstmailen må ikke fejle
// verificeringen — kontoen er aktiv uanset — men logges så en fejlkonfigureret
// mailer ikke er usynlig.
func VerifyEmail(userRepo *repo.UsersRepo, token string) error {
	userID, err := ParseActionToken(token, "verify_email")
	if err != nil {
		return err
	}
	if err := userRepo.UpdateFields(int(userID), map[string]interface{}{"is_active": true}); err != nil {
		return err
	}

	user, err := userRepo.FindByID(int(userID))
	if err != nil {
		return nil
	}
	if err := sendWelcomeEmail(user.Email, user.FirstName); err != nil {
		log.Printf("failed to send welcome email to %s: %v", user.Email, err)
	}
	return nil
}

// sendWelcomeEmail sender velkomstbrevet efter en vellykket e-mailbekræftelse.
func sendWelcomeEmail(email, firstName string) error {
	subject := "Velkommen til LejeMatch!"
	html := `
	<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto;">
			<div style="padding: 24px 24px 0;">
				` + EmailHeader() + `
			</div>
			<div style="padding: 0 24px 24px;">
				<h2 style="color: #006644;">Velkomstbrev fra skaberne bag LejeMatch</h2>

				<p>Kære ` + firstName + `</p>

				<p>TUSIND TAK, fordi du har oprettet dig som bruger på Danmarks nye, smarte boligportal.
				Vi håber, at vi kan hjælpe dig i lige netop den situation du står i. Hvis du ikke allerede
				har fanget konceptet, så får du her en lille guide:</p>

				<h3 style="color: #006644; margin-top: 32px;">Skal du udleje?</h3>
				<ol style="padding-left: 20px;">
					<li>Klik på &rsquo;Opret opslag&rsquo; på LejeMatch.dk.</li>
					<li>Udfyld felter, som er påkrævet.<br />Herunder en kort beskrivelse af lejemålet og hvad du ønsker af lejeren.</li>
					<li>Udfyld kontaktformularen, så din fremtidige lejer kan kontakte dig.</li>
					<li>Vedhæft min. 3 billeder af lejemålet (Vigtigt: Undervurdér ikke effekten af gode billeder).</li>
					<li>1, 2, 3 &ndash; nu er du klar til at poste din lejlighed (klap klap).</li>
					<li>Scroll gennem listen af potentielle lejere, som har oprettet en annonce under &rsquo;Lejere&rsquo;.</li>
				</ol>

				<h3 style="color: #006644; margin-top: 32px;">Skal du leje?</h3>
				<ol style="padding-left: 20px;">
					<li>Klik på &rsquo;Opret opslag&rsquo; på LejeMatch.dk.</li>
					<li>Udfyld de påkrævede felter.<br />Herunder en kort beskrivelse af dig/jer som lejere.
					Skriv hvem du er, hvad du laver, om du har erfaring med at bo ude osv. (Med andre ord: Sælg din case ;)).</li>
					<li>Udfyld kontaktformularen, så du kan kontaktes.</li>
					<li>Vedhæft minimum 3 billeder af dig/jer.</li>
					<li>Voila, så er du klar til at poste.</li>
					<li>Scroll gennem listen af boliger på &rsquo;Boliger&rsquo;-siden.</li>
				</ol>

				<p style="margin-top: 32px;">Endnu engang TAK for at oprette dig på LejeMatch &ndash; held og lykke fra hele holdet.</p>
				` + EmailSignature() + `
			</div>
		</body>
	</html>
	`
	return SendEmail(email, subject, html)
}

// RequestPasswordReset sender et nulstil-adgangskode-link. Returnerer om
// kontoen findes, så handleren bevidst kan vise en tydelig besked i stedet
// for den anbefalede uspecifikke besked (fravalgt efter aftale — se
// forgot_password.go for afvejningen).
func RequestPasswordReset(userRepo *repo.UsersRepo, email string) (bool, error) {
	user, err := userRepo.GetByEmailWithPassword(email)
	if err != nil {
		return false, nil
	}

	token, err := GenerateActionToken(user.ID, "reset_password", passwordResetTTL)
	if err != nil {
		return true, err
	}

	link := config.AppConfigInstance.FrontendURL + "/nulstil-adgangskode/" + token
	subject := "Nulstil din adgangskode hos LejeMatch"
	html := `
	<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<h2>Nulstil din adgangskode</h2>
			<p>Hej ` + user.FirstName + `,</p>
			<p>Klik på linket for at vælge en ny adgangskode:</p>
			<p><a href="` + link + `">Nulstil adgangskode</a></p>
			<p style="color: #666; font-size: 12px;">Linket udløber om 1 time. Har du ikke bedt om dette, kan du ignorere denne e-mail.</p>
		</body>
	</html>
	`
	return true, SendEmail(user.Email, subject, html)
}

// ResetPassword sætter en ny adgangskode ud fra et gyldigt reset-token.
func ResetPassword(userRepo *repo.UsersRepo, token, newPassword string) error {
	userID, err := ParseActionToken(token, "reset_password")
	if err != nil {
		return err
	}

	hashed, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return userRepo.UpdateFields(int(userID), map[string]interface{}{"password": hashed})
}
