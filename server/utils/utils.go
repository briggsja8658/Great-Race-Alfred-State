package utils

import (
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"greatRace/server/types"

	"github.com/fatih/color"
)

func GenerateMultipleIDs(idsNeeded int) []int {
	var ids []int
	for i := 0; i < idsNeeded; i++ {
		// Seed the random number generator with the current time
		rand.Seed(time.Now().UnixNano())

		// Generate a random integer between 0 and 10 Billion
		var randomInt = rand.Intn(10000000001)
		ids = append(ids, randomInt)
	}
	return ids
}

func GenerateID() int {
	// Generate a random integer between 0 and 10 Billion
	randomID := rand.Intn(10000000001)
	return randomID
}

func RandomIntRange(intRange int) int {
	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Generate a random integer between 0 and 10 Billion
	var randomInt = rand.Intn(intRange + 1)
	return randomInt
}

func SetAppVars() types.AppVars {

	var workingDir, dirError = os.Getwd()
	if dirError != nil {
		log.Fatalf("There was an error finding working directory\n%v\n\n", dirError)
	}

	var appVars = types.AppVars{
		DBPath:   workingDir + "/server/db/sqliteDB.sqlite",
		LogPath:  workingDir + "/app.log",
		KeyPath:  workingDir + "/server/keys/server.key",
		CertPath: workingDir + "/server/keys/server.crt",
	}

	return appVars
}

// CurrentFuncName returns a string of the current function
func CurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	return funcName
}

// LogFatal logs the err and stack trace then ends the program
func LogFatal(message string, err error, logPath string) {
	logFile, _ := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log.SetOutput(logFile)
	log.Printf("A FATAL ERROR HAS OCCURED\n%s\n%v\n%s\n\n\n", message, err, debug.Stack())
	logFile.Close()
	redLog := color.New(color.BgHiRed)
	redLog.Printf("A FATAL ERROR HAS OCCURED\n%s\n%v\n%s\n\n\n", message, err, debug.Stack())
	os.Exit(1)
}

// LogErr log the err and stack trace without ending the program
func LogErr(message string, err error, logPath string) {
	logFile, _ := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log.SetOutput(logFile)
	log.Printf("A non-fatal error has occurred\n%s\n%v\n%s\n\n\n", message, err, debug.Stack())
	cyanLog := color.New(color.BgHiCyan)
	cyanLog.Printf("A non-fatal error has occurred\n%s\n%v\n%s\n\n\n", message, err, debug.Stack())
	logFile.Close()
}

func LogMessage(message string, logPath string) {
	logFile, _ := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log.SetOutput(logFile)
	log.Printf("General message\n%s\n\n\n", message)
	greenLog := color.New(color.BgHiGreen)
	greenLog.Printf("General message\n%s\n\n\n", message)
	logFile.Close()
}
