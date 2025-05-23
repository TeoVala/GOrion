package commands

import (
	envloader "GOrion/internal/env"
	handler "GOrion/internal/filehandler"
	"GOrion/internal/logging"
	"GOrion/internal/orm/gen"
	"GOrion/internal/router"
	"GOrion/internal/server"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CommandCallback defines the signature for command handler functions
type CommandCallback func(args []string)

// CommandRegistry manages all available commands
type Command struct {
	Handler     CommandCallback
	SubCommands map[string]*Command
}

type CommandRegistry struct {
	Commands map[string]*Command
}

func (cr *CommandRegistry) RegisterCommand(commandName string, handler CommandCallback) *Command {
	registeredCommand := &Command{
		Handler:     handler,
		SubCommands: make(map[string]*Command), // Create a list of subcommands for later
	}
	cr.Commands[commandName] = registeredCommand
	return registeredCommand
}

func (cr *CommandRegistry) RegisterSubCommand(parentCmdName string, subCommandName string, handler CommandCallback) *Command {
	parent, exists := cr.Commands[parentCmdName]
	if !exists {
		logging.LogAndPrint("Cannot register subcommand for non-existent parent command: %s", parentCmdName)
		os.Exit(1)
	}

	registeredSubCommand := &Command{
		Handler: handler,
	}

	parent.SubCommands[subCommandName] = registeredSubCommand
	return registeredSubCommand
}

func NewCommandRegistry() *CommandRegistry {
	registry := &CommandRegistry{
		Commands: make(map[string]*Command),
	}

	//Register main commands
	registry.RegisterCommand("runserver", registry.RunServer)
	registry.RegisterCommand("routelist", registry.RouteList)
	registry.RegisterCommand("genmodels", registry.GenerateModels)

	// Register make command with subcommands
	registry.RegisterCommand("make", registry.ParentNilCommand) // Register as empty for err purposes
	registry.RegisterSubCommand("make", "route", registry.CreateRoute)
	registry.RegisterSubCommand("make", "model", registry.CreateModel)
	registry.RegisterSubCommand("make", "handler", registry.CreateHandler)

	return registry
}

func (cr *CommandRegistry) PrintAllCommands() {
	logging.LogAndPrint("    Available commands:")

	for cmdName, _ := range cr.Commands {
		logging.LogAndPrint("      %s", cmdName)
	}
}

func (cr *CommandRegistry) PrintAllSubCommands(cmd string) {
	logging.LogAndPrint("    Available Subcommands:")

	for cmdName, _ := range cr.Commands[cmd].SubCommands {
		logging.LogAndPrint("      %s", cmdName)
	}
}

func (cr *CommandRegistry) RunServer(args []string) {
	var validArgs = [1]string{"port"}
	var foundInvalidArg bool
	var err = []string{}

	// Load projects environment variables
	env := envloader.LoadEnvVariables()

	var userSetPort string = env.AppPort

	var userValidArgs = make(map[string]string)

	for _, arg := range args {
		foundInvalidArg = true

		for _, prefix := range []string{"--", "-"} {
			if strings.HasPrefix(arg, prefix) {
				argParts := strings.SplitN(strings.TrimPrefix(arg, prefix), ":", 2)

				// Check if values passed
				if len(argParts) != 2 || argParts[1] == "" {
					break
				}

				// argParts[0] Name argParts[1] Value
				for _, validArg := range validArgs {
					if argParts[0] == validArg {
						foundInvalidArg = false
						userValidArgs[argParts[0]] = argParts[1]
						break
					}
				}
			}
		}
		if foundInvalidArg {
			err = append(err, arg)
		}
	}

	if len(err) >= 1 {
		logging.LogAndPrint("Encountered problem with command: %s", err)
		os.Exit(1)
	}

	for cmd, arg := range userValidArgs {
		switch cmd {
		case "port":
			// Convert port number to string
			// Common HTTP ports: 80 (standard), 443 (HTTPS)
			// Development server default: 8000
			// Recommended for development: 8000-8999
			// Ports below 1024 require root/admin privileges

			portNum, err := strconv.Atoi(arg)
			if err != nil {
				logging.LogAndPrint("Error converting to int: %s", err)
				os.Exit(1)
			}

			if portNum >= 0 && portNum <= 65535 {
				userSetPort = arg
				break
			}

			logging.LogAndPrint("Invalid port number: %s", arg)
			os.Exit(1)
		}
	}

	serverInstanse := server.NewServer(userSetPort)
	serverInstanse.ServerRun()
}

func (cr *CommandRegistry) RouteList(args []string) {
	server.ServerRunOnlyRoutes()
	router.GetAllRoutes()
}

func (cr *CommandRegistry) GenerateModels(args []string) {
	gen.Generate()
}

// This command is placed as a placeholder
// for functions that need subcommands to run
func (cr *CommandRegistry) ParentNilCommand(args []string) {
	logging.LogAndPrint("No subcommand submitted")
	var cmd string = os.Args[1]
	cr.PrintAllSubCommands(cmd)
	os.Exit(1)
}

func (cr *CommandRegistry) CreateRoute(args []string) {
	// fmt.Println(args)
	if len(args) == 0 {
		logging.LogAndPrint("Error: Unknown Subcommand: No valid name for route added")
		os.Exit(1)
	}

	var validArgs = [1]string{"name"}
	var foundInvalidArg bool
	var userValidArgs = make(map[string]string)
	var err = []string{}

	var routeName string

	// TODO Make this its own function and merge with runserver same code
	for _, arg := range args {
		foundInvalidArg = true

		for _, prefix := range []string{"--", "-"} {
			if strings.HasPrefix(arg, prefix) {
				argParts := strings.SplitN(strings.TrimPrefix(arg, prefix), ":", 2)

				// Check if values passed
				if len(argParts) != 2 || argParts[1] == "" {
					break
				}

				// argParts[0] Name argParts[1] Value
				for _, validArg := range validArgs {
					if argParts[0] == validArg {
						foundInvalidArg = false
						userValidArgs[argParts[0]] = argParts[1]
						break
					}
				}
			}
		}
		if foundInvalidArg {
			err = append(err, arg)
		}
	}

	for cmd, arg := range userValidArgs {
		switch cmd {
		case "name":
			routeName = arg
		}
	}

	handler.CreateRoute(routeName)
}

func (cr *CommandRegistry) CreateModel(args []string) {
	// TODO Create models doesnt work
	fmt.Println("Running create model with args:", args)
}

func (cr *CommandRegistry) CreateMiddleware(args []string) {
	// TODO Create middleware doesnt work
	fmt.Println("Running create middleware with args:", args)
}

func (cr *CommandRegistry) CreateHandler(args []string) {
	// fmt.Println(args)
	if len(args) == 0 {
		logging.LogAndPrint("Error: Unknown Subcommand: No valid name for handler added")
		os.Exit(1)
	}

	var validArgs = [1]string{"name"}
	var foundInvalidArg bool
	var userValidArgs = make(map[string]string)
	var err = []string{}

	var handlerName string

	// TODO Make this its own function and merge with runserver same code
	for _, arg := range args {
		foundInvalidArg = true

		for _, prefix := range []string{"--", "-"} {
			if strings.HasPrefix(arg, prefix) {
				argParts := strings.SplitN(strings.TrimPrefix(arg, prefix), ":", 2)

				// Check if values passed
				if len(argParts) != 2 || argParts[1] == "" {
					break
				}

				// argParts[0] Name argParts[1] Value
				for _, validArg := range validArgs {
					if argParts[0] == validArg {
						foundInvalidArg = false
						userValidArgs[argParts[0]] = argParts[1]
						break
					}
				}
			}
		}
		if foundInvalidArg {
			err = append(err, arg)
		}
	}

	for cmd, arg := range userValidArgs {
		switch cmd {
		case "name":
			handlerName = arg
		}
	}

	handler.CreateHandler(handlerName)
}

// Run executes the specified command with the given arguments
func (cr *CommandRegistry) ExecuteCommand(args []string) {
	action, params, err := cr.ParseCommand(args)
	if err != nil {
		logging.LogAndPrint("Error: Problem with running command")
		cr.PrintAllCommands()
		os.Exit(1)
	}
	action(params)
}

func (cr *CommandRegistry) ParseCommand(args []string) (CommandCallback, []string, error) {
	var subcommand string
	var isCreator bool = false

	var cmd string = os.Args[1]
	var parameters []string = os.Args[2:]

	if strings.Contains(os.Args[1], ":") {
		parts := strings.SplitN(os.Args[1], ":", 2)
		cmd = parts[0]
		subcommand = parts[1]
		isCreator = true
	}

	if isCreator {
		action, err := cr.Commands[cmd].SubCommands[subcommand]
		if !err {
			logging.LogAndPrint("Error: Unknown Subcommand: '%s'", subcommand)
			cr.PrintAllSubCommands(cmd)
			os.Exit(1)
		}

		return action.Handler, parameters, nil
	}

	action, err := cr.Commands[cmd]
	if !err {
		logging.LogAndPrint("Error: Unknown Command '%s'", cmd)
		cr.PrintAllCommands()
		os.Exit(1)
	}

	return action.Handler, parameters, nil
}
