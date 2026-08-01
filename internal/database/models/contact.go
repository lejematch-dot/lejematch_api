package models

import (
	"strconv"

	"gorm.io/gorm"
)

// ContactTargetType angiver om beskeden vedrører et Listing eller en SeekerListing.
type ContactTargetType string

const (
	ContactTargetListing ContactTargetType = "listing"
	ContactTargetSeeker  ContactTargetType = "seeker"
)

// ContactRelationshipType angiver om afsenderen er en del af et par eller en
// venneflok — kun relevant når NumPeople > 1, ellers tom streng.
type ContactRelationshipType string

const (
	ContactRelationshipCouple  ContactRelationshipType = "par"
	ContactRelationshipFriends ContactRelationshipType = "venner"
	ContactRelationshipOther   ContactRelationshipType = "andet"
)

// IsValid tillader også tom streng, da feltet kun er relevant når NumPeople > 1.
func (r ContactRelationshipType) IsValid() bool {
	switch r {
	case "", ContactRelationshipCouple, ContactRelationshipFriends, ContactRelationshipOther:
		return true
	default:
		return false
	}
}

func (r ContactRelationshipType) Label() string {
	switch r {
	case ContactRelationshipCouple:
		return "Par"
	case ContactRelationshipFriends:
		return "Venner"
	case ContactRelationshipOther:
		return "Andet"
	default:
		return ""
	}
}

// ContactAgeRange er afsenderens aldersinterval, oplyst af afsenderen selv
// ved kontakt — bruges til at give modtageren et hurtigt overblik, ikke til
// automatisk fravælgelse (se citynorm-lignende overvejelser i CLAUDE.md).
type ContactAgeRange string

const (
	ContactAgeUnder25 ContactAgeRange = "under25"
	ContactAge26To35  ContactAgeRange = "26-35"
	ContactAge35Plus  ContactAgeRange = "35+"
)

func (a ContactAgeRange) IsValid() bool {
	switch a {
	case ContactAgeUnder25, ContactAge26To35, ContactAge35Plus:
		return true
	default:
		return false
	}
}

func (a ContactAgeRange) Label() string {
	switch a {
	case ContactAgeUnder25:
		return "Under 25"
	case ContactAge26To35:
		return "26-35"
	case ContactAge35Plus:
		return "35+"
	default:
		return ""
	}
}

// ContactEmployment er afsenderens beskæftigelsessituation, oplyst af
// afsenderen selv ved kontakt.
type ContactEmployment string

const (
	ContactEmploymentFullTime ContactEmployment = "fast_job"
	ContactEmploymentStudent  ContactEmployment = "studerende"
	ContactEmploymentRetired  ContactEmployment = "pensionist"
	ContactEmploymentOther    ContactEmployment = "andet"
)

func (e ContactEmployment) IsValid() bool {
	switch e {
	case ContactEmploymentFullTime, ContactEmploymentStudent, ContactEmploymentRetired, ContactEmploymentOther:
		return true
	default:
		return false
	}
}

func (e ContactEmployment) Label() string {
	switch e {
	case ContactEmploymentFullTime:
		return "Fast job"
	case ContactEmploymentStudent:
		return "Studerende"
	case ContactEmploymentRetired:
		return "Pensionist"
	case ContactEmploymentOther:
		return "Andet"
	default:
		return ""
	}
}

// Contact er en besked sendt fra én bruger (Sender) til en anden (Recipient)
// om et konkret opslag. Gemmes i databasen ud over at blive sendt som e-mail,
// så modtageren kan se historikken i sit dashboard.
//
// NumPeople/RelationshipType/AgeRange/Employment oplyses af afsenderen selv
// og vises til modtageren som info i beskedlisten — bevidst IKKE eksponeret
// som filter/sortering i UI'en, da det ellers ville gøre LejeMatch til et
// aktivt værktøj til systematisk fravælgelse af lejere fremfor blot at
// understøtte modtagerens egen, individuelle vurdering.
//
// Kun udfyldt for TargetType=listing (en lejer der kontakter en udlejer om
// deres bolig) — det er lejeren der beskriver sig selv der. For
// TargetType=seeker (en udlejer der kontakter en lejers "søger bolig"-opslag)
// giver felterne ikke mening for afsenderen og står tomme/0.
type Contact struct {
	gorm.Model

	SenderID         uint              `gorm:"not null;index"`
	Sender           User              `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE"`
	RecipientID      uint              `gorm:"not null;index"`
	Recipient        User              `gorm:"foreignKey:RecipientID;constraint:OnDelete:CASCADE"`
	TargetType       ContactTargetType `gorm:"not null"`
	TargetID         uint              `gorm:"not null"`
	Message          string            `gorm:"not null"`
	SenderPhone      string
	NumPeople        int `gorm:"not null;default:1"`
	RelationshipType ContactRelationshipType
	AgeRange         ContactAgeRange   `gorm:"not null"`
	Employment       ContactEmployment `gorm:"not null"`
}

// NumPeopleSummary formaterer antal personer + evt. par/venner-label til
// visning, fx "3 (Venner)" eller "1".
func (c Contact) NumPeopleSummary() string {
	summary := strconv.Itoa(c.NumPeople)
	if label := c.RelationshipType.Label(); label != "" {
		summary += " (" + label + ")"
	}
	return summary
}
