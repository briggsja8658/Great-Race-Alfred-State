package types

type Location struct {
	LocationID   int
	PositionID   int
	LocationName string
	ReadableName string
}

type Team struct {
	TeamID        int
	TeamName      string
	StartingPoint int
	CreatedBy     int
}

type UserEntry struct {
	TeamMemberID int
	Name         string
	TeamID       int
}

type Progess struct {
	ProgressID    int
	Finished      bool
	TimeStarted   int
	TimeCompleted int
	StartingPoint int
	CurrentPoint  int
	EndPoint      int
	TeamID        int
}
