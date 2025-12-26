package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"greatRace/server/db"
	"greatRace/server/types"
	"greatRace/server/utils"

	"golang.org/x/net/websocket"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

var newTeamSSEChan chan string
var routerAppVars *types.AppVars
var routeCache map[string]*template.Template
var jsCache map[string][]byte
var cssCache map[string][]byte
var imgCache map[string][]byte
var wsServer types.WSServer

func NewRouter(appVars *types.AppVars, port string) *chi.Mux {
	mux := chi.NewRouter()

	wsServer = newWSServer()

	//Cors config
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"localhost" + port, "127.0.0.1" + port, "it.alfredstate.edu" + port},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum cache duration of preflight requests
	})

	mux.Use(cors.Handler)

	//Set up global varialbes for the router
	routerAppVars = appVars
	routeCache = createTemplateCache()
	jsCache = createJSCache()
	cssCache = createCSSCache()
	imgCache = createImgCache()

	//CONTENT
	//js
	mux.Get("/js/*", getJSFile)
	//css
	mux.Get("/css/*", getCSSFile)
	//img
	mux.Get("/img/*", getImgFile)

	//PAGES
	//index
	mux.Get("/", index)

	//END POINTS
	//server side events
	mux.Get("/sse/newTeam", newTeamSSE)

	//gets
	mux.Get("/data/locations", getLocations)
	mux.Get("/data/pictures", getPictures)
	mux.Get("/data/teams", getCurrentTeams)
	mux.Get("/data/currentTarget", getCurrentTarget)
	//mux.Get("/data/nextTarget", nextTarget)
	mux.Get("/data/updateProgress", updateProgress)

	//posts
	mux.Post("/data/newUser", newUser)
	mux.Post("/data/checkTeamName", checkTeamName)

	mux.NotFound(pageNotfound)
	return mux
}

// PAGE HANDLERS
func index(writer http.ResponseWriter, request *http.Request) {
	userID, userFound := getUserIDFromCookie(request)
	stringMapData := map[string]string{}
	if userFound {
		userData := db.FindUser(&userID, routerAppVars)
		stringMapData = utils.CreateStringMapUserData(&userData)
	} else {
		stringMapData = utils.CreateStringMapDefaultUserData()
	}
	renderTemplate(
		writer,
		"index.page.tmpl",
		&types.TemplateData{
			StringMap: stringMapData,
		},
	)
}

func pageNotfound(writer http.ResponseWriter, response *http.Request) {

}

// WEB SOCKET
func newWSServer() types.WSServer {
	return types.WSServer{
		WsConnections: make(map[*websocket.Conn]bool),
	}
}

// DATA HANDLERS
func getCurrentTeams(writer http.ResponseWriter, request *http.Request) {
	teams := db.GetAllTeams(routerAppVars)
	marshaledJSON, err := json.Marshal(teams)
	if err != nil {
		utils.LogErr("There was an error creating JSON object", err, routerAppVars.LogPath)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(marshaledJSON)
	}
}

