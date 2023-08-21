/*
# AWS Organizations Visualiser

TODO: Add description
*/
package main

import (
	"context"
	"fmt"
	"os"

	"log"

	"github.com/CentricaDevOps/aws-organizations-visualiser/display/cli"
	"github.com/CentricaDevOps/aws-organizations-visualiser/generation"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
)

// Logs is a struct that defines the logging configuration of the application
type Logs struct {
	Enabled bool
}

// Println is a function that prints a line to the console if logging is enabled
// for the application
func (l *Logs) Println(v ...interface{}) {
	if l.Enabled {
		log.Println(v...)
	}
}

// logs is a global variable that is used to log information about the
// application
var logs Logs

// init is called before main and is used to check the permissions of the user
// running the application and to set up the logging configuration
func setupLogging() {
	// Set up logging
	logs := Logs{Enabled: false}
	ll := os.Getenv("LOGS_ENABLED")
	if ll == "true" {
		logs.Enabled = true
		logs.Println("Logging initialised")
	}
}

// checkPermissions is a function that checks the permissions of the user running
// the application
func checkPermissions() (context.Context, *organizations.Client, error) {
	// Check permissions with a dry run of one command this application will run
	logs.Println("Checking permissions...")
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Error loading aws config, are you sure you are logged in?")
		logs.Println(err)
		return nil, nil, err
	}
	orgClient := organizations.NewFromConfig(cfg)
	_, err = orgClient.ListRoots(ctx, &organizations.ListRootsInput{})
	if err != nil {
		fmt.Println("You do not have permission to run the ListRoots command.")
		fmt.Println("Check that the account you are using has AWS Organizations enabled and that you are logged in with the correct permissions.")
		logs.Println(err)
		return nil, nil, err
	}
	logs.Println(" - Permissions OK!")
	return ctx, orgClient, nil
}

// main is the entry point of the application, it is called when the application
// is executed and is used to call the main logic of the application.
func main() {
	// STAGE 0: Set up the logging and check permissions
	setupLogging()
	ctx, cfg, err := checkPermissions()
	if err != nil {
		fmt.Println("Error checking permissions")
		logs.Println(err)
		return
	}

	// STAGE 1: Sort out the input flags
	// TODO: Decide on what flags to use and sort them out

	// STAGE 2: Run the main logic of the application to generate the data
	// structure
	context.Background()
	tree, err := generation.GenerateStructure(ctx, cfg)
	if err != nil {
		fmt.Println("Error generating structure")
		logs.Println(err)
		return
	}

	// STAGE 3: Determine the output format and output the data structure
	// TODO: Implement the output logic of the application
	cli.Display(tree)
}
