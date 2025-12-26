package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"greatRace/server/types"
	"greatRace/server/utils"

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

// FindIfTeamExist from the sqlite3 database adnn return a string slice
func FindIfTeamExist(userData *types.TeamData, appVars *types.AppVars) (bool, error) {
	db, err := sql.Open("sqlite3", appVars.DBPath)
	if err != nil {
		utils.LogErr("There was an error opening the db", err, appVars.LogPath)
		return true, err
	}

	teams, err := db.Query("SELECT * FROM team")
	if err != nil {
		utils.LogErr(
			"There was an error querying the team database id "+utils.CurrentFuncName(),
			err,
			appVars.LogPath,
		)
		return true, err
	}

	var team types.Team
	for teams.Next() {
		// Scan the row values into the struct fields
		err := teams.Scan(&team.TeamID, &team.TeamName)
		if err != nil {
			fmt.Printf("\nthere was an error with the scan in %s", utils.CurrentFuncName())
			return true, err
		}
		if userData.TeamName == team.TeamName {
			return true, nil
		}
	}

	db.Close()
	return false, nil
}

// FindIfUserExist from the sqlite3 database adnn return a string slice
func FindIfUserExist(userData *types.User, appVars *types.AppVars) (bool, error) {
	db, err := sql.Open("sqlite3", appVars.DBPath)
	if err != nil {
		utils.LogErr("There was an error opening the db", err, appVars.LogPath)
		return true, err
	}

	users, err := db.Query("SELECT * FROM users")
	if err != nil {
		utils.LogErr(
			"There was an error querying the team database id "+utils.CurrentFuncName(),
			err,
			appVars.LogPath,
		)
		return true, err
	}

	var user types.User
	for users.Next() {
		// Scan the row values into the struct fields
		err := users.Scan(&user.ID, &user.Name, &user.TeamID)
		if err != nil {
			fmt.Printf("\nthere was an error with the scan in %s", utils.CurrentFuncName())
			return true, err
		}

		if userData.ID == user.ID {
			return true, nil
		}
	}

	db.Close()
	return false, nil
}

// FindUser by id and return a user type
func FindUser(userID *int, appVars *types.AppVars) types.User {
	db, err := sql.Open("sqlite3", appVars.DBPath)
	if err != nil {
		utils.LogErr("There was an error when opening database to find an existing user", err, appVars.LogPath)
	}
	results, err := db.Query(`
		SELECT users.userID, users.name, users.teamID, team.teamName
		FROM users
		INNER JOIN team on team.teamID = users.teamID
		WHERE userID = ?
	`, userID)
	if err != nil {
		utils.LogErr("There was an error when quering the database for a userID", err, appVars.LogPath)
	}

	var user types.User
	for results.Next() {
		err = results.Scan(&user.ID, &user.Name, &user.TeamID, &user.TeamName)
	}
	if err != nil {
		utils.LogErr("There was an error when quering the database for a userID", err, appVars.LogPath)
	}

	return user
}

// CreateNewUser creates a user in the database
func CreateNewUser(userData *types.NewUser, appVars *types.AppVars) error {
	db, err := sql.Open("sqlite3", appVars.DBPath)
	if err != nil {
		utils.LogErr("There was an error in opening the db", err, appVars.LogPath)
		db.Close()
		return err
	}

	sqlStatement, err := db.Prepare("INSERT INTO users(userID, name, teamID) VALUES(?, ?, ?)")
	if err != nil {
		utils.LogErr("There was an error in preparing sql statment", err, appVars.LogPath)
		db.Close()
		return err
	}

	_, err = sqlStatement.Exec(fmt.Sprint(userData.UserID), fmt.Sprint(userData.UserName), fmt.Sprint(userData.TeamID))
	if err != nil {
		utils.LogErr("There was an error in exec the sql statement", err, appVars.LogPath)
		db.Close()
		return err
	}

	teamByID, err := db.Query(
		"SELECT * FROM team WHERE teamID = ?;",
		fmt.Sprint(userData.TeamID),
	)
	if err != nil {
		utils.LogErr("There was an error in exec the sql statement", err, appVars.LogPath)
		db.Close()
		return err
	}

	var currentTeam types.Team

	for teamByID.Next() {
		err = teamByID.Scan(&currentTeam.TeamID, &currentTeam.TeamName, &currentTeam.StartingPoint, &currentTeam.CreatedBy)
		if err != nil {
			utils.LogErr("There was an error when scanning results", err, appVars.LogPath)
		}
	}

	db.Close()
	return nil
}

