package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/vlarrat-theodo/lbc-foosball/db"
	"github.com/vlarrat-theodo/lbc-foosball/models"
	"net/http"
	"strings"
)

type goal struct {
	Scorer   string `json:"scorer"`
	Opponent string `json:"opponent"`
	Player   string `json:"player"`
	Gamelle  bool   `json:"gamelle"`
}

func (g goal) isPissette() (pissetteGoal bool) {
	return g.Player == "p9"
}

var authorizedPlayers = [...]string{"p1", "p2", "p3", "p4", "p5", "p6", "p7", "p8", "p9", "p10", "p11"}
var demiPlayers = [...]string{"p4", "p5", "p6", "p7", "p8"}

func errorResponse(errorMessage string, errorStatusCode int) (APIResponse events.APIGatewayProxyResponse, APIError error) {
	errorMessage = strings.ReplaceAll(errorMessage, "\"", "\\\"")
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf("{\"error\": \"%s\"}", errorMessage),
		StatusCode: errorStatusCode,
	}, nil
}

func checkPlayerExists(playerToCheck string) (existingPlayer bool) {
	for _, authorizedPlayer := range authorizedPlayers {
		if playerToCheck == authorizedPlayer {
			return true
		}
	}
	return false
}

func isPlayerDemi(playerToCheck string) (playerDemi bool) {
	for _, demiPlayer := range demiPlayers {
		if playerToCheck == demiPlayer {
			return true
		}
	}
	return false
}

func updateScore(scoreToUpdate *models.Score, newGoal goal) (updateScoreError error) {
	// Check that submitted goal and score correspond to same users
	if !((newGoal.Scorer == scoreToUpdate.User1Id && newGoal.Opponent == scoreToUpdate.User2Id) || (newGoal.Scorer == scoreToUpdate.User2Id && newGoal.Opponent == scoreToUpdate.User1Id)) {
		return errors.New("goal and score do not correspond to same users")
	}

	// Check that submitted goal player belongs to authorized values
	if !checkPlayerExists(newGoal.Player) {
		return fmt.Errorf(`submitted goal player "%s" does not exist`, newGoal.Player)
	}

	// Handle "pissette" case: nothing happens when goal is scored by player "p9"
	if newGoal.isPissette() {
		return nil
	}

	// Handle "gamelle" case: opponent loses 1 point and scorer scores no point
	if newGoal.Gamelle {
		// "gamelle" case has only effect when not scored from "demi" player
		if !isPlayerDemi(newGoal.Player) {
			scoreToUpdate.ScorePoints(newGoal.Opponent, -1)
		}
		return nil
	}

	// Handle "demi" case: add 2 points in balance when goal scored by midfielder
	if isPlayerDemi(newGoal.Player) {
		scoreToUpdate.GoalsInBalance += 2
		return nil
	}

	// Handle "goals_in_balance" case: add points in balance to scorer instead of only 1 point
	if scoreToUpdate.GoalsInBalance > 0 {
		scoreToUpdate.ScorePoints(newGoal.Scorer, scoreToUpdate.GoalsInBalance)
		scoreToUpdate.GoalsInBalance = 0
	} else { // Classic case
		scoreToUpdate.ScorePoints(newGoal.Scorer, 1)
	}

	// Handle end of sets (when one user turns 10 points)
	if scoreToUpdate.IsSetFinished() {
		scoreToUpdate.ChangeSet(newGoal.Scorer)
	}

	return nil
}

func normalizeScoreToJSON(scoreToNormalize models.Score) (JSONScore string) {
	var normalizeScored = ""

	normalizeScored += fmt.Sprintf("{\n")
	normalizeScored += fmt.Sprintf("\t\"%s\": {\n", scoreToNormalize.User1Id)
	normalizeScored += fmt.Sprintf("\t\t\"sets\": \"%d\"\n,", scoreToNormalize.User1Sets)
	normalizeScored += fmt.Sprintf("\t\t\"points\": \"%d\"\n", scoreToNormalize.User1Points)
	normalizeScored += fmt.Sprintf("\t},\n")
	normalizeScored += fmt.Sprintf("\t\"%s\": {\n", scoreToNormalize.User2Id)
	normalizeScored += fmt.Sprintf("\t\t\"sets\": \"%d\"\n,", scoreToNormalize.User2Sets)
	normalizeScored += fmt.Sprintf("\t\t\"points\": \"%d\"\n", scoreToNormalize.User2Points)
	normalizeScored += fmt.Sprintf("\t},\n")
	normalizeScored += fmt.Sprintf("\t\"goals_in_balance\": \"%d\"\n", scoreToNormalize.GoalsInBalance)
	normalizeScored += fmt.Sprintf("}\n")

	return normalizeScored
}

func handler(request events.APIGatewayProxyRequest) (APIResponse events.APIGatewayProxyResponse, APIError error) {
	var databaseConnection *pop.Connection
	var requestError, dbError, updateScoreError error
	var validateError *validate.Errors
	var submittedGoal = goal{}
	var goalScore = models.Score{}
	var databaseConnector = db.DatabaseConnector{}

	databaseConnection, dbError = databaseConnector.GetConnection()
	if dbError != nil {
		return errorResponse(fmt.Sprintf("Failed to connect to database: %s", dbError), http.StatusInternalServerError)
	}
	defer databaseConnection.Close()

	requestError = json.Unmarshal([]byte(request.Body), &submittedGoal)
	if requestError != nil {
		return errorResponse(fmt.Sprintf("Bad request body: %s", requestError), http.StatusBadRequest)
	}

	existingScoreQuery := databaseConnection.Where("user1_id = ? AND user2_id = ? OR user1_id = ? AND user2_id = ?", submittedGoal.Scorer, submittedGoal.Opponent, submittedGoal.Opponent, submittedGoal.Scorer)
	scoreAlreadyExists, dbError := existingScoreQuery.Exists(models.Score{})

	if dbError != nil {
		return errorResponse(fmt.Sprintf("Failed to connect to database: %s", dbError), http.StatusInternalServerError)
	}

	if scoreAlreadyExists {
		dbError = existingScoreQuery.First(&goalScore)
		if dbError != nil {
			return errorResponse(fmt.Sprintf("Failed to retrieve existing score: %s", dbError), http.StatusInternalServerError)
		}
	} else {
		goalScore.User1Id = submittedGoal.Scorer
		goalScore.User2Id = submittedGoal.Opponent
	}

	updateScoreError = updateScore(&goalScore, submittedGoal)

	if updateScoreError != nil {
		return errorResponse(fmt.Sprintf("Failed to create/update score: %s", updateScoreError), http.StatusInternalServerError)
	}

	validateError, dbError = databaseConnection.ValidateAndSave(&goalScore)

	if validateError != nil && len(validateError.Errors) != 0 {
		return errorResponse(fmt.Sprintf("Failed to create/update score: %s", validateError), http.StatusInternalServerError)
	}
	if dbError != nil {
		return errorResponse(fmt.Sprintf("Failed to create/update score: %s", dbError), http.StatusInternalServerError)
	}

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       normalizeScoreToJSON(goalScore),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
