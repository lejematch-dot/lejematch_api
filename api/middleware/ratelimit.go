package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// LoginLimiter begrænser login-forsøg til 10 pr. minut pr. IP, så en bot ikke
// kan brute-force adgangskoder.
var LoginLimiter = limiter.New(limiter.Config{
	Max:        10,
	Expiration: 1 * time.Minute,
})

// RegisterLimiter begrænser oprettelse af nye konti til 5 pr. time pr. IP.
var RegisterLimiter = limiter.New(limiter.Config{
	Max:        5,
	Expiration: 1 * time.Hour,
})

// ForgotPasswordLimiter og ResendVerificationLimiter begrænser hver til 5 pr.
// time pr. IP, så de ikke kan bruges til at spamme andres indbakker med mails.
var ForgotPasswordLimiter = limiter.New(limiter.Config{
	Max:        5,
	Expiration: 1 * time.Hour,
})

var ResendVerificationLimiter = limiter.New(limiter.Config{
	Max:        5,
	Expiration: 1 * time.Hour,
})

// ReportLimiter begrænser rapporter til 20 pr. time pr. IP, så
// rapport-funktionen ikke selv kan bruges som et spam-værktøj.
var ReportLimiter = limiter.New(limiter.Config{
	Max:        20,
	Expiration: 1 * time.Hour,
})

// PublicUploadLimiter begrænser det uautentificerede upload-endpoint (brugt
// til profilbillede under oprettelse, før kontoen — og dermed en JWT —
// findes) til 20 pr. time pr. IP.
var PublicUploadLimiter = limiter.New(limiter.Config{
	Max:        20,
	Expiration: 1 * time.Hour,
})
