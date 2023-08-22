// # AWS Organizations Visualiser
//
// This is a tool that can be used to visualise the structure of an AWS
// Organizations structure. It can be used to generate a JSON representation of the
// structure or to display the structure in the CLI.
//
// ## Usage
//
// This tool generates a structure that represents the AWS Organizations structure
// and then, based on the flags passed in, either displays the structure in the CLI
// or outputs the structure to a JSON file or both.
//
// ### Flags
//
// Usage:
//
//	aws-organizations-visualiser [flags]
//
// Flags:
//
//	-include-json
//	      Include the JSON representation of the AWS Organizations structure in the output (default true)
//	-include-visual
//	      Include the visual representation of the AWS Organizations structure in the output (default true)
//	-o string
//	      The output file for the JSON representation of the AWS Organizations structure (default "output.json")
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/CentricaDevOps/aws-organizations-visualiser/display/cli"
	"github.com/CentricaDevOps/aws-organizations-visualiser/display/json"
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
		fmt.Println(v...)
	}
}

// logs is a global variable that is used to log information about the
// application
var logs Logs

// init is called before main and is used to check the permissions of the user
// running the application and to set up the logging configuration
func setupLogging(logging string) {
	// Set up logging
	value, err := strconv.ParseBool(logging)
	if err != nil {
		value = false
	}
	logs = Logs{Enabled: value}
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
	// STAGE 1: Sort out the input flags
	jsonPtr := flag.Bool("include-json", true, "Include the JSON representation of the AWS Organizations structure in the output")
	visualPtr := flag.Bool("include-visual", true, "Include the visual representation of the AWS Organizations structure in the output")
	outputPtr := flag.String("o", "output.json", "The output file for the JSON representation of the AWS Organizations structure")
	flag.Parse()

	// STAGE 2: Set up the logging and check permissions
	ll := os.Getenv("LOGS_ENABLED")
	setupLogging(ll)
	ctx, cfg, err := checkPermissions()
	if err != nil {
		fmt.Println("Error checking permissions")
		logs.Println(err)
		return
	}
	logs.Println("Output file:", *outputPtr)

	// STAGE 3: Run the main logic of the application to generate the data
	// structure
	context.Background()
	tree, err := generation.GenerateStructure(ctx, cfg)
	if err != nil {
		fmt.Println("Error generating structure")
		logs.Println(err)
		return
	}

	// STAGE 4: Determine the output format and output the data structure
	// If no output format is specified, exit
	if !*visualPtr && !*jsonPtr {
		fmt.Println("No output format specified, exiting...")
		return
	}
	// If the visual output format is specified, display the data structure on
	// the CLI
	if *visualPtr {
		cli.Display(tree)
	}

	// If the JSON output format is specified, output the data structure to a
	// JSON file with the given name
	if *jsonPtr {
		jsonTree, err := json.Create(tree)
		if err != nil {
			fmt.Println("Error generating JSON")
			logs.Println(err)
			return
		}
		err = json.OutputToFile(jsonTree, *outputPtr)
		if err != nil {
			fmt.Println("Error outputting JSON to file")
			logs.Println(err)
			return
		}
	}
}