// CreateNewTeam enters a new team into the database and returns the id of that team
func CreateNewTeam(userData *types.NewUser, appVars *types.AppVars) (int, error) {
	db, err := sql.Open("sqlite3", appVars.DBPath)
	if err != nil {
		utils.LogErr("There was an error in opening the db", err, appVars.LogPath)
		db.Close()
		return 0, err
	}

	teamsFound, err := db.Query(
		"SELECT * FROM team WHERE createdBy = ? LIMIT 1;",
		fmt.Sprint(userData.UserID),
	)
	if err != nil {
		utils.LogErr("There was an error in exec the sql statement", err, appVars.LogPath)
		db.Close()
	}

	numberOfTeams := FindRowsLength(teamsFound)
	if numberOfTeams == 0 {
		sqlStatment, err := db.Prepare("INSERT INTO team(teamID, teamName, startingPoint, createdBy) VALUES(?,?,?,?)")
		if err != nil {
			utils.LogErr("There was an error in preparing sql statment", err, appVars.LogPath)
			db.Close()
			return 0, err
		}

		startingPoint := selectStaringPoint(db, appVars)
		endingPoint := findEndPointID(startingPoint, db, appVars)
		newTeamID := utils.GenerateID()
		_, err = sqlStatment.Exec(fmt.Sprint(strconv.Itoa(newTeamID)), fmt.Sprint(userData.TeamName), fmt.Sprint(startingPoint), fmt.Sprint(userData.UserID))
		if err != nil {
			utils.LogErr("There was an error in entering the new team in the db", err, appVars.LogPath)
			db.Close()
			return 0, err
		}

		sqlStatment, err = db.Prepare("INSERT INTO progress(progressID, finished, timeStarted, timeCompleted, startingPoint, currentPoint, endPoint, teamID) VALUES(?,?,?,?,?,?,?,?)")
		if err != nil {
			utils.LogErr("There was an error in preparing sql statment", err, appVars.LogPath)
			db.Close()
			return 0, err
		}

		progressID := utils.GenerateID()
		_, err = sqlStatment.Exec(
			fmt.Sprint(strconv.Itoa(progressID)), //progressID
			fmt.Sprint(false),                    //finished
			fmt.Sprint(0),                        //timeStarted
			fmt.Sprint(0),                        //timeCompleted
			fmt.Sprint(startingPoint),            //startingPoint
			fmt.Sprint(startingPoint),            //currentPoint (staring point and current point in the begining will be the same)
			fmt.Sprint(endingPoint),              //endPoition
			fmt.Sprint(newTeamID))                //teamID

		if err != nil {
			utils.LogErr("There was an error in entering the new team in the db", err, appVars.LogPath)
			db.Close()
			return 0, err
		}
		return newTeamID, nil
	}

	db.Close()
	return 0, nil
}

func GetAllTeams(appVars *types.AppVars) []types.TeamData {
	db, err := sql.Open("sqlite3", appVars.DBPath)
	if err != nil {
		utils.LogErr("There was an error when opening database", err, appVars.LogPath)
	}
	results, err := db.Query(`SELECT * FROM team`)
	if err != nil {
		utils.LogErr("There was an error when quering the database", err, appVars.LogPath)
	}

	var teamData types.TeamData
	var teamsData []types.TeamData
	for results.Next() {
		err = results.Scan(&teamData.TeamID, &teamData.TeamName, &teamData.StartingPoint, &teamData.CreatedBy)
		teamsData = append(teamsData, teamData)
		if err != nil {
			utils.LogErr("There was an error when scanning results", err, appVars.LogPath)
		}
	}

	return teamsData
}

func CheckTeamName(teamName string, appVars *types.AppVars) (bool, error) {
	found := false
	db, err := sql.Open("sqlite3", appVars.DBPath)
	if err != nil {
		utils.LogErr("There was an error when opening database", err, appVars.LogPath)
	}

	results, err := db.Query(`SELECT * FROM team`)
	if err != nil {
		utils.LogErr("There was an error when quering the database", err, appVars.LogPath)
	}

	var teamEntry types.Team
	for results.Next() {
		err = results.Scan(&teamEntry.TeamID, &teamEntry.TeamName, &teamEntry.StartingPoint, &teamEntry.CreatedBy)
		if err != nil {
			utils.LogErr("There was an error when scanning results", err, appVars.LogPath)
		}

		if strings.EqualFold(teamEntry.TeamName, teamName) {
			found = true
		}
	}

	db.Close()
	return found, err
}

