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

// Score represents current status of foosball match between two users.
//
type Score struct {
	ID             uuid.UUID `json:"id" db:"id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	User1Id        string    `json:"user1_id" db:"user1_id"`
	User2Id        string    `json:"user2_id" db:"user2_id"`
	User1Points    int       `json:"user1_points" db:"user1_points"`
	User2Points    int       `json:"user2_points" db:"user2_points"`
	User1Sets      int       `json:"user1_sets" db:"user1_sets"`
	User2Sets      int       `json:"user2_sets" db:"user2_sets"`
	GoalsInBalance int       `json:"goals_in_balance" db:"goals_in_balance"`
}

// ScorePoints add points to submitted scorer.
//
func (s *Score) ScorePoints(scorerID string, pointsToAdd int) {
	switch scorerID {
	case s.User1Id:
		s.User1Points += pointsToAdd
	case s.User2Id:
		s.User2Points += pointsToAdd
	}
}

// IsSetFinished check if current set is finished.
//
func (s *Score) IsSetFinished() (finishedSet bool) {
	const pointsToWinSet int = 10

	return s.User1Points >= pointsToWinSet || s.User2Points >= pointsToWinSet
}

// ChangeSet add 1 set to the winner and set points and balance to 0.
//
func (s *Score) ChangeSet(winnerID string) {
	s.User1Points = 0
	s.User2Points = 0
	s.GoalsInBalance = 0

	switch winnerID {
	case s.User1Id:
		s.User1Sets++
	case s.User2Id:
		s.User2Sets++
	}
}

// String returns string representation of Score.
//
func (s Score) String() (scoreString string) {
	jp, marshalError := json.Marshal(s)
	if marshalError != nil {
		log.Println(marshalError)
		return ""
	}
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
//
func (s *Score) Validate(tx *pop.Connection) (validatorErrors *validate.Errors, validationError error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.User1Id, Name: "User1Id"},
		&validators.StringIsPresent{Field: s.User2Id, Name: "User2Id"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
//
func (s *Score) ValidateCreate(tx *pop.Connection) (validatorErrors *validate.Errors, validationError error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
//
func (s *Score) ValidateUpdate(tx *pop.Connection) (validatorErrors *validate.Errors, validationError error) {
	return validate.NewErrors(), nil
}
