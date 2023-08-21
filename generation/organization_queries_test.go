package generation

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/stretchr/testify/require"
)

// TestGetOUsForParent tests the getOUsForParent function calls the mock
// correctly and returns the data correctly.
func TestGetOUsForParent(t *testing.T) {
	// This tests that the getOUForParent function returns the correct OUs for a
	// given parent OU using a mock library

	// Create the mock
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
			if *params.ParentId != "ou-1234" {
				return nil, fmt.Errorf("expected parent ID to be ou-1234, got %s", *params.ParentId)
			}

			// Create the mock output
			output := &organizations.ListOrganizationalUnitsForParentOutput{
				OrganizationalUnits: []types.OrganizationalUnit{
					{
						Id:   aws.String("ou-1234"),
						Name: aws.String("TestOU"),
					},
				},
			}

			return output, nil
		},
	}

	// Call the function
	ctx := context.Background()
	ous, err := getOUsForParent(ctx, &mockClient, "ou-1234")
	require.NoError(t, err)

	// Check that the correct OU was returned
	require.Equal(t, 1, len(ous))
	require.Equal(t, "ou-1234", ous[0].Id)
	require.Equal(t, "TestOU", ous[0].Name)
}

// TestGetOUsForParentError tests the getOUsForParent function calls the mock
// correctly and returns the error correctly
func TestGetOUsForParentError(t *testing.T) {
	// This tests that the getOUForParent function returns the correct OUs for a
	// given parent OU using a mock library

	// Create the mock
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
			if *params.ParentId != "ou-1234" {
				return nil, fmt.Errorf("expected parent ID to be ou-1234, got %s", *params.ParentId)
			}

			// Create the mock output
			output := &organizations.ListOrganizationalUnitsForParentOutput{
				OrganizationalUnits: []types.OrganizationalUnit{
					{
						Id:   aws.String("ou-1234"),
						Name: aws.String("TestOU"),
					},
				},
			}

			return output, nil
		},
	}

	// Call the function
	ctx := context.Background()
	ous, err := getOUsForParent(ctx, &mockClient, "ou-4321")
	require.Error(t, err)
	require.Nil(t, ous)
}

// TestGetOUsForParentError tests the getOUsForParent function that it doesn't
// loop forever if rate limits occur.
func TestGetOUsForParentErrorRateLimit(t *testing.T) {
	// This tests that the getOUForParent function returns the correct OUs for a
	// given parent OU using a mock library

	// Check current time:
	timeBefore := time.Now()
	// Create the mock
	mockClient := ListOrganizationalUnitsForParentMock{
		ListOrganizationalUnitsForParentFunc: func(
			ctx context.Context,
			params *organizations.ListOrganizationalUnitsForParentInput,
			optFns ...func(*organizations.Options),
		) (
			*organizations.ListOrganizationalUnitsForParentOutput,
			error,
		) {
			return nil, fmt.Errorf("exceeded maximum number of attempts")
		},
	}

	// Call the function
	ctx := context.Background()
	ous, err := getOUsForParent(ctx, &mockClient, "ou-4321")
	require.Error(t, err)
	require.Nil(t, ous)

	// Check that the function didn't loop forever
	timeAfter := time.Now()

	// The code loops 5 times and waits 5 seconds each time, so the time taken
	// should be less than 30 seconds
	expectedTimeDiff := 30 * time.Second
	require.True(t, timeAfter.Sub(timeBefore) < expectedTimeDiff)
}

// TestGetRootId tests the getRootId function calls the mock correctly and
// returns the data correctly.
func TestGetRootId(t *testing.T) {
	// This tests that the getRootId function returns the correct root ID using a
	// mock library

	// Create the mock
	mockClient := ListRootsMock{
		ListRootsFunc: func(
			ctx context.Context,
			params *organizations.ListRootsInput,
			optFns ...func(*organizations.Options),
		) (
			*organizations.ListRootsOutput,
			error,
		) {
			// Create the mock output
			output := &organizations.ListRootsOutput{
				Roots: []types.Root{
					{
						Id: aws.String("r-1234"),
					},
				},
			}

			return output, nil
		},
	}

	// Call the function
	ctx := context.Background()
	rootId, err := getRootId(ctx, &mockClient)
	require.NoError(t, err)

	// Check that the correct root ID was returned
	require.Equal(t, "r-1234", rootId)
}

// TestGetRootIdError tests the getRootId function calls the mock correctly and
// returns the error correctly.
func TestGetRootIdError(t *testing.T) {
	// This tests that the getRootId function returns the correct root ID using a
	// mock library

	// Create the mock
	mockClient := ListRootsMock{
		ListRootsFunc: func(
			ctx context.Context,
			params *organizations.ListRootsInput,
			optFns ...func(*organizations.Options),
		) (
			*organizations.ListRootsOutput,
			error,
		) {
			// Return an error
			return nil, fmt.Errorf("Testing Error")
		},
	}

	// Call the function
	ctx := context.Background()
	rootId, err := getRootId(ctx, &mockClient)
	require.Error(t, err)
	require.Equal(t, "", rootId)
}

// TestGetAllAccountsFromOUID tests the getAccountsFromOU function calls the mock
// correctly and returns the data correctly.
func TestGetAllAccountsFromOUID(t *testing.T) {
	// This tests that the getAccountsFromOU function returns the correct accounts
	// for a given OU using a mock library

	// Create the mock
	mockClient := mockListAccountsForParentPager{
		PageNum: 0,
		Pages: []organizations.ListAccountsForParentOutput{
			{
				Accounts: []types.Account{
					{
						Id:   aws.String("a-1234"),
						Name: aws.String("TestAccount"),
					},
				},
			},
			{
				Accounts: []types.Account{
					{
						Id:   aws.String("a-5678"),
						Name: aws.String("TestAccount2"),
					},
				},
			},
		},
	}

	// Call the function
	ctx := context.Background()
	accounts, err := getAllAccountsFromOUID(ctx, &mockClient, "ou-1234")
	require.NoError(t, err)

	// Check that the correct accounts were returned
	require.Equal(t, 2, len(accounts))
	require.Equal(t, "a-1234", *accounts[0].Id)
	require.Equal(t, "TestAccount", *accounts[0].Name)
	require.Equal(t, "a-5678", *accounts[1].Id)
	require.Equal(t, "TestAccount2", *accounts[1].Name)
}

// TestGetAllAccountsFromOUIDError tests the getAccountsFromOU function calls the
// mock correctly and returns the error correctly.
func TestGetAllAccountsFromOUIDError(t *testing.T) {
	// This tests that the getAccountsFromOU function returns the correct accounts
	// for a given OU using a mock library

	// Create the mock
	mockClient := mockListAccountsForParentPager{
		PageNum: 0,
		Pages: []organizations.ListAccountsForParentOutput{
			{
				Accounts: []types.Account{
					{
						Id:   aws.String("a-1234"),
						Name: aws.String("TestAccount"),
					},
				},
			},
		},
		Error: fmt.Errorf("Testing Error"),
	}

	// Call the function
	ctx := context.Background()
	accounts, err := getAllAccountsFromOUID(ctx, &mockClient, "ou-1234")
	require.Error(t, err)
	require.Nil(t, accounts)
}