func FindRowsLength(rows *sql.Rows) int {
	count := 0
	for rows.Next() {
		count++
	}
	return count
}

func GetLocations(appVars *types.AppVars) []types.Location {
	db, _ := sql.Open("sqlite3", appVars.DBPath)

	results, err := db.Query(`SELECT * FROM locations`)
	if err != nil {
		utils.LogErr("There was an error when quering the database", err, appVars.LogPath)
	}

	var locations []types.Location
	for results.Next() {
		var location types.Location
		err = results.Scan(&location.LocationID, &location.PositionID, &location.ReadableName, &location.LocationName)
		if err != nil {
			utils.LogErr("There was an error when scanning results", err, appVars.LogPath)
		}

		locations = append(locations, location)
	}

	return locations
}

func GetAllPictures(appVars *types.AppVars) []types.Picture {
	db, _ := sql.Open("sqlite3", appVars.DBPath)

	results, err := db.Query(`SELECT * FROM pictures`)
	if err != nil {
		utils.LogErr("There was an error when quering the database", err, appVars.LogPath)
	}

	var pictures []types.Picture
	for results.Next() {
		var picture types.Picture
		err = results.Scan(&picture.PictureID, &picture.PictureLocation, &picture.PicturePath)
		if err != nil {
			utils.LogErr("There was an error when scanning results", err, appVars.LogPath)
		}

		pictures = append(pictures, picture)
	}

	return pictures
}

// GetProgressLocation finds the progress of the team and returns the id of the target
func GetCurrentTarget(teamID int, appVars *types.AppVars) int {
	db, _ := sql.Open("sqlite3", appVars.DBPath)

	results, err := db.Query(`SELECT * FROM progress WHERE teamID = ? LIMIT 1;`, teamID)
	if err != nil {
		utils.LogErr("There was an error when quering the database", err, appVars.LogPath)
	}

	var progress types.Progess
	for results.Next() {
		err = results.Scan(
			&progress.ProgressID,
			&progress.Finished,
			&progress.TimeStarted,
			&progress.TimeCompleted,
			&progress.StartingPoint,
			&progress.CurrentPoint,
			&progress.EndPoint,
			&progress.TeamID,
		)
		if err != nil {
			utils.LogErr("There was an error when scanning results", err, appVars.LogPath)
		}
	}
	db.Close()
	return progress.CurrentPoint
}

func GetNextTarget(teamID int, appVars *types.AppVars) int {
	db, _ := sql.Open("sqlite3", appVars.DBPath)

	results, err := db.Query(`SELECT * FROM progress WHERE teamID = ? LIMIT 1;`, teamID)
	if err != nil {
		utils.LogErr("There was an error when quering the database", err, appVars.LogPath)
	}

	var progress types.Progess
	for results.Next() {
		err = results.Scan(
			&progress.ProgressID,
			&progress.Finished,
			&progress.TimeStarted,
			&progress.TimeCompleted,
			&progress.StartingPoint,
			&progress.CurrentPoint,
			&progress.EndPoint,
			&progress.TeamID,
		)
		if err != nil {
			utils.LogErr("There was an error when scanning results", err, appVars.LogPath)
		}
	}

	progress.CurrentPoint = progress.CurrentPoint + 1
	db.Close()
	return progress.CurrentPoint
}

func UpdateProgress(teamID int, appVars *types.AppVars) error {
	db, _ := sql.Open("sqlite3", appVars.DBPath)

	results, err := db.Query(`SELECT * FROM progress WHERE teamID = ?;`, teamID)
	if err != nil {
		utils.LogErr("There was an error when quering the database", err, appVars.LogPath)
	}

	var progress types.Progess
	for results.Next() {
		err = results.Scan(
			&progress.ProgressID,
			&progress.Finished,
			&progress.TimeStarted,
			&progress.TimeCompleted,
			&progress.StartingPoint,
			&progress.CurrentPoint,
			&progress.EndPoint,
			&progress.TeamID,
		)
		if err != nil {
			utils.LogErr("There was an error when scanning results", err, appVars.LogPath)
		}
	}

	newLocationID := findNewLocationID(progress.CurrentPoint, db, appVars)

	sqlStatement, err := db.Prepare("UPDATE progress SET CurrentPoint = ? WHERE teamID = ?")
	if err != nil {
		utils.LogErr("There was an error in preparing sql statment", err, appVars.LogPath)
	}

	_, err = sqlStatement.Exec(fmt.Sprint(newLocationID), fmt.Sprint(progress.TeamID))
	if err != nil {
		utils.LogErr("There was an error in exec the sql statement", err, appVars.LogPath)
	}
	db.Close()
	return err
}

