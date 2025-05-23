package main

import (
	"GOrion/internal/commands"
	"GOrion/internal/logging"
	"log"
	"os"
)

type App struct {
	logFile         *os.File
	commandRegistry *commands.CommandRegistry
}

func NewApp() (*App, error) {
	// Initialize logging
	logFile, err := logging.InitLog()
	if err != nil {
		log.Fatalf("Failed to initialize logging: %v", err)
	}

	return &App{
		logFile:         logFile,
		commandRegistry: commands.NewCommandRegistry(), // Initialize command registry
	}, nil
}

func (app *App) Close() {
	// Defer closing the log file until main exits.
	logging.CloseLogFile(app.logFile)
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.Close()

	// Command validation
	if len(os.Args) < 2 {
		logging.LogAndPrint("Usage: go run manage.go <command> [arguments]")
		app.commandRegistry.PrintAllCommands()

		os.Exit(1)
	}

	app.commandRegistry.ExecuteCommand(os.Args)
}
