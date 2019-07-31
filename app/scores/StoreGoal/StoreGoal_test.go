package main

import (
	"github.com/gobuffalo/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vlarrat-theodo/lbc-foosball/models"
	"testing"
	"time"
)

// TestUpdateScoreUsersIncoherence tests updateScore function for users incoherence between goal and score submitted.
//
func TestUpdateScoreUsersIncoherence(t *testing.T) {
	var updateScoreError error
	var score models.Score

	assertHandler := assert.New(t)
	now := time.Now()
	initialUUID, _ := uuid.NewV4()

	initialScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    7,
		User2Points:    3,
		User1Sets:      2,
		User2Sets:      1,
		GoalsInBalance: 0,
	}

	bothDifferentGoal := goal{
		Scorer:   "user3",
		Opponent: "user4",
		Player:   "p1",
		Gamelle:  false,
	}

	score = initialScore
	updateScoreError = updateScore(&score, bothDifferentGoal)
	assertHandler.NotNil(updateScoreError, "Both users different between goal and score: updateScore function should raise an error")
	assertHandler.Equal(initialScore, score, "Both users different between goal and score: updateScore function should not modify score")

	firstDifferentCase1Goal := goal{
		Scorer:   "user2",
		Opponent: "user3",
		Player:   "p1",
		Gamelle:  false,
	}

	score = initialScore
	updateScoreError = updateScore(&score, firstDifferentCase1Goal)
	assertHandler.NotNil(updateScoreError, "First user different between goal and score (case 1): updateScore function should raise an error")
	assertHandler.Equal(initialScore, score, "First user different between goal and score (case 1): updateScore function should not modify score")

	firstDifferentCase2Goal := goal{
		Scorer:   "user3",
		Opponent: "user2",
		Player:   "p1",
		Gamelle:  false,
	}

	score = initialScore
	updateScoreError = updateScore(&score, firstDifferentCase2Goal)
	assertHandler.NotNil(updateScoreError, "First user different between goal and score (case 2): updateScore function should raise an error")
	assertHandler.Equal(initialScore, score, "First user different between goal and score (case 2): updateScore function should not modify score")

	secondDifferentCase1Goal := goal{
		Scorer:   "user1",
		Opponent: "user3",
		Player:   "p1",
		Gamelle:  false,
	}

	score = initialScore
	updateScoreError = updateScore(&score, secondDifferentCase1Goal)
	assertHandler.NotNil(updateScoreError, "Second user different between goal and score (case 1): updateScore function should raise an error")
	assertHandler.Equal(initialScore, score, "Second user different between goal and score (case 1): updateScore function should not modify score")

	secondDifferentCase2Goal := goal{
		Scorer:   "user3",
		Opponent: "user1",
		Player:   "p1",
		Gamelle:  false,
	}

	score = initialScore
	updateScoreError = updateScore(&score, secondDifferentCase2Goal)
	assertHandler.NotNil(updateScoreError, "Second user different between goal and score (case 2): updateScore function should raise an error")
	assertHandler.Equal(initialScore, score, "Second user different between goal and score (case 2): updateScore function should not modify score")

	bothSameCase1Goal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p1",
		Gamelle:  false,
	}

	score = initialScore
	updateScoreError = updateScore(&score, bothSameCase1Goal)
	assertHandler.Nil(updateScoreError, "Both users same between goal and score (case 1): updateScore function should not raise an error")
	assertHandler.NotEqual(initialScore, score, "Both users same between goal and score (case 1): updateScore function should modify score")

	bothSameCase2Goal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p1",
		Gamelle:  false,
	}

	score = initialScore
	updateScoreError = updateScore(&score, bothSameCase2Goal)
	assertHandler.Nil(updateScoreError, "Both users same between goal and score (case 2): updateScore function should not raise an error")
	assertHandler.NotEqual(initialScore, score, "Both users same between goal and score (case 2): updateScore function should modify score")

}

