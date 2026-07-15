package services

import (
	"Lejematch/internal/database/repo"
	"log"
	"strconv"
	"time"
	_ "time/tzdata" // bager IANA-tidszonedatabasen ind i binaren — runtime-imaget (alpine) har den ikke installeret
)

var copenhagen *time.Location

func init() {
	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		// Uden dette kører den daglige digest på UTC i stedet for dansk tid
		// uden nogen synlig fejl — så det skal larme, ikke bare falde tilbage.
		log.Printf("contactDigest: could not load Europe/Copenhagen, falling back to UTC: %v", err)
		loc = time.UTC
	}
	copenhagen = loc
}

// StartDailyContactDigest starter en baggrundsrutine der hver aften kl. 22
// (dansk tid) sender kontakt@lejematch.dk en mail med antallet af kontakter
// oprettet den dag mellem lejere og udlejere — ingen navne, opslag eller
// beskeder, kun et samlet tal.
func StartDailyContactDigest() {
	go func() {
		for {
			now := time.Now().In(copenhagen)
			next := time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, copenhagen)
			if !next.After(now) {
				next = next.Add(24 * time.Hour)
			}
			time.Sleep(time.Until(next))
			sendDailyContactDigest()
		}
	}()
}

// contactDigestDays angiver hvor mange dage tilbage (inkl. i dag) den
// daglige mail viser antal kontakter for.
const contactDigestDays = 3

// TriggerContactDigestNow sender den daglige oversigt med det samme, uden at
// vente på kl. 22. Bruges af admin-endpointet til at sende en udeblevet mail
// manuelt, eller til at teste at afsendelsen virker.
func TriggerContactDigestNow() {
	sendDailyContactDigest()
}

func sendDailyContactDigest() {
	now := time.Now().In(copenhagen)
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, copenhagen)

	contactsRepo := repo.NewContactsRepo()

	rows := ""
	for i := 0; i < contactDigestDays; i++ {
		dayEnd := startOfToday.AddDate(0, 0, -i+1)
		dayStart := startOfToday.AddDate(0, 0, -i)
		if i == 0 {
			dayEnd = now // dagen er ikke slut endnu
		}

		count, err := contactsRepo.CountBetween(dayStart, dayEnd)
		if err != nil {
			log.Printf("daily contact digest: failed to count contacts for %s: %v", dayStart.Format("02-01-2006"), err)
			return
		}

		label := dayStart.Format("02-01-2006")
		switch i {
		case 0:
			label = "I dag (" + label + ")"
		case 1:
			label = "I går (" + label + ")"
		case 2:
			label = "I forgårs (" + label + ")"
		}

		rows += `<p><strong>` + label + `:</strong> ` + strconv.FormatInt(count, 10) + ` kontakter</p>`
	}

	subject := "Daglig oversigt: kontakter på LejeMatch"
	html := `
	<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<h2>Daglig oversigt</h2>
			<p>Antal kontakter oprettet mellem lejere og udlejere på LejeMatch:</p>
			` + rows + `
		</body>
	</html>
	`

	if err := SendEmail("kontakt@lejematch.dk", subject, html); err != nil {
		log.Printf("daily contact digest: failed to send email: %v", err)
	}
}
