package db

import (
	"database/sql"
	"fmt"
	"greatRace/server/types"
	"greatRace/server/utils"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3" //Reference the package without directly using. This is typically done for drivers
)

// DBInit checks if the database file exist and if it's not found creates one. Then fills it with the default data
func DBInit(appVars *types.AppVars) {
	fileInfo, err := os.Stat(appVars.DBPath)
	if fileInfo == nil { //If no database was found nil will be the result
		os.Create(appVars.DBPath) //create the database file
	} else if err != nil {
		utils.LogFatal("There was a error in checking if the Sqlite3 database exist", err, appVars.LogPath)
	}
	createDB(appVars)
}

// CreateDB check to see if database exist and if not create it type of *Sqlite is the expected paramiter and error is the return type
func createDB(appVars *types.AppVars) {
	db, err := sql.Open("sqlite3", appVars.DBPath)
	if err != nil {
		utils.LogFatal("There was an error opening the database when populating it with data", err, appVars.LogPath)
	} else {
		createTables(db, appVars)
	}
	db.Close()
}

// createTables if the database is empty.
func createTables(db *sql.DB, appVars *types.AppVars) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
			userID INTEGER PRIMARY KEY,
			firstName STRING,
			lastName STRING,
			finished BOOLEAN,
			timeStarted INTEGER,
			timeCompleted INTEGER,
			pioneerStadium INTEGER,
			taParishHall INTEGER,
			studentDevelopmentCenter INTEGER,
			hindelLibrary INTEGER,
			financialAid INTEGER,
			studentLeadershipCenter INTEGER,
			ejBrownHall INTEGER,
			mailCenter INTEGER,
			orvisActivitiesCenter INTEGER,
			baseballField INTEGER,
			softballField INTEGER
		);`)
	if err != nil {
		utils.LogFatal("There was an error creating users table", err, appVars.LogPath)
	}
}

// CreateNewUser creates a user in the database
func CreateNewUser(userData *types.NewUser, appVars *types.AppVars) error {
	db, err := sql.Open("sqlite3", appVars.DBPath)
	if err != nil {
		utils.LogErr("There was an error in opening the db", err, appVars.LogPath)
		db.Close()
		return err
	}

	sqlStatement, err := db.Prepare("INSERT INTO users(userID, firstName, lastName, finished, timeStarted, timeCompleted, pioneerStadium, taParishHall, studentDevelopmentCenter, hindelLibrary, financialAid, studentLeadershipCenter, ejBrownHall, mailCenter, orvisActivitiesCenter, baseballField, softballField) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		utils.LogErr("There was an error in preparing sql statment", err, appVars.LogPath)
		db.Close()
		return err
	}

	_, err = sqlStatement.Exec(
		fmt.Sprint(userData.UserID),        //userID
		fmt.Sprint(userData.UserName),      //firstName
		fmt.Sprint(userData.TeamID),        //lastName
		fmt.Sprint(false),                  //finished
		fmt.Sprint(time.Now().UnixMilli()), //timeStarted
		fmt.Sprint(0),                      //timeCompleted
		fmt.Sprint(0),                      //pioneerStadium
		fmt.Sprint(0),                      //taParishHall
		fmt.Sprint(0),                      //studentDevelopmentCenter
		fmt.Sprint(0),                      //hindelLibrary
		fmt.Sprint(0),                      //financialAid
		fmt.Sprint(0),                      //studentLeadershipCenter
		fmt.Sprint(0),                      //ejBrownHall
		fmt.Sprint(0),                      //mailCenter
		fmt.Sprint(0),                      //orvisActivitiesCenter
		fmt.Sprint(0),                      //baseballField
		fmt.Sprint(0),                      //softballField
	)
	if err != nil {
		utils.LogErr("There was an error creating user", err, appVars.LogPath)
		return err
	}

	db.Close()
	return nil
}