// TestUpdateScoreExistingPlayers tests updateScore function for existing and non existing players.
//
func TestUpdateScoreExistingPlayers(t *testing.T) {
	var updateScoreError error
	var score models.Score

	assertHandler := assert.New(t)
	now := time.Now()
	initialUUID, _ := uuid.NewV4()

	initialScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    7,
		User2Points:    3,
		User1Sets:      2,
		User2Sets:      1,
		GoalsInBalance: 0,
	}

	notExistingPlayerGoal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "zizou",
		Gamelle:  false,
	}

	score = initialScore
	updateScoreError = updateScore(&score, notExistingPlayerGoal)
	assertHandler.NotNil(updateScoreError, "Goal from not existing player: updateScore function should raise an error")
	assertHandler.Equal(initialScore, score, "Goal from not existing player: updateScore function should not modify score")

	existingPlayerGoal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p1",
		Gamelle:  false,
	}

	score = initialScore
	updateScoreError = updateScore(&score, existingPlayerGoal)
	assertHandler.Nil(updateScoreError, "Goal from existing player: updateScore function should not raise an error")
	assertHandler.NotEqual(initialScore, score, "Goal from existing player: updateScore function should modify score")

}

// TestUpdateScoreRegularGoal tests updateScore function for classic goals.
//
func TestUpdateScoreRegularGoal(t *testing.T) {
	var score models.Score

	assertHandler := assert.New(t)
	now := time.Now()
	initialUUID, _ := uuid.NewV4()

	initialScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    7,
		User2Points:    3,
		User1Sets:      2,
		User2Sets:      1,
		GoalsInBalance: 0,
	}

	firstUserGoal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p1",
		Gamelle:  false,
	}

	awaitedFirstGoalScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    8,
		User2Points:    3,
		User1Sets:      2,
		User2Sets:      1,
		GoalsInBalance: 0,
	}

	score = initialScore
	_ = updateScore(&score, firstUserGoal)
	assertHandler.Equal(awaitedFirstGoalScore, score, "Regular goal from user1 (not winning set): score not updated as expected")

	secondUserGoal := goal{
		Scorer:   "user2",
		Opponent: "user1",
		Player:   "p1",
		Gamelle:  false,
	}

	awaitedSecondGoalScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    7,
		User2Points:    4,
		User1Sets:      2,
		User2Sets:      1,
		GoalsInBalance: 0,
	}

	score = initialScore
	_ = updateScore(&score, secondUserGoal)
	assertHandler.Equal(awaitedSecondGoalScore, score, "Regular goal from user2 (not winning set): score not updated as expected")

}