func selectStaringPoint(db *sql.DB, appVars *types.AppVars) int {
	results, err := db.Query(`SELECT * FROM locations`)
	if err != nil {
		utils.LogErr("There was an error when quering the database", err, appVars.LogPath)
	}

	var locations []types.Location
	var location types.Location
	for results.Next() {
		err = results.Scan(&location.LocationID, &location.PositionID, &location.ReadableName, &location.LocationName)
		if err != nil {
			utils.LogErr("Error scanning db results", err, appVars.LogPath)
		}
		locations = append(locations, location)
	}

	locationSeed := utils.RandomIntRange(len(locations))
	locationID := locations[locationSeed].LocationID

	return locationID
}

func findNumberOfLocations(db *sql.DB, appVars *types.AppVars) int {
	results, err := db.Query("SELECT * FROM locations")
	if err != nil {
		utils.LogErr("There was an error reading from locations table", err, appVars.LogPath)
	}

	numberOfLocations := FindRowsLength(results)
	return numberOfLocations
}

// CreateDB check to see if database exist and if not create it type of *Sqlite is the expected paramiter and error is the return type
func createDB(appVars *types.AppVars) {
	db, err := sql.Open("sqlite3", appVars.DBPath)
	if err != nil {
		utils.LogFatal("There was an error opening the database when populating it with data", err, appVars.LogPath)
	}

	transaction, _ := startTransaction(db)
	createTables(transaction, appVars)
	transaction.Commit()

	transaction, _ = startTransaction(db)
	insertDefaultData(transaction, db, appVars)
	transaction.Commit()

	db.Close()
}

// createTables if the database is empty.
func createTables(transaction *sql.Tx, appVars *types.AppVars) {
	_, err := transaction.Exec(`CREATE TABLE IF NOT EXISTS team (
			teamID INTEGER PRIMARY KEY,
			teamName CHAR(50),
			startingPoint INTEGER,
			createdBy INTEGER,
			FOREIGN KEY (createdBy) REFERENCES users (userID),
			FOREIGN KEY (startingPoint) REFERENCES locations (locationID)
		);`)
	if err != nil {
		utils.LogFatal("There was an error creating team table", err, appVars.LogPath)
	}

	_, err = transaction.Exec(`CREATE TABLE IF NOT EXISTS users (
			userID INTEGER PRIMARY KEY,
			name CHAR(50),
			teamID INTEGER,
			FOREIGN KEY (teamID) REFERENCES team (teamID)
		);`)
	if err != nil {
		utils.LogFatal("There was an error creating users table", err, appVars.LogPath)
	}

	_, err = transaction.Exec(`CREATE TABLE IF NOT EXISTS progress (
			progressID INTEGER PRIMARY KEY,
			finished BOOLEAN,
			timeStarted INTEGER,
			timeCompleted INTEGER,
			startingPoint INTEGER,
			currentPoint INTEGER,
			endPoint INTEGER,
			teamID INTEGER,
			FOREIGN KEY (startingPoint) REFERENCES locations (locationID),
			FOREIGN KEY (currentPoint) REFERENCES locations (locationID),
			FOREIGN KEY (endPoint) REFERENCES locations (locationID),
			FOREIGN KEY (teamID) REFERENCES team (teamID)
		);`)
	if err != nil {
		utils.LogFatal("There was an error creating progress table", err, appVars.LogPath)
	}

	_, err = transaction.Exec(`CREATE TABLE IF NOT EXISTS locations (
			locationID INTEGER PRIMARY KEY,
			positionID INTEGER,
			readableName CHAR(50),
			locationName CHAR(50)
		);`)
	if err != nil {
		utils.LogFatal("There was an error creating locaitons table", err, appVars.LogPath)
	}

	_, err = transaction.Exec(`CREATE TABLE IF NOT EXISTS pictures (
			pictureID INTEGER PRIMARY KEY,
			pictureLocation CHAR(100),
			picturePath CHAR(100)
		);`)
	if err != nil {
		utils.LogFatal("There was an error creating picture table", err, appVars.LogPath)
	}
}

