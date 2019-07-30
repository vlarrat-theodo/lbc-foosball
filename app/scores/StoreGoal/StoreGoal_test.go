package main

import (
	"github.com/gobuffalo/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vlarrat-theodo/lbc-foosball/models"
	"testing"
	"time"
)

func TestUpdateScoreUsersIncoherence(t *testing.T) {
	var updateScoreError error
	var score models.Score

	assertHandler := assert.New(t)
	now := time.Now()
	initialUUID, _ := uuid.NewV4()

	initialScore := models.Score{
		ID: initialUUID,
		CreatedAt: now,
		UpdatedAt: now,
		User1Id: "user1",
		User2Id: "user2",
		User1Points: 7,
		User2Points: 3,
		User1Sets: 2,
		User2Sets: 1,
	}

	bothDifferentGoal := Goal{
		Scorer: "user3",
		Opponent: "user4",
	}

	firstDifferentCase1Goal := Goal{
		Scorer: "user2",
		Opponent: "user3",
	}

	firstDifferentCase2Goal := Goal{
		Scorer: "user3",
		Opponent: "user2",
	}

	secondDifferentCase1Goal := Goal{
		Scorer: "user1",
		Opponent: "user3",
	}

	secondDifferentCase2Goal := Goal{
		Scorer: "user3",
		Opponent: "user1",
	}

	bothSameCase1Goal := Goal{
		Scorer: "user1",
		Opponent: "user2",
	}

	bothSameCase2Goal := Goal{
		Scorer: "user1",
		Opponent: "user2",
	}

	score = initialScore
	updateScoreError = updateScore(&score, &bothDifferentGoal)
	assertHandler.NotNil(updateScoreError, "Both users different between goal and score: updateScore function should raise an error")
	assertHandler.Equal(initialScore, score, "Both users different between goal and score: updateScore function should not modify score")

	score = initialScore
	updateScoreError = updateScore(&score, &firstDifferentCase1Goal)
	assertHandler.NotNil(updateScoreError, "First user different between goal and score (case 1): updateScore function should raise an error")
	assertHandler.Equal(initialScore, score, "First user different between goal and score (case 1): updateScore function should not modify score")

	score = initialScore
	updateScoreError = updateScore(&score, &firstDifferentCase2Goal)
	assertHandler.NotNil(updateScoreError, "First user different between goal and score (case 2): updateScore function should raise an error")
	assertHandler.Equal(initialScore, score, "First user different between goal and score (case 2): updateScore function should not modify score")

	score = initialScore
	updateScoreError = updateScore(&score, &secondDifferentCase1Goal)
	assertHandler.NotNil(updateScoreError, "Second user different between goal and score (case 1): updateScore function should raise an error")
	assertHandler.Equal(initialScore, score, "Second user different between goal and score (case 1): updateScore function should not modify score")

	score = initialScore
	updateScoreError = updateScore(&score, &secondDifferentCase2Goal)
	assertHandler.NotNil(updateScoreError, "Second user different between goal and score (case 2): updateScore function should raise an error")
	assertHandler.Equal(initialScore, score, "Second user different between goal and score (case 2): updateScore function should not modify score")

	score = initialScore
	updateScoreError = updateScore(&score, &bothSameCase1Goal)
	assertHandler.Nil(updateScoreError, "Both users same between goal and score (case 1): updateScore function should not raise an error")

	score = initialScore
	updateScoreError = updateScore(&score, &bothSameCase2Goal)
	assertHandler.Nil(updateScoreError, "Both users same between goal and score (case 2): updateScore function should not raise an error")

}