// TestUpdateScoreWinningSet tests updateScore function for goals leading to set winning.
//
func TestUpdateScoreWinningSet(t *testing.T) {
	var score models.Score

	assertHandler := assert.New(t)
	now := time.Now()
	initialUUID, _ := uuid.NewV4()

	firstUserGoal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p1",
		Gamelle:  false,
	}

	secondUserGoal := goal{
		Scorer:   "user2",
		Opponent: "user1",
		Player:   "p1",
		Gamelle:  false,
	}

	firstUserToWinScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    9,
		User2Points:    4,
		User1Sets:      2,
		User2Sets:      1,
		GoalsInBalance: 0,
	}

	awaitedFirstUserToWinAfterFirstUserGoalScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    0,
		User2Points:    0,
		User1Sets:      3,
		User2Sets:      1,
		GoalsInBalance: 0,
	}

	score = firstUserToWinScore
	_ = updateScore(&score, firstUserGoal)
	assertHandler.Equal(awaitedFirstUserToWinAfterFirstUserGoalScore, score, "User1 winning set (without goals in balance): score not updated as expected")

	awaitedFirstUserToWinAfterSecondUserGoalScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    9,
		User2Points:    5,
		User1Sets:      2,
		User2Sets:      1,
		GoalsInBalance: 0,
	}

	score = firstUserToWinScore
	_ = updateScore(&score, secondUserGoal)
	assertHandler.Equal(awaitedFirstUserToWinAfterSecondUserGoalScore, score, "User2 scored while user1 about to win: score not updated as expected")

	secondUserToWinScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    5,
		User2Points:    9,
		User1Sets:      6,
		User2Sets:      7,
		GoalsInBalance: 0,
	}

	awaitedSecondUserToWinAfterFirstUserGoalScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    6,
		User2Points:    9,
		User1Sets:      6,
		User2Sets:      7,
		GoalsInBalance: 0,
	}

	score = secondUserToWinScore
	_ = updateScore(&score, firstUserGoal)
	assertHandler.Equal(awaitedSecondUserToWinAfterFirstUserGoalScore, score, "User1 scored while user2 about to win: score not updated as expected")

	awaitedSecondUserToWinAfterSecondUserGoalScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    0,
		User2Points:    0,
		User1Sets:      6,
		User2Sets:      8,
		GoalsInBalance: 0,
	}

	score = secondUserToWinScore
	_ = updateScore(&score, secondUserGoal)
	assertHandler.Equal(awaitedSecondUserToWinAfterSecondUserGoalScore, score, "User2 winning set (without goals in balance): score not updated as expected")

	firstUserToWinScoreWithGoalsInBalance := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    9,
		User2Points:    4,
		User1Sets:      2,
		User2Sets:      1,
		GoalsInBalance: 6,
	}

	awaitedFirstUserToWinAfterFirstUserGoalScoreWithGoalsInBalance := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    0,
		User2Points:    0,
		User1Sets:      3,
		User2Sets:      1,
		GoalsInBalance: 0,
	}

	score = firstUserToWinScoreWithGoalsInBalance
	_ = updateScore(&score, firstUserGoal)
	assertHandler.Equal(awaitedFirstUserToWinAfterFirstUserGoalScoreWithGoalsInBalance, score, "User1 winning set (with goals in balance): score not updated as expected")

	secondUserToWinScoreWithGoalsInBalance := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    5,
		User2Points:    9,
		User1Sets:      6,
		User2Sets:      7,
		GoalsInBalance: 2,
	}

	awaitedSecondUserToWinAfterSecondUserGoalScoreWithGoalsInBalance := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    0,
		User2Points:    0,
		User1Sets:      6,
		User2Sets:      8,
		GoalsInBalance: 0,
	}

	score = secondUserToWinScoreWithGoalsInBalance
	_ = updateScore(&score, secondUserGoal)
	assertHandler.Equal(awaitedSecondUserToWinAfterSecondUserGoalScoreWithGoalsInBalance, score, "User2 winning set (with goals in balance): score not updated as expected")

}

// TestUpdateScorePissetteCase tests updateScore function for goals in "pissette" case.
//
func TestUpdateScorePissetteCase(t *testing.T) {
	var score models.Score

	assertHandler := assert.New(t)
	now := time.Now()
	initialUUID, _ := uuid.NewV4()

	initialScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    5,
		User2Points:    4,
		User1Sets:      2,
		User2Sets:      1,
		GoalsInBalance: 0,
	}

	classicPlayerGoal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p1",
		Gamelle:  false,
	}

	score = initialScore
	_ = updateScore(&score, classicPlayerGoal)
	assertHandler.NotEqual(initialScore, score, "Classic player goal: score should be modified")

	pissettePlayerGoal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p9",
		Gamelle:  false,
	}

	score = initialScore
	_ = updateScore(&score, pissettePlayerGoal)
	assertHandler.Equal(initialScore, score, "Pissette player goal without gamelle: score should not be modified")

	pissettePlayerGamelleGoal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p9",
		Gamelle:  true,
	}

	score = initialScore
	_ = updateScore(&score, pissettePlayerGamelleGoal)
	assertHandler.Equal(initialScore, score, "Pissette player goal with gamelle: score should not be modified")

}

