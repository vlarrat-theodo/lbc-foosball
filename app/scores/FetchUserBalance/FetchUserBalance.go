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

type scoreBalance struct {
	Won  int `json:"won"`
	Lost int `json:"lost"`
}

func errorResponse(errorMessage string, errorStatusCode int) (APIResponse events.APIGatewayProxyResponse, APIError error) {
	errorMessage = strings.ReplaceAll(errorMessage, "\"", "\\\"")
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf("{\"error\": \"%s\"}", errorMessage),
		StatusCode: errorStatusCode,
	}, nil
}

func handler(request events.APIGatewayProxyRequest) (APIResponse events.APIGatewayProxyResponse, APIError error) {
	var databaseConnection *pop.Connection
	var databaseConnector = db.DatabaseConnector{}
	var dbError, marshalError error
	var requestedUserID string
	var requestedUserScores []models.Score
	var requestedUserBalance scoreBalance
	var requestedUserBalanceInJSON []byte

	databaseConnection, dbError = databaseConnector.GetConnection()
	if databaseConnection != nil {
		defer databaseConnection.Close()
	}
	if dbError != nil {
		return errorResponse(fmt.Sprintf("Failed to connect to database: %s", dbError), http.StatusInternalServerError)
	}

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

func main() {
	lambda.Start(handler)
}