func newUser(writer http.ResponseWriter, request *http.Request) {
	//Decode the JSON data from the request body into a User struct
	var newUserData types.NewUser
	err := json.NewDecoder(request.Body).Decode(&newUserData)
	if err != nil {
		utils.LogErr("Client sent invalid JSON data", err, routerAppVars.LogPath)
		http.Error(writer, "Invalid JSON data", http.StatusBadRequest)
	} else {
		if newUserData.TeamID != 0 {
			err = db.CreateNewUser(&newUserData, routerAppVars)
			if err != nil {
				utils.LogErr("There was an error entering a new user into the database", err, routerAppVars.LogPath)
				http.Error(writer, "Error creating youruserNameValue, newTeamNameValue profile, please try again later", http.StatusInternalServerError)
			} else {
				marshaledJSON, err := json.Marshal(&newUserData)
				if err != nil {
					utils.LogErr("There was an error creating the json to send to the client", err, routerAppVars.LogPath)
					http.Error(writer, "Error creating your profile, please try again later", http.StatusInternalServerError)
				} else {
					setCookie(&newUserData, writer, request)
					writer.Header().Set("Content-Type", "application/json")
					writer.Write(marshaledJSON)
				}
			}
		} else {
			teamID, err := db.CreateNewTeam(&newUserData, routerAppVars)
			if err != nil {
				utils.LogErr("There was an error entering a new team into the database", err, routerAppVars.LogPath)
				http.Error(writer, "Error creating your profile, please try again later", http.StatusInternalServerError)
			} else {
				newUserData.TeamID = teamID
				err = db.CreateNewUser(&newUserData, routerAppVars)
				if err != nil {
					utils.LogErr("There was an error entering a new user into the database", err, routerAppVars.LogPath)
					http.Error(writer, "Error creating your profile, please try again later", http.StatusInternalServerError)
				} else {
					marshaledJSON, err := json.Marshal(&newUserData)
					if err != nil {
						utils.LogErr("There was an error creating the json to send to the client", err, routerAppVars.LogPath)
						http.Error(writer, "Error creating your profile, please try again later", http.StatusInternalServerError)
					} else {
						setCookie(&newUserData, writer, request)
						writer.Header().Set("Content-Type", "application/json")
						writer.Write(marshaledJSON)
					}
				}
			}
		}
	}
	cookies := request.Cookies()
	for _, cookie := range cookies {
		fmt.Println("Name: ", cookie.Name)
	}
}

func checkTeamName(writer http.ResponseWriter, request *http.Request) {
	var checkNames types.CheckNames
	err := json.NewDecoder(request.Body).Decode(&checkNames)
	if err != nil {
		utils.LogErr("Client sent invalid JSON data", err, routerAppVars.LogPath)
		http.Error(writer, "Invalid JSON data", http.StatusBadRequest)
	} else {
		found, err := db.CheckTeamName(checkNames.TeamName, routerAppVars)
		if err != nil {
			utils.LogErr("There was an error checking userName in the database", err, routerAppVars.LogPath)
			http.Error(writer, "Error creating your profile, please try again later", http.StatusInternalServerError)
		} else {
			marshaledJSON, err := json.Marshal(found)
			if err != nil {
				utils.LogErr("There was an error creating json to send to the client", err, routerAppVars.LogPath)
				http.Error(writer, "Error creating your profile, please try again later", http.StatusInternalServerError)
			} else {
				writer.Header().Set("Content-Type", "application/json")
				writer.Write(marshaledJSON)
			}
		}

	}
}

func getLocations(writer http.ResponseWriter, request *http.Request) {
	locations := db.GetLocations(routerAppVars)
	marshaledJSON, err := json.Marshal(locations)
	if err != nil {
		utils.LogErr("There was an error creating json to send to the client", err, routerAppVars.LogPath)
		http.Error(writer, "Error creating your profile, please try again later", http.StatusInternalServerError)
	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(marshaledJSON)
	}
}

func getCurrentTarget(writer http.ResponseWriter, request *http.Request) {
	teamID, _ := getTeamIDFromCookie(request)
	currentTarget := db.GetCurrentTarget(teamID, routerAppVars)
	marshaledJSON, err := json.Marshal(currentTarget)
	if err != nil {
		utils.LogErr("There was an error creating json to send to the client", err, routerAppVars.LogPath)
		http.Error(writer, "Error creating your profile, please try again later", http.StatusInternalServerError)
	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(marshaledJSON)
	}
}

func updateProgress(writer http.ResponseWriter, request *http.Request) {
	teamID, _ := getTeamIDFromCookie(request)
	err := db.UpdateProgress(teamID, routerAppVars)
	if err != nil {
		utils.LogErr("There was an error creating json to send to the client", err, routerAppVars.LogPath)
		http.Error(writer, "Error creating your profile, please try again later", http.StatusInternalServerError)
	} else {
		writer.WriteHeader(http.StatusOK)
	}
}

func getPictures(writer http.ResponseWriter, request *http.Request) {
	allPictures := db.GetAllPictures(routerAppVars)
	marshaledJSON, err := json.Marshal(allPictures)
	if err != nil {
		utils.LogErr("There was an error creating json to send to the client", err, routerAppVars.LogPath)
		http.Error(writer, "Error creating your profile, please try again later", http.StatusInternalServerError)
	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(marshaledJSON)
	}
}

