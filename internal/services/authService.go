package services

import (
	"Lejematch/config"
	"Lejematch/internal/database/repo"
	"Lejematch/internal/security"
	"errors"
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

// VerifyEmail aktiverer kontoen tilhørende et gyldigt bekræftelses-token.
func VerifyEmail(userRepo *repo.UsersRepo, token string) error {
	userID, err := ParseActionToken(token, "verify_email")
	if err != nil {
		return err
	}
	return userRepo.UpdateFields(int(userID), map[string]interface{}{"is_active": true})
}

// RequestPasswordReset sender et nulstil-adgangskode-link. Svarer altid succes
// til klienten (håndteres i handleren) uanset om e-mailen findes, for ikke at
// afsløre registrerede e-mails.
func RequestPasswordReset(userRepo *repo.UsersRepo, email string) error {
	user, err := userRepo.GetByEmailWithPassword(email)
	if err != nil {
		return nil
	}

	token, err := GenerateActionToken(user.ID, "reset_password", passwordResetTTL)
	if err != nil {
		return err
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
	return SendEmail(user.Email, subject, html)
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
