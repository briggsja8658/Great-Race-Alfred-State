module greatRace

go 1.18

require (
	//Things manually imported
	github.com/fatih/color v1.15.0 // direct
	github.com/go-chi/chi v1.5.4 // direct
	github.com/mattn/go-sqlite3 v1.14.17 // direct
)

require (
	github.com/go-chi/cors v1.2.1 // indirect

	//Dependencies for the packages imported
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	golang.org/x/sys v0.6.0 // indirect
)

require (
	github.com/r3labs/sse/v2 v2.10.0 // indirect
	golang.org/x/net v0.0.0-20191116160921-f9c825593386 // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
)