func insertDefaultData(transaction *sql.Tx, db *sql.DB, appVars *types.AppVars) {
	var allLocations []types.Location
	locationRows, err := transaction.Query("SELECT * FROM locations")
	for locationRows.Next() {
		var location types.Location
		// Scan the row values into the struct fields
		err := locationRows.Scan(&location.LocationID, &location.PositionID, &location.ReadableName, &location.LocationName)
		if err != nil {
			utils.LogFatal("There was an error creating the database", err, appVars.LogPath)
		}
		allLocations = append(allLocations, location)
	}

	sqlStatement, err := transaction.Prepare("INSERT INTO locations(locationID, positionID, readableName, locationName) VALUES(?, ?, ?, ?) ")
	if err != nil {
		utils.LogFatal("There was an error preping locations data", err, appVars.LogPath)
	}

	positionID := 1
	for x := 0; x < len(appVars.LocationNames); x++ {
		found := false
		for y := 0; y < len(allLocations); y++ {
			if allLocations[y].LocationName == appVars.LocationNames[x] {
				found = true
			}
		}
		if !found {
			locationID := utils.GenerateID()
			_, err := sqlStatement.Exec(fmt.Sprint(locationID), fmt.Sprint(positionID), fmt.Sprint(appVars.ReadableNames[x]), fmt.Sprint(appVars.LocationNames[x]))
			if err != nil {
				utils.LogFatal("There was an error inserting locations data", err, appVars.LogPath)
			}
		}
		positionID++
	}

	sortedLocations, err := utils.SortDir(appVars.PicturesPath)
	if err != nil {
		utils.LogFatal("There was an error sorting the pictures dir", err, appVars.LogPath)
	}

	dbPictures, err := findStoredPictures(db)
	dbPictureNames := getPictureNames(dbPictures)
	if err != nil {
		utils.LogFatal("There was an error stored pictures", err, appVars.LogPath)
	}

	dirPictures, err := getDirPictures(sortedLocations)
	if err != nil {
		utils.LogFatal("There was an error finding pictureNames", err, appVars.LogPath)
	}

	dirPicturePaths := createPicturePath(dirPictures)
	newPictures := utils.FindNewString(dirPicturePaths, dbPictureNames)
	newLocations := findNewLocations(newPictures, dirPictures)

	transaction.Exec(`DELETE FROM pictures`) //Delete the prior pictures to avoid duplications
	sqlPictures, err := transaction.Prepare("INSERT INTO pictures(pictureID, picturePath, pictureLocation) VALUES(?,?,?)")
	if err != nil {
		utils.LogFatal("There was an error preping picture data", err, appVars.LogPath)
	}

	for x := 0; x < len(newPictures); x++ {
		pictureIDs := utils.GenerateMultipleIDs(len(newPictures))
		_, err := sqlPictures.Exec(fmt.Sprint(pictureIDs[x]), fmt.Sprintf(newPictures[x]), fmt.Sprint(newLocations[x]))
		if err != nil {
			utils.LogFatal("There was an error inserting picture data", err, appVars.LogPath)
		}
	}
}

// startTransaction create a new transaction in sqlite3. This function assumes that there will be an open connection.
func startTransaction(db *sql.DB) (*sql.Tx, error) {
	transaction, err := db.Begin()
	if err != nil {
		fmt.Printf("\nThere was an creating a new transaction in %s", utils.CurrentFuncName())
		return nil, err
	} else {
		return transaction, nil
	}
}

// findStoredPictures from the sqlite3 database and return a string slice
func findStoredPictures(db *sql.DB) ([]types.Picture, error) {

	var pictures []types.Picture
	storedPictures, err := db.Query("SELECT * FROM pictures")
	if err != nil {
		fmt.Printf("\nThere was an error with the db.Query in %s", utils.CurrentFuncName())
		return nil, err
	}

	count := 0
	storedPictures.Scan(&count)
	if count != 0 {
		for storedPictures.Next() {
			var picture types.Picture

			// Scan the row values into the struct fields
			err := storedPictures.Scan(&picture.PictureID, &picture.PicturePath, &picture.PictureLocation, &picture.LocationID)
			if err != nil {
				fmt.Printf("\nthere was an error with the scan in %s", utils.CurrentFuncName())
				return nil, err
			}
			pictures = append(pictures, picture)
		}
		//Checks for errors during the loop of rows
		err = storedPictures.Err()
		if err != nil {
			fmt.Printf("There was an error in reading rows in %s", utils.CurrentFuncName())
			return nil, err
		}
	}

	return pictures, nil
}

