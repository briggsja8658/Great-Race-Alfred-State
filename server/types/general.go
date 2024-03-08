package types

import "golang.org/x/net/websocket"

type User struct {
	ID        int
	FirstName string
	LastName  int
}

type DirPicture struct {
	Location    string
	PictureName string
}

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}

type TeamData struct {
	TeamID        int
	TeamName      string
	StartingPoint int
	CreatedBy     int
}

type TeamDataJSON struct {
	TeamID   int    `json:"teamID"`
	TeamName string `json:"teamName"`
}

type NewUser struct {
	UserID   int
	UserName string
	TeamID   int
	TeamName string
}

type CheckNames struct {
	UserName string
	TeamName string
}

type WSServer struct {
	WsConnections map[*websocket.Conn]bool
}
