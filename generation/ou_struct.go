package generation

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

// --- OU ----------------------------------------------------------------------
// OU is a struct that represents an OU in the AWS Organizations structure.
// It can be used to represent the entire structure or a substructure in the
// style of a tree.
type OU struct {
	Id       string
	Name     string
	Children []*OU
	Accounts []types.Account
}

// addChildren adds the given OUs to the OU's children slice.
func (o *OU) addChildren(children []*OU) {
	o.Children = append(o.Children, children...)
}

// GetChildren returns the OU's child OUs.
func (o *OU) GetChildren() []*OU {
	return o.Children
}

// GetId returns the OU's ID.
func (o *OU) GetId() string {
	return o.Id
}

// GetName returns the OU's name.
func (o *OU) GetName() string {
	return o.Name
}

// GetAccounts returns a list of the accounts in the OU.
func (o *OU) GetAccounts() []types.Account {
	return o.Accounts
}

func (parent *OU) fillOuTree(ctx context.Context, api ListOrganizationalUnitsForParent) error {
	// Get the OUs for the parent OU.
	ous, err := getOUsForParent(ctx, api, parent.Id)
	if err != nil {
		return err
	}
	if len(ous) == 0 {
		return nil
	}

	// Recursively fill the tree with the OUs.
	for i := range ous {
		err := ous[i].fillOuTree(ctx, api)
		if err != nil {
			return err
		}
	}

	// Append the OUs to the parent OU.
	parent.addChildren(ous)

	return nil
}

// fillAccountsRecursive fills the OU tree with the accounts in the OUs.
func (parent *OU) fillAccountsRecursive(ctx context.Context, api *organizations.Client) (*OU, error) {
	// Get the accounts for the parent OU.
	accounts, err := getAccountsFromOU(ctx, api, parent.Id, parent.Name)
	if err != nil {
		return nil, err
	}
	parent.Accounts = accounts

	// Recursively fill the tree with the accounts.
	for i := range parent.Children {
		ou, err := parent.Children[i].fillAccountsRecursive(ctx, api)
		if err != nil {
			return nil, err
		}
		parent.Children[i] = ou
	}

	return parent, nil
}