func getPictureNames(pictures []types.Picture) []string {
	var pictureNames []string
	for x := 0; x < len(pictures); x++ {
		pictureNames = append(pictureNames, pictures[x].PicturePath)
	}
	return pictureNames
}

func getDirPictures(picDir []string) ([]types.DirPicture, error) {
	var dirPictures []types.DirPicture
	for x := 0; x < len(picDir); x++ { //Start the loop with the directory that was given
		_, err := os.Stat(picDir[x])
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Nothing found in %s\n", picDir[x])
			} else {
				return nil, err
			}
		}

		files, _ := os.ReadDir(picDir[x])
		for y := 0; y < len(files); y++ {
			if files[y].IsDir() {
				nestedFiles, err := os.ReadDir(picDir[x] + files[y].Name()) //If file[y] is a directory step into it and get the pictures
				if err != nil {
					return nil, err
				}
				for z := 0; z < len(nestedFiles); z++ {
					if !nestedFiles[x].IsDir() {
						dirPicture := types.DirPicture{
							Location:    filepath.Base(picDir[x]),
							PictureName: nestedFiles[z].Name(),
						}
						dirPictures = append(dirPictures, dirPicture)
					}
				}
			} else {
				dirPicture := types.DirPicture{
					Location:    filepath.Base(picDir[x]),
					PictureName: files[y].Name(),
				}
				dirPictures = append(dirPictures, dirPicture)
			}
		}
	}
	return dirPictures, nil
}

func createPicturePath(dirPictures []types.DirPicture) []string {
	var picturePath []string
	for x := 0; x < len(dirPictures); x++ {
		picturePath = append(picturePath, ("/img/" + dirPictures[x].PictureName))
	}
	return picturePath
}

func findNewLocations(newPictures []string, dirPictures []types.DirPicture) []string {
	var newLocations []string
	for x := 0; x < len(dirPictures); x++ {
		for y := 0; y < len(newPictures); y++ {
			if ("/img/" + dirPictures[x].PictureName) == newPictures[y] {
				newLocations = append(newLocations, dirPictures[x].Location)
			}
		}
	}
	return newLocations
}

func findNewLocationID(locationID int, db *sql.DB, appVars *types.AppVars) int {

	locations, err := db.Query("SELECT * FROM locations")
	if err != nil {
		utils.LogErr("Error reading from DB", err, appVars.LogPath)
	}

	var allLocations []types.Location
	for locations.Next() {
		var location types.Location
		locations.Scan(&location.LocationID, &location.PositionID, &location.ReadableName, &location.LocationName)
		allLocations = append(allLocations, location)
	}

	var newPositionID int
	var newLocationID int
	for x := 0; x < len(allLocations); x++ {
		if locationID == allLocations[x].LocationID {
			positionID := allLocations[x].PositionID
			if positionID == 12 {
				newPositionID = 1
			} else {
				newPositionID = positionID + 1
			}
		}
	}

	for x := 0; x < len(allLocations); x++ {
		if newPositionID == allLocations[x].PositionID {
			newLocationID = allLocations[x].LocationID
		}
	}
	return newLocationID
}

func findEndPointID(locationID int, db *sql.DB, appVars *types.AppVars) int {

	locations, err := db.Query("SELECT * FROM locations")
	if err != nil {
		utils.LogErr("Error reading from DB", err, appVars.LogPath)
	}

	var allLocations []types.Location
	for locations.Next() {
		var location types.Location
		locations.Scan(&location.LocationID, &location.PositionID, &location.ReadableName, &location.LocationName)
		allLocations = append(allLocations, location)
	}

	var newPositionID int
	var endPointID int
	for x := 0; x < len(allLocations); x++ {
		if locationID == allLocations[x].LocationID {
			positionID := allLocations[x].PositionID
			if positionID == 1 {
				newPositionID = 12
			} else {
				newPositionID = positionID - 1
			}
		}
	}

	for x := 0; x < len(allLocations); x++ {
		if newPositionID == allLocations[x].PositionID {
			endPointID = allLocations[x].LocationID
		}
	}
	return endPointID
}
