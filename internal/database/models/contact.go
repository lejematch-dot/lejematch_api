package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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

// IntSlice gemmes som jsonb — bruges til Ages, da antallet af aldre varierer
// med NumPeople.
type IntSlice []int

func (s IntSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *IntSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan IntSlice")
	}
	return json.Unmarshal(bytes, s)
}

// Contact er en besked sendt fra én bruger (Sender) til en anden (Recipient)
// om et konkret opslag. Gemmes i databasen ud over at blive sendt som e-mail,
// så modtageren kan se historikken i sit dashboard.
//
// NumPeople/RelationshipType/Ages/Employment/HasPets oplyses af afsenderen
// selv og vises til modtageren som info i beskedlisten — bevidst IKKE
// eksponeret som filter/sortering i UI'en, da det ellers ville gøre
// LejeMatch til et aktivt værktøj til systematisk fravælgelse af lejere
// fremfor blot at understøtte modtagerens egen, individuelle vurdering.
//
// Kun udfyldt for TargetType=listing (en lejer der kontakter en udlejer om
// deres bolig) — det er lejeren der beskriver sig selv der. For
// TargetType=seeker (en udlejer der kontakter en lejers "søger bolig"-opslag)
// giver felterne ikke mening for afsenderen og står tomme/0/false.
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
	// Ages har én alder per person i NumPeople (maks. 5 — dropdown-loftet).
	Ages       IntSlice          `gorm:"type:jsonb"`
	Employment ContactEmployment `gorm:"not null"`
	HasPets    bool
}

// NumPeopleSummary formaterer antal personer + evt. par/venner-label til
// visning, fx "3 (Venner)" eller "5+". Dropdown-feltet i UI'en stopper ved
// "5+", så 5 betyder reelt "5 eller flere".
func (c Contact) NumPeopleSummary() string {
	summary := strconv.Itoa(c.NumPeople)
	if c.NumPeople >= 5 {
		summary = "5+"
	}
	if label := c.RelationshipType.Label(); label != "" {
		summary += " (" + label + ")"
	}
	return summary
}

// AgesSummary formaterer aldrene til visning, fx "28" eller "28, 31 og 34".
func (c Contact) AgesSummary() string {
	if len(c.Ages) == 0 {
		return ""
	}
	if len(c.Ages) == 1 {
		return strconv.Itoa(c.Ages[0])
	}
	parts := make([]string, len(c.Ages))
	for i, age := range c.Ages {
		parts[i] = strconv.Itoa(age)
	}
	summary := parts[0]
	for i := 1; i < len(parts)-1; i++ {
		summary += ", " + parts[i]
	}
	summary += " og " + parts[len(parts)-1]
	return summary
}

func (c Contact) PetsLabel() string {
	if c.HasPets {
		return "Ja"
	}
	return "Nej"
}
