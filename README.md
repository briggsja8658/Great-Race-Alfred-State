# Alfred Great Race
- High level
    * Alfred Great Race is a Go web app for running The Great Race event. It serves HTML templates and static assets, and stores state in a local SQLite database.


# Features
- Primary Goal
    * Run an event flow with location-based progress, rotating targets, and basic team management.

- Participants can
    * Create or join teams
    * View current target and upcoming locations
    * Update progress through the race


# Tech
- Server
    * Go (net/http)
    * Templates with html/template
    * Routing with chi
- Client
    * Vanilla JS (ES modules)
    * Bootstrap
- DB
    * SQLite (local file)


# Prerequisites
- Go (1.18+ recommended)
- TLS cert/key files for local HTTPS
    * `server/keys/server.crt`
    * `server/keys/server.key`


# Run the project
- Clone the repository (if not already cloned)
    * git clone `https://github.com/briggsja8658/Alfred-Great-Race.git`
    * cd Alfred-Great-Race
- Start
    * go run .
- Go to
    * `https://localhost:6789` in your browser.


# Project structure (top-level)
- `main.go` - application entry point
- `client/` - front-end assets (JS, images, styles)
- `server/` - server code, templates, and static assets
- `server/db/` - SQLite schema and data access
- `server/router/` - HTTP handlers and routing
- `server/templates/` - HTML templates
- `server/pictures/` - image assets


# Notes
- Data and logs
    * SQLite database at `server/db/sqliteDB.sqlite` (created on first run)
    * Log file at `app.log`
- Git
    * `.gitignore` excludes logs, sqlite files, and TLS key/cert files.