// TestUpdateScoreGamelleCase tests updateScore function for goals in "gamelle" case.
//
func TestUpdateScoreGamelleCase(t *testing.T) {
	var score models.Score

	assertHandler := assert.New(t)
	now := time.Now()
	initialUUID, _ := uuid.NewV4()

	classicGoal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p1",
		Gamelle:  false,
	}

	gamelleGoal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p1",
		Gamelle:  true,
	}

	initialClassicScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    3,
		User2Points:    2,
		User1Sets:      6,
		User2Sets:      5,
		GoalsInBalance: 0,
	}

	awaitedAfterClassicGoalClassicScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    4,
		User2Points:    2,
		User1Sets:      6,
		User2Sets:      5,
		GoalsInBalance: 0,
	}

	score = initialClassicScore
	_ = updateScore(&score, classicGoal)
	assertHandler.Equal(awaitedAfterClassicGoalClassicScore, score, "Classic goal: score not updated as expected")

	awaitedAfterGamelleGoalClassicScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    3,
		User2Points:    1,
		User1Sets:      6,
		User2Sets:      5,
		GoalsInBalance: 0,
	}

	score = initialClassicScore
	_ = updateScore(&score, gamelleGoal)
	assertHandler.Equal(awaitedAfterGamelleGoalClassicScore, score, "Gamelle goal (classic case): score not updated as expected")

	initialZeroScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    2,
		User2Points:    0,
		User1Sets:      6,
		User2Sets:      5,
		GoalsInBalance: 0,
	}

	awaitedAfterGamelleGoalZeroScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    2,
		User2Points:    -1,
		User1Sets:      6,
		User2Sets:      5,
		GoalsInBalance: 0,
	}

	score = initialZeroScore
	_ = updateScore(&score, gamelleGoal)
	assertHandler.Equal(awaitedAfterGamelleGoalZeroScore, score, "Gamelle goal (negative case): score not updated as expected")

}

// TestUpdateScoreDemiCase tests updateScore function for goals in "demi" case.
//
func TestUpdateScoreDemiCase(t *testing.T) {
	var score models.Score

	assertHandler := assert.New(t)
	now := time.Now()
	initialUUID, _ := uuid.NewV4()

	demiGoal := goal{
		Scorer:   "user2",
		Opponent: "user1",
		Player:   "p4",
		Gamelle:  false,
	}

	initialScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    5,
		User2Points:    4,
		User1Sets:      1,
		User2Sets:      2,
		GoalsInBalance: 0,
	}

	awaitedAfterDemiGoalScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    5,
		User2Points:    4,
		User1Sets:      1,
		User2Sets:      2,
		GoalsInBalance: 2,
	}

	score = initialScore
	_ = updateScore(&score, demiGoal)
	assertHandler.Equal(awaitedAfterDemiGoalScore, score, "Demi goal: score not updated as expected")

	classicGoal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p1",
		Gamelle:  false,
	}

	awaitedAfterDemiThenClassicGoalScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    7,
		User2Points:    4,
		User1Sets:      1,
		User2Sets:      2,
		GoalsInBalance: 0,
	}

	score = awaitedAfterDemiGoalScore
	_ = updateScore(&score, classicGoal)
	assertHandler.Equal(awaitedAfterDemiThenClassicGoalScore, score, "Classic goal after demi goal: score not updated as expected")

	awaitedAfterDemiThenDemiGoalScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    5,
		User2Points:    4,
		User1Sets:      1,
		User2Sets:      2,
		GoalsInBalance: 4,
	}

	score = awaitedAfterDemiGoalScore
	_ = updateScore(&score, demiGoal)
	assertHandler.Equal(awaitedAfterDemiThenDemiGoalScore, score, "Demi goal after demi goal: score not updated as expected")

	gamelleGoal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p1",
		Gamelle:  true,
	}

	awaitedAfterDemiThenGamelleGoalScore := models.Score{
		ID:             initialUUID,
		CreatedAt:      now,
		UpdatedAt:      now,
		User1Id:        "user1",
		User2Id:        "user2",
		User1Points:    5,
		User2Points:    3,
		User1Sets:      1,
		User2Sets:      2,
		GoalsInBalance: 2,
	}

	score = awaitedAfterDemiGoalScore
	_ = updateScore(&score, gamelleGoal)
	assertHandler.Equal(awaitedAfterDemiThenGamelleGoalScore, score, "Gamelle goal after demi goal: score not updated as expected")

	demiGamelleGoal := goal{
		Scorer:   "user1",
		Opponent: "user2",
		Player:   "p4",
		Gamelle:  true,
	}

	score = awaitedAfterDemiGoalScore
	_ = updateScore(&score, demiGamelleGoal)
	assertHandler.Equal(awaitedAfterDemiGoalScore, score, "Gamelle by demi goal: score should not be modified")

}
