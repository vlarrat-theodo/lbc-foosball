package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gobuffalo/pop"
	"github.com/vlarrat-theodo/lbc-foosball/db"
	"github.com/vlarrat-theodo/lbc-foosball/models"
	"net/http"
	"strings"
)

// scoreBalance represents sum of sets won and lost by one user.
//
type scoreBalance struct {
	Won  int `json:"won"`
	Lost int `json:"lost"`
}

// errorResponse formats API HTTP responses sent when an error occurs.
//
func errorResponse(errorMessage string, errorStatusCode int) (APIResponse events.APIGatewayProxyResponse, APIError error) {
	errorMessage = strings.ReplaceAll(errorMessage, "\"", "\\\"")
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf("{\"error\": \"%s\"}", errorMessage),
		StatusCode: errorStatusCode,
	}, nil
}

// handler is the main function launched by Lambda.
//
// In this Lambda, it will:
//     - retrieve user_id from API request
//     - retrieve from DB all scores regarding requested user
//     - calculate sum of won and lost sets by requested user
//     - send HTTP JSON response containing this information
//
func handler(request events.APIGatewayProxyRequest) (APIResponse events.APIGatewayProxyResponse, APIError error) {
	var databaseConnection *pop.Connection
	var databaseConnector = db.DatabaseConnector{}
	var dbError, marshalError error
	var requestedUserID string
	var requestedUserScores []models.Score
	var requestedUserBalance scoreBalance
	var requestedUserBalanceInJSON []byte

	databaseConnection, dbError = databaseConnector.GetConnection()
	if dbError != nil {
		return errorResponse(fmt.Sprintf("Failed to connect to database: %s", dbError), http.StatusInternalServerError)
	}
	defer databaseConnection.Close()

	requestedUserID = request.QueryStringParameters["user_id"]
	if requestedUserID == "" {
		return errorResponse("Bad request: you must provide a value for 'user_id' parameter", http.StatusBadRequest)
	}

	dbError = databaseConnection.Where("user1_id = ? or user2_id = ?", requestedUserID, requestedUserID).All(&requestedUserScores)
	if dbError != nil {
		return errorResponse(fmt.Sprintf("Failed to retrieve user's scores for user_id '%s'", requestedUserID), http.StatusInternalServerError)
	}

	for _, requestedUserScore := range requestedUserScores {
		switch requestedUserID {
		case requestedUserScore.User1Id:
			requestedUserBalance.Won += requestedUserScore.User1Sets
			requestedUserBalance.Lost += requestedUserScore.User2Sets
		case requestedUserScore.User2Id:
			requestedUserBalance.Won += requestedUserScore.User2Sets
			requestedUserBalance.Lost += requestedUserScore.User1Sets
		}
	}

	requestedUserBalanceInJSON, marshalError = json.Marshal(requestedUserBalance)
	if marshalError != nil {
		return errorResponse(fmt.Sprintf("Failed to JSONify user balance: %s", marshalError), http.StatusInternalServerError)
	}

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(requestedUserBalanceInJSON),
		StatusCode: http.StatusOK,
	}, nil
}

// Main launches Lambda function.
//
func main() {
	lambda.Start(handler)
}
