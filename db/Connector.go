package db

import (
	"github.com/gobuffalo/pop"
	"os"
)

// DatabaseConnector represents handler to connect to database specified in environment variables.
//
type DatabaseConnector struct {
}

// GetConnection establishes connection to database according to parameters stored in environment variables.
//
func (p DatabaseConnector) GetConnection() (db *pop.Connection, err error) {
	var dbConnectionsDetails pop.ConnectionDetails

	dbConnectionsDetails.Dialect = os.Getenv("DB_DIALECT")
	dbConnectionsDetails.Host = os.Getenv("DB_HOST")
	dbConnectionsDetails.Port = os.Getenv("DB_PORT")
	dbConnectionsDetails.Database = os.Getenv("DB_NAME")
	dbConnectionsDetails.User = os.Getenv("DB_USERNAME")
	dbConnectionsDetails.Password = os.Getenv("DB_PASSWORD")
	dbConnectionsDetails.RawOptions = "sslmode=" + os.Getenv("DB_SSLMODE")
	dbConnection, dbError := pop.NewConnection(&dbConnectionsDetails)

	if dbError != nil {
		return dbConnection, dbError
	}

	dbError = dbConnection.Open()

	return dbConnection, dbError
}
