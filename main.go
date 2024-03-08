package main

import (
	"fmt"
	"greatRace/server/db"
	"greatRace/server/router"
	"greatRace/server/utils"
	"net/http"
)

// main is the entry point of the program
func main() {
	//Set configuration for the app
	appVars := utils.SetAppVars()
	db.DBInit(&appVars)
	const port = ":6789"
	router := router.NewRouter(&appVars, port)
	err := http.ListenAndServeTLS(port, appVars.CertPath, appVars.KeyPath, router)
	utils.LogMessage("Server running on port "+port, appVars.LogPath)
	fmt.Println("Server running on port " + port)
	if err != nil {
		utils.LogFatal("There was an error starting the server", err, appVars.LogPath)
	}

}
