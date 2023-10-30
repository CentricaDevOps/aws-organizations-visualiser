// # Generation Package
//
// Package generation provides the generation of the code for the organizations
// visualisation application.
//
// This package has a single exported function, GenerateStructure, which takes in
// an Organizations Client and returns a custom tree structure that contains all
// the information about the organization.

// It also provides some helper functions for the struct to make it easier to
// work with.
package generation

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

// GenerateStructure takes in an Organizations Client and returns a custom tree
// structure that contains all the information about the organization.
func GenerateStructure(ctx context.Context, orgClient *organizations.Client) (*OU, error) {
	// Get the root of the organization
	rootId, err := getRootId(ctx, orgClient)
	if err != nil {
		return nil, err
	}

	// Initialise the tree
	tree := &OU{
		Id:       rootId,
		Name:     "Root",
		Children: []*OU{},
		Accounts: []types.Account{},
	}

	// Get the OUs
	err = tree.fillOuTree(ctx, orgClient)
	if err != nil {
		return nil, err
	}

	// Get the accounts
	tree, err = tree.fillAccountsRecursive(ctx, orgClient)
	if err != nil {
		return nil, err
	}

	// Remove suspended accounts
	tree = tree.removeSuspendedAccounts()

	return tree, nil
}