// FILE HANDLERS
func getJSFile(writer http.ResponseWriter, request *http.Request) {
	reqPath := request.URL.Path
	jsFile := jsCache[reqPath]
	fileSize := len(jsFile)
	if fileSize == 0 {
		utils.LogErr("The file of "+reqPath+" was not found", nil, routerAppVars.LogPath)
		http.NotFound(writer, request) //Send a 404 to the client
	} else if fileSize > 0 {
		writer.Header().Set("Content-Type", "application/javascript")
		_, err := writer.Write(jsFile)
		if err != nil {
			utils.LogErr("There was an error sending a javascript file to the client", err, routerAppVars.LogPath)
		}
	} else {
		utils.LogErr("There was a general error trying to find a JS file", nil, routerAppVars.LogPath)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getCSSFile(writer http.ResponseWriter, request *http.Request) {
	reqPath := request.URL.Path
	cssFile := cssCache[reqPath]
	fileSize := len(cssFile)
	if fileSize == 0 {
		utils.LogErr("The file of "+reqPath+" was not found", nil, routerAppVars.LogPath)
		http.NotFound(writer, request) //Send a 404 to the client
	} else if fileSize > 0 {
		writer.Header().Set("Content-Type", "text/css")
		status, err := writer.Write(cssFile)
		if status == 0 {
			utils.LogErr("There was an error sending a css file to the client", err, routerAppVars.LogPath)
		}
	} else {
		//Something else happened
		utils.LogErr("There was a general error trying to find a css file", nil, routerAppVars.LogPath)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getImgFile(writer http.ResponseWriter, request *http.Request) {
	reqPath := request.URL.Path
	imgFile := imgCache[reqPath]
	fileSize := len(imgFile)
	if fileSize == 0 {
		utils.LogErr("The file of "+reqPath+" was not found", nil, routerAppVars.LogPath)
		http.NotFound(writer, request) //Send a 404 to the client
	} else if fileSize > 0 {
		writer.Header().Set("Content-Type", "image/webp")
		_, err := writer.Write(imgFile)
		if err != nil {
			utils.LogErr("There was an error sending a img file to the client", err, routerAppVars.LogPath)
		}
	} else {
		//Something else happened
		utils.LogErr("There was a general error trying to find a img file", nil, routerAppVars.LogPath)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

// SSE HANDLERS
func newTeamSSE(writer http.ResponseWriter, request *http.Request) {

}
func newTeamSSETrigger() {

}

// HELPER FUNCTIONS
// RenderTemplate renders a cached template with the given values
func renderTemplate(
	writer http.ResponseWriter,
	pageName string,
	templateData *types.TemplateData,
) {
	cachePage := routeCache[pageName]
	buf := new(bytes.Buffer)
	err := cachePage.Execute(buf, templateData)
	if err != nil {
		utils.LogErr("Error when executing chacePage", err, routerAppVars.LogPath)
	}
	_, err = buf.WriteTo(writer)
	if err != nil {
		utils.LogErr("error writing template to browser", err, routerAppVars.LogPath)
	}
}

// getUserIDFromCookie if there is one.
// If there is return the id with a bool of true stating that the user was found
// If there is no user then return 0 with a bool of false
func getUserIDFromCookie(request *http.Request) (int, bool) {
	currentCookie, _ := request.Cookie("greatRaceCookie")

	if currentCookie != nil {
		cookieValues := strings.Split(currentCookie.Value, ",")
		userIDInt, err := strconv.Atoi(cookieValues[0])
		if err != nil {
			utils.LogErr("There was an error getting the cookie from the request", err, routerAppVars.LogPath)
		}
		return userIDInt, true
	} else {
		return 0, false //There is no cookie
	}
}

func getTeamIDFromCookie(request *http.Request) (int, bool) {
	currentCookie, _ := request.Cookie("greatRaceCookie")

	if currentCookie != nil {
		cookieValues := strings.Split(currentCookie.Value, ",")
		userIDInt, err := strconv.Atoi(cookieValues[2])
		if err != nil {
			utils.LogErr("There was an error getting the cookie from the request", err, routerAppVars.LogPath)
		}
		return userIDInt, true
	} else {
		return 0, false //There is no cookie
	}
}

// createTemplateCache creates a template cache as a map
func createTemplateCache() map[string]*template.Template {
	schema := template.FuncMap{}
	templateCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./server/templates/*.page.tmpl")
	if err != nil {
		utils.LogFatal("There was an error creating template cache", err, routerAppVars.LogPath)
	}

	for x := 0; x < len(pages); x++ {
		pageName := filepath.Base(pages[x])
		pageString, err := template.New(pageName).Funcs(schema).ParseFiles(pages[x])
		if err != nil {
			utils.LogFatal("There was an error creating template cache", err, routerAppVars.LogPath)
		}

		layoutString, err := filepath.Glob("./server/templates/*.layout.tmpl")
		if err != nil {
			utils.LogFatal("There was an error creating template cache", err, routerAppVars.LogPath)
		}

		if len(layoutString) > 0 {
			pageString, err = pageString.ParseGlob("./server/templates/*.layout.tmpl")
			if err != nil {
				utils.LogFatal("There was an error creating template cache", err, routerAppVars.LogPath)
			}
		}

		templateCache[pageName] = pageString
	}

	return templateCache
}

// createJSCache creates a map of JS byte data
func createJSCache() map[string][]byte {
	jsCache := map[string][]byte{}

	_ = filepath.WalkDir("./client/js/", func(path string, fileInfo os.DirEntry, err error) error {
		if err != nil {
			utils.LogFatal("There was an error in creating the jsCache", err, routerAppVars.LogPath)
		}
		if !fileInfo.IsDir() {
			basePath := filepath.Base(path)
			fileContent, err := os.ReadFile(path)
			if err != nil {
				utils.LogFatal("There was an error in creating the jsCache", err, routerAppVars.LogPath)
			}
			jsCache["/js/"+basePath] = fileContent
		}
		return err
	})

	return jsCache
}

func createCSSCache() map[string][]byte {
	cssCache := map[string][]byte{}

	_ = filepath.WalkDir("./client/css/", func(path string, fileInfo os.DirEntry, err error) error {
		if err != nil {
			utils.LogFatal("There was an error in creating the cssCache", err, routerAppVars.LogPath)
		}
		if !fileInfo.IsDir() {
			basePath := filepath.Base(path)
			fileContent, err := os.ReadFile(path)
			if err != nil {
				utils.LogFatal("There was an error in creating the cssCache", err, routerAppVars.LogPath)
			}
			cssCache["/css/"+basePath] = fileContent
		}
		return err
	})

	return cssCache
}

func createImgCache() map[string][]byte {
	imgCache := map[string][]byte{}

	_ = filepath.WalkDir("./server/pictures/", func(path string, fileInfo os.DirEntry, err error) error {
		if err != nil {
			utils.LogFatal("There was an error in creating the imgCache", err, routerAppVars.LogPath)
		}
		if !fileInfo.IsDir() && filepath.Ext(fileInfo.Name()) == ".webp" {
			basePath := filepath.Base(path)
			fileContent, err := os.ReadFile(path)
			if err != nil {
				utils.LogFatal("There was an error in creating the imgCache", err, routerAppVars.LogPath)
			}
			imgCache["/img/"+basePath] = fileContent

		}
		return err
	})

	return imgCache
}

func setSEEHeader(writer http.ResponseWriter, request *http.Request) http.ResponseWriter {
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	return writer
}

func setCookie(newUserData *types.NewUser, writer http.ResponseWriter, request *http.Request) {
	values := []string{strconv.Itoa(newUserData.UserID), newUserData.TeamName, strconv.Itoa(newUserData.TeamID)}
	valuesString := strings.Join(values, ",")
	http.SetCookie(writer, &http.Cookie{
		Name:     "greatRaceCookie",
		Domain:   "127.0.0.1:6789",
		Path:     "/",
		Secure:   true,
		Value:    valuesString,
		Expires:  time.Now().Add(720 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
