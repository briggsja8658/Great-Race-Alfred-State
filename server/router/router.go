package router

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"greatRace/server/db"
	"greatRace/server/types"
	"greatRace/server/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

var routerAppVars *types.AppVars
var routeCache map[string]*template.Template
var imgCache map[string][]byte

func NewRouter(appVars *types.AppVars, port string) *chi.Mux {
	mux := chi.NewRouter()

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
	imgCache = createImgCache()

	//CONTENT
	mux.Get("/img/*", getImgFile)

	//PAGES
	//gets
	mux.Get("/", index)
	mux.Get("/pioneerStadium", pioneerStadium)
	mux.Get("/taParishHall", taParishHall)
	mux.Get("/studentDevelopmentCenter", studentDevelopmentCenter)
	mux.Get("/hindelLibrary", hindelLibrary)
	mux.Get("/financialAid", financialAid)
	mux.Get("/studentLeadershipCenter", studentLeadershipCenter)
	mux.Get("/pioneerCenter", pioneerCenter)
	mux.Get("/ejBrownHall", ejBrownHall)
	mux.Get("/mailCenter", mailCenter)
	mux.Get("/orvisActivitiesCenter", orvisActivitiesCenter)
	mux.Get("/baseballField", baseballField)
	mux.Get("/softballField", softballField)

	//posts
	mux.Post("/data/createUser", createUser)
	mux.Post("/pioneerStadium", recordPioneerStadium)
	mux.Post("/taParishHall", recordTaParishHall)
	mux.Post("/studentDevelopmentCenter", recordStudentDevelopmentCenter)
	mux.Get("/studentLeadershipCenter", recordStudentLeadershipCenter)
	mux.Post("/hindelLibrary", recordHindelLibrary)
	mux.Post("/financialAid", recordFinancialAid)
	mux.Post("/pioneerCenter", recordPioneerCenter)
	mux.Post("/ejBrownHall", recordEjBrownHall)
	mux.Post("/mailCenter", recordMailCenter)
	mux.Post("/orvisActivitiesCenter", recordOrvisActivitiesCenter)
	mux.Post("/baseballField", recordBaseballField)
	mux.Post("/softballField", recordSoftballField)

	mux.NotFound(pageNotfound)
	return mux
}

// FILE HANDLERS
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

// PAGE HANDLERS
func index(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"index.tmpl",
		nil,
	)
}

func pioneerStadium(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"pioneerStadium.tmpl",
		nil,
	)
}

func recordPioneerStadium(writer http.ResponseWriter, request *http.Request) {

}

func taParishHall(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"taParishHall.tmpl",
		nil,
	)
}

func recordTaParishHall(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"taParishHall.tmpl",
		nil,
	)
}

func studentDevelopmentCenter(writer http.ResponseWriter, request *http.Request) {

}
func recordStudentDevelopmentCenter(writer http.ResponseWriter, request *http.Request) {

}

func hindelLibrary(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"hindelLibrary.tmpl",
		nil,
	)
}

func recordHindelLibrary(writer http.ResponseWriter, request *http.Request) {

}

func financialAid(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"financialAid.tmpl",
		nil,
	)
}

func recordFinancialAid(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"financialAid.tmpl",
		nil,
	)
}

func studentLeadershipCenter(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"studentLeadershipCenter.tmpl",
		nil,
	)
}

func recordStudentLeadershipCenter(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"studentLeadershipCenter.tmpl",
		nil,
	)
}

func pioneerCenter(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"pioneerCenter.tmpl",
		nil,
	)
}

func recordPioneerCenter(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"pioneerCenter.tmpl",
		nil,
	)
}

func ejBrownHall(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"ejBrownHall.tmpl",
		nil,
	)
}

func recordEjBrownHall(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"ejBrownHall.tmpl",
		nil,
	)
}

func mailCenter(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"mailCenter.tmpl",
		nil,
	)
}

func recordMailCenter(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"mailCenter.tmpl",
		nil,
	)
}

func orvisActivitiesCenter(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"orvisActivitiesCenter.tmpl",
		nil,
	)
}

func recordOrvisActivitiesCenter(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"orvisActivitiesCenter.tmpl",
		nil,
	)
}

func baseballField(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"baseballField.tmpl",
		nil,
	)
}

func recordBaseballField(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"baseballField.tmpl",
		nil,
	)
}

func softballField(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"softballField.tmpl",
		nil,
	)
}

func recordSoftballField(writer http.ResponseWriter, request *http.Request) {
	renderTemplate(
		writer,
		"softballField.tmpl",
		nil,
	)
}

func pageNotfound(writer http.ResponseWriter, response *http.Request) {
	renderTemplate(
		writer,
		"notFound.tmpl",
		nil,
	)
}

func createUser(writer http.ResponseWriter, request *http.Request) {
	var newUser types.NewUser
	err := json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		utils.LogErr("Client sent invalid JSON data", err, routerAppVars.LogPath)
		http.Error(writer, "Invalid JSON data", http.StatusBadRequest)
	} else {
		err = db.CreateNewUser(&newUser, routerAppVars)
		if err != nil {
			utils.LogErr("There was an error entering a new user into the database", err, routerAppVars.LogPath)
			http.Error(writer, "Error creating your profile please try again later", http.StatusInternalServerError)
		} else {
			writer.WriteHeader(http.StatusOK)
		}
	}
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

// createTemplateCache creates a template cache as a map
func createTemplateCache() map[string]*template.Template {
	schema := template.FuncMap{}
	templateCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./server/templates/pages/*.tmpl")
	if err != nil {
		utils.LogFatal("There was an error creating template cache", err, routerAppVars.LogPath)
	}

	for x := 0; x < len(pages); x++ {
		pageName := filepath.Base(pages[x])
		pageString, err := template.New(pageName).Funcs(schema).ParseFiles(pages[x])
		if err != nil {
			utils.LogFatal("There was an error creating template cache", err, routerAppVars.LogPath)
		}

		layoutString, err := filepath.Glob("./server/templates/templateBase.tmpl")
		if err != nil {
			utils.LogFatal("There was an error creating template cache", err, routerAppVars.LogPath)
		}

		if len(layoutString) > 0 {
			pageString, err = pageString.ParseGlob("./server/templates/templateBase.tmpl")
			if err != nil {
				utils.LogFatal("There was an error creating template cache", err, routerAppVars.LogPath)
			}
		}

		templateCache[pageName] = pageString
	}

	return templateCache
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
