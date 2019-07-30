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


type Goal struct {
	Scorer		string `json:"scorer"`
	Opponent	string `json:"opponent"`
}


func errorResponse(errorMessage string, errorStatusCode int) (events.APIGatewayProxyResponse, error) {
	errorMessage = strings.ReplaceAll(errorMessage, "\"", "\\\"")
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf("{\"error\": \"%s\"}", errorMessage),
		StatusCode: errorStatusCode,
	}, nil
}


func updateScore(scoreToUpdate *models.Score, newGoal *Goal)  error {
	var scorerPointsToAdd, opponentPointsToAdd uint

	// Check that submitted goal and score correspond to same users
	if !((newGoal.Scorer == scoreToUpdate.User1Id && newGoal.Opponent == scoreToUpdate.User2Id) || (newGoal.Scorer == scoreToUpdate.User2Id && newGoal.Opponent == scoreToUpdate.User1Id)){
		return errors.New("goal and score do not correspond to same users")
	}

	scorerPointsToAdd = 1
	opponentPointsToAdd = 0

	if newGoal.Scorer == scoreToUpdate.User1Id {
		scoreToUpdate.User1Points += scorerPointsToAdd
		scoreToUpdate.User2Points += opponentPointsToAdd
	} else {
		scoreToUpdate.User1Points += opponentPointsToAdd
		scoreToUpdate.User2Points += scorerPointsToAdd
	}
	return nil
}


func normalizeScoreToJSON (scoreToNormalize *models.Score)  string {
	var normalizeScored = ""

	normalizeScored += fmt.Sprintf("{\n")
	normalizeScored += fmt.Sprintf("\t\"%s\": {\n", scoreToNormalize.User1Id)
	normalizeScored += fmt.Sprintf("\t\t\"points\": \"%d\"\n", scoreToNormalize.User1Points)
	normalizeScored += fmt.Sprintf("\t},\n")
	normalizeScored += fmt.Sprintf("\t\"%s\": {\n", scoreToNormalize.User2Id)
	normalizeScored += fmt.Sprintf("\t\t\"points\": \"%d\"\n", scoreToNormalize.User2Points)
	normalizeScored += fmt.Sprintf("\t}\n")
	normalizeScored += fmt.Sprintf("}\n")

	return normalizeScored
}


func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var databaseConnection *pop.Connection
	var requestError, dbError, updateScoreError error
	var validateError *validate.Errors
	var submittedGoal = &Goal{}
	var goalScore = &models.Score{}
	var databaseConnector = db.DatabaseConnector{}

	databaseConnection, dbError = databaseConnector.GetConnection()
	if databaseConnection != nil {
		defer databaseConnection.Close()
	}
	if dbError != nil {
		return errorResponse(fmt.Sprintf("Failed to connect to database: %s", dbError), http.StatusInternalServerError)
	}
	requestError = json.Unmarshal([]byte(request.Body), submittedGoal)

	if requestError != nil {
		return errorResponse(fmt.Sprintf("Bad request body: %s", requestError), http.StatusBadRequest)
	}

	existingScoreQuery := databaseConnection.Where("user1_id = ? AND user2_id = ? OR user1_id = ? AND user2_id = ?", submittedGoal.Scorer, submittedGoal.Opponent, submittedGoal.Opponent, submittedGoal.Scorer)
	scoreAlreadyExists, dbError := existingScoreQuery.Exists(models.Score{})

	if dbError != nil {
		return errorResponse(fmt.Sprintf("Failed to connect to database: %s", dbError), http.StatusInternalServerError)
	}

	if scoreAlreadyExists {
		dbError = existingScoreQuery.First(goalScore)
		if dbError != nil {
			return errorResponse(fmt.Sprintf("Failed to retrieve existing score: %s", dbError), http.StatusInternalServerError)
		}
	} else {
		goalScore.User1Id = submittedGoal.Scorer
		goalScore.User1Points = 0
		goalScore.User2Id = submittedGoal.Opponent
		goalScore.User2Points = 0
	}

	updateScoreError = updateScore(goalScore, submittedGoal)

	if updateScoreError != nil {
		return errorResponse(fmt.Sprintf("Failed to create/update score: %s", updateScoreError), http.StatusInternalServerError)
	}

	validateError, dbError = databaseConnection.ValidateAndSave(goalScore)

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
