package services

import (
	"Lejematch/internal/database/models"
	"Lejematch/internal/database/repo"
	"log"
	"sort"
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
// (dansk tid) sender kontakt@lejematch.dk en mail med dagens kontakter,
// fordelt på udlejer- og lejer-opslag. Om søndagen sendes derudover et
// ugentligt recap lige efter den daglige mail.
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
			if time.Now().In(copenhagen).Weekday() == time.Sunday {
				sendWeeklyContactDigest()
			}
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

// TriggerWeeklyContactDigestNow sender det ugentlige recap med det samme,
// uden at vente til søndag. Samme formål som TriggerContactDigestNow.
func TriggerWeeklyContactDigestNow() {
	sendWeeklyContactDigest()
}

// targetBreakdown bygger HTML der viser antal kontakter fordelt på
// udlejer-opslag (Listing) og lejer-opslag (SeekerListing), med hvilke
// konkrete opslag der er blevet kontaktet, sorteret efter flest kontakter.
func targetBreakdown(contacts []*models.Contact) string {
	listingCounts := map[uint]int{}
	seekerCounts := map[uint]int{}
	for _, c := range contacts {
		switch c.TargetType {
		case models.ContactTargetListing:
			listingCounts[c.TargetID]++
		case models.ContactTargetSeeker:
			seekerCounts[c.TargetID]++
		}
	}

	listingsRepo := repo.NewListingsRepo()
	seekersRepo := repo.NewSeekersRepo()

	html := `<h3 style="margin-top: 24px;">Udlejer-opslag (` + strconv.Itoa(len(listingCounts)) + ` opslag kontaktet)</h3>`
	html += listBreakdown(listingCounts, func(id uint) string {
		listing, err := listingsRepo.FindByID(int(id))
		if err != nil {
			return "Opslag #" + strconv.FormatUint(uint64(id), 10) + " (slettet)"
		}
		return listing.Title
	})

	html += `<h3 style="margin-top: 24px;">Lejer-opslag (` + strconv.Itoa(len(seekerCounts)) + ` opslag kontaktet)</h3>`
	html += listBreakdown(seekerCounts, func(id uint) string {
		seeker, err := seekersRepo.FindByID(int(id))
		if err != nil {
			return "Opslag #" + strconv.FormatUint(uint64(id), 10) + " (slettet)"
		}
		return seeker.Title
	})

	return html
}

func listBreakdown(counts map[uint]int, titleFor func(uint) string) string {
	if len(counts) == 0 {
		return `<p style="color: #666;">Ingen kontakter.</p>`
	}

	type row struct {
		id    uint
		count int
	}
	rows := make([]row, 0, len(counts))
	for id, count := range counts {
		rows = append(rows, row{id, count})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].count != rows[j].count {
			return rows[i].count > rows[j].count
		}
		return rows[i].id < rows[j].id
	})

	html := `<ul>`
	for _, r := range rows {
		html += `<li>` + titleFor(r.id) + `: ` + strconv.Itoa(r.count) + ` kontakt(er)</li>`
	}
	html += `</ul>`
	return html
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

	todayContacts, err := contactsRepo.FindBetween(startOfToday, now)
	if err != nil {
		log.Printf("daily contact digest: failed to load today's contacts: %v", err)
		return
	}

	subject := "Daglig oversigt: kontakter på LejeMatch"
	html := `
	<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<h2>Daglig oversigt</h2>
			<p>Antal kontakter oprettet mellem lejere og udlejere på LejeMatch:</p>
			` + rows + `
			<p>Dagens kontakter, fordelt på opslag:</p>
			` + targetBreakdown(todayContacts) + `
		</body>
	</html>
	`

	if err := SendEmail("kontakt@lejematch.dk", subject, html); err != nil {
		log.Printf("daily contact digest: failed to send email: %v", err)
	}
}

func sendWeeklyContactDigest() {
	now := time.Now().In(copenhagen)
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, copenhagen)
	weekStart := startOfToday.AddDate(0, 0, -6) // inkl. i dag = 7 dage

	contactsRepo := repo.NewContactsRepo()

	weekContacts, err := contactsRepo.FindBetween(weekStart, now)
	if err != nil {
		log.Printf("weekly contact digest: failed to load week's contacts: %v", err)
		return
	}

	dayRows := ""
	for i := 6; i >= 0; i-- {
		dayStart := startOfToday.AddDate(0, 0, -i)
		dayEnd := dayStart.AddDate(0, 0, 1)
		if i == 0 {
			dayEnd = now
		}
		count, err := contactsRepo.CountBetween(dayStart, dayEnd)
		if err != nil {
			log.Printf("weekly contact digest: failed to count contacts for %s: %v", dayStart.Format("02-01-2006"), err)
			return
		}
		dayRows += `<p><strong>` + dayStart.Format("Mon 02-01-2006") + `:</strong> ` + strconv.FormatInt(count, 10) + ` kontakter</p>`
	}

	subject := "Ugentligt recap: kontakter på LejeMatch"
	html := `
	<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<h2>Ugentligt recap</h2>
			<p>I alt <strong>` + strconv.Itoa(len(weekContacts)) + `</strong> kontakter de sidste 7 dage (` + weekStart.Format("02-01-2006") + ` &ndash; ` + now.Format("02-01-2006") + `):</p>
			` + dayRows + `
			<p>Ugens kontakter, fordelt på opslag:</p>
			` + targetBreakdown(weekContacts) + `
		</body>
	</html>
	`

	if err := SendEmail("kontakt@lejematch.dk", subject, html); err != nil {
		log.Printf("weekly contact digest: failed to send email: %v", err)
	}
}
