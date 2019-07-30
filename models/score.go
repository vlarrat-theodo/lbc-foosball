package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"log"
	"time"
)

type Score struct {
	ID          	uuid.UUID   `json:"id" db:"id"`
	CreatedAt   	time.Time	`json:"created_at" db:"created_at"`
	UpdatedAt   	time.Time   `json:"updated_at" db:"updated_at"`
	User1Id			string		`json:"user1_id" db:"user1_id"`
	User2Id			string		`json:"user2_id" db:"user2_id"`
	User1Points		uint    	`json:"user1_points" db:"user1_points"`
	User2Points		uint    	`json:"user2_points" db:"user2_points"`
	User1Sets		int    		`json:"user1_sets" db:"user1_sets"`
	User2Sets		int    		`json:"user2_sets" db:"user2_sets"`
}

func (s Score) String() string {
	jp, marshalError := json.Marshal(s)
	if marshalError != nil {
		log.Println(marshalError)
		return ""
	}
	return string(jp)
}

type Scores []Score

func (s Scores) String() string {
	jp, marshalError := json.Marshal(s)
	if marshalError != nil {
		log.Println(marshalError)
		return ""
	}
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (s *Score) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.User1Id, Name: "User1Id"},
		&validators.StringIsPresent{Field: s.User2Id, Name: "User2Id"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Score) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Score) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