func TestUpdateScoreRegularGoal(t *testing.T) {
	var score models.Score

	assertHandler := assert.New(t)
	now := time.Now()
	initialUUID, _ := uuid.NewV4()

	initialScore := models.Score{
		ID: initialUUID,
		CreatedAt: now,
		UpdatedAt: now,
		User1Id: "user1",
		User2Id: "user2",
		User1Points: 7,
		User2Points: 3,
		User1Sets: 2,
		User2Sets: 1,
	}

	firstUserGoal := Goal{
		Scorer: "user1",
		Opponent: "user2",
	}

	awaitedFirstGoalScore := models.Score{
		ID: initialUUID,
		CreatedAt: now,
		UpdatedAt: now,
		User1Id: "user1",
		User2Id: "user2",
		User1Points: 8,
		User2Points: 3,
		User1Sets: 2,
		User2Sets: 1,
	}

	secondUserGoal := Goal{
		Scorer: "user2",
		Opponent: "user1",
	}

	awaitedSecondGoalScore := models.Score{
		ID: initialUUID,
		CreatedAt: now,
		UpdatedAt: now,
		User1Id: "user1",
		User2Id: "user2",
		User1Points: 7,
		User2Points: 4,
		User1Sets: 2,
		User2Sets: 1,
	}

	score = initialScore
	_ = updateScore(&score, &firstUserGoal)
	assertHandler.Equal(awaitedFirstGoalScore, score, "Regular goal from user1 (not winning set): score not updated as expected")

	score = initialScore
	_ = updateScore(&score, &secondUserGoal)
	assertHandler.Equal(awaitedSecondGoalScore, score, "Regular goal from user2 (not winning set): score not updated as expected")

}

func TestUpdateScoreWinningSet(t *testing.T) {
	var score models.Score

	assertHandler := assert.New(t)
	now := time.Now()
	initialUUID, _ := uuid.NewV4()

	firstUserGoal := Goal{
		Scorer: "user1",
		Opponent: "user2",
	}

	secondUserGoal := Goal{
		Scorer: "user2",
		Opponent: "user1",
	}

	firstUserToWinScore := models.Score{
		ID: initialUUID,
		CreatedAt: now,
		UpdatedAt: now,
		User1Id: "user1",
		User2Id: "user2",
		User1Points: 9,
		User2Points: 4,
		User1Sets: 2,
		User2Sets: 1,
	}

	secondUserToWinScore := models.Score{
		ID: initialUUID,
		CreatedAt: now,
		UpdatedAt: now,
		User1Id: "user1",
		User2Id: "user2",
		User1Points: 5,
		User2Points: 9,
		User1Sets: 6,
		User2Sets: 7,
	}

	awaitedFirstUserToWinAfterFirstUserGoalScore := models.Score{
		ID: initialUUID,
		CreatedAt: now,
		UpdatedAt: now,
		User1Id: "user1",
		User2Id: "user2",
		User1Points: 0,
		User2Points: 0,
		User1Sets: 3,
		User2Sets: 1,
	}

	awaitedFirstUserToWinAfterSecondUserGoalScore := models.Score{
		ID: initialUUID,
		CreatedAt: now,
		UpdatedAt: now,
		User1Id: "user1",
		User2Id: "user2",
		User1Points: 9,
		User2Points: 5,
		User1Sets: 2,
		User2Sets: 1,
	}

	awaitedSecondUserToWinAfterFirstUserGoalScore := models.Score{
		ID: initialUUID,
		CreatedAt: now,
		UpdatedAt: now,
		User1Id: "user1",
		User2Id: "user2",
		User1Points: 6,
		User2Points: 9,
		User1Sets: 6,
		User2Sets: 7,
	}

	awaitedSecondUserToWinAfterSecondUserGoalScore := models.Score{
		ID: initialUUID,
		CreatedAt: now,
		UpdatedAt: now,
		User1Id: "user1",
		User2Id: "user2",
		User1Points: 0,
		User2Points: 0,
		User1Sets: 6,
		User2Sets: 8,
	}

	score = firstUserToWinScore
	_ = updateScore(&score, &firstUserGoal)
	assertHandler.Equal(awaitedFirstUserToWinAfterFirstUserGoalScore, score, "User1 winning set: score not updated as expected")

	score = firstUserToWinScore
	_ = updateScore(&score, &secondUserGoal)
	assertHandler.Equal(awaitedFirstUserToWinAfterSecondUserGoalScore, score, "User2 scored while user1 about to win: score not updated as expected")

	score = secondUserToWinScore
	_ = updateScore(&score, &firstUserGoal)
	assertHandler.Equal(awaitedSecondUserToWinAfterFirstUserGoalScore, score, "User1 scored while user2 about to win: score not updated as expected")

	score = secondUserToWinScore
	_ = updateScore(&score, &secondUserGoal)
	assertHandler.Equal(awaitedSecondUserToWinAfterSecondUserGoalScore, score, "User2 winning set: score not updated as expected")

}
