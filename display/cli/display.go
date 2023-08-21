// # Display/CLI
//
// This package contains the code for the CLI display of the AWS accounts and
// OUs. It uses the tree structure generated in the generation package to
// display the accounts and OUs in a visually appealing way.
package cli

import (
	"fmt"

	"github.com/CentricaDevOps/aws-organizations-visualiser/generation"
)

// Display is a function that takes in the tree structure and displays it in a
// visually appealing way in the CLI.
func Display(tree *generation.OU) {
	fmt.Println("Displaying tree structure...")
}
