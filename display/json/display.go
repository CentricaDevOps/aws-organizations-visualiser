// # Display/JSON
//
// This package contains the code for the JSON display of the AWS accounts and
// OUs. It uses the tree structure generated in the generation package to
// create a JSON representation of the accounts and OUs.
package json

import (
	"os"

	"github.com/CentricaDevOps/aws-organizations-visualiser/generation"
)

// Create is a function that takes in the tree structure and creates a JSON
// representation of it.
func Create(tree *generation.OU) ([]byte, error) {
	// Use the existing structure to create a JSON representation of the tree.
	return tree.ToJSON()
}

// OutputToFile is a function that takes in the json representaiton of the tree and
// outputs it to the given file
func OutputToFile(json []byte, filename string) error {
	// Open the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the JSON to the file
	_, err = file.Write(json)
	if err != nil {
		return err
	}
	return nil
}
