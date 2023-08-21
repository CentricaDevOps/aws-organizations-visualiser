package generation

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/stretchr/testify/require"
)

// TestOUCreation tests the creation of an OU struct.
func TestOuCreation(t *testing.T) {
	// Create an OU struct.
	ou := &OU{
		Id:   "ou-1234",
		Name: "TestOU",
	}

	// Check that the OU was created correctly.
	require.Equal(t, "ou-1234", ou.Id, "OU ID was not set correctly")
	require.Equal(t, "TestOU", ou.Name, "OU name was not set correctly")
	require.Equal(t, 0, len(ou.Children), "OU children was not set correctly")
	require.Equal(t, 0, len(ou.Accounts), "OU accounts was not set correctly")
}

// TestOuAddChildren tests the addChildren method of the OU struct.
func TestOuAddChildren(t *testing.T) {
	// Create an OU struct.
	ou := &OU{
		Id:   "ou-1234",
		Name: "TestOU",
	}

	// Create a child OU struct.
	child := &OU{
		Id:   "ou-5678",
		Name: "TestChildOU",
	}

	// Add the child OU to the parent OU.
	ou.addChildren([]*OU{child})

	// Check that the child OU was added correctly.
	require.Equal(t, 1, len(ou.Children), "OU children was not set correctly")
	require.Equal(t, child, ou.Children[0], "OU children was not set correctly")

	// Add another child OU to the parent OU.
	child2 := &OU{
		Id:   "ou-9012",
		Name: "TestChildOU2",
	}

	// Add the child OU to the parent OU.
	ou.addChildren([]*OU{child2})

	// Check that the child OU was added correctly.
	require.Equal(t, 2, len(ou.Children), "OU children was not set correctly")
	require.Equal(t, child2, ou.Children[1], "OU children was not set correctly")

	// Check that the ids of the children are correct.
	require.Equal(t, "ou-5678", ou.Children[0].Id, "OU children was not set correctly")
	require.Equal(t, "ou-9012", ou.Children[1].Id, "OU children was not set correctly")

	// Check that the names of the children are correct.
	require.Equal(t, "TestChildOU", ou.Children[0].Name, "OU children was not set correctly")
	require.Equal(t, "TestChildOU2", ou.Children[1].Name, "OU children was not set correctly")
}

// TestOUGetters tests the getters of the OU struct.
func TestOuGetters(t *testing.T) {
	// Create an OU struct.
	ou := &OU{
		Id:   "ou-1234",
		Name: "TestOU",
	}

	// Check that the getters return the correct values.
	require.Equal(t, "ou-1234", ou.GetId(), "OU ID getter returned incorrect value")
	require.Equal(t, "TestOU", ou.GetName(), "OU name getter returned incorrect value")
	require.Equal(t, 0, len(ou.GetChildren()), "OU children getter returned incorrect value")
	require.Equal(t, 0, len(ou.GetAccounts()), "OU accounts getter returned incorrect value")
}

// TestOuFillOuTree tests the fillOuTree method of the OU struct.
func TestOuFillOuTree(t *testing.T) {
	// Create a mock API.
	mockClient := ListOrganizationalUnitsForParentMock{
		ListOrganizationalUnitsForParentFunc: func(
			ctx context.Context,
			params *organizations.ListOrganizationalUnitsForParentInput,
			optFns ...func(*organizations.Options),
		) (
			*organizations.ListOrganizationalUnitsForParentOutput,
			error,
		) {
			// Check that the correct parent ID was passed in
			// If not, return an empty list of OUs.
			if *params.ParentId != "ou-1234" {
				return &organizations.ListOrganizationalUnitsForParentOutput{
					OrganizationalUnits: []types.OrganizationalUnit{},
				}, nil
			}

			// Create the mock output
			output := &organizations.ListOrganizationalUnitsForParentOutput{
				OrganizationalUnits: []types.OrganizationalUnit{
					{
						Id:   aws.String("ou-5678"),
						Name: aws.String("TestChildOU"),
					},
				},
			}

			return output, nil
		},
	}

	// Create an OU struct.
	ou := &OU{
		Id:   "ou-1234",
		Name: "TestOU",
	}

	// Fill the OU tree.
	ctx := context.Background()
	err := ou.fillOuTree(ctx, &mockClient)
	require.NoError(t, err, "fillOuTree returned an error")

	// Check that the child OU was added correctly.
	require.Equal(t, 1, len(ou.Children), "OU children was not set correctly")
	require.Equal(t, "ou-5678", ou.Children[0].Id, "OU children was not set correctly")
	require.Equal(t, "TestChildOU", ou.Children[0].Name, "OU children was not set correctly")
}
