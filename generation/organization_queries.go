package generation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

func getOUsForParent(ctx context.Context, api ListOrganizationalUnitsForParent, parentId string) ([]*OU, error) {
	// Retry 5 times if the API call fails due to rate limits
	for i := 0; i < 5; i++ {
		// Get the child OUs of the root OU.
		ouList, err := api.ListOrganizationalUnitsForParent(ctx, &organizations.ListOrganizationalUnitsForParentInput{
			ParentId: &parentId,
		})
		if err != nil {
			if strings.Contains(err.Error(), "exceeded maximum number of attempts") {
				time.Sleep(5 * time.Second)
				continue
			}
			return nil, err
		}
		ous := make([]*OU, len(ouList.OrganizationalUnits))
		for i, ou := range ouList.OrganizationalUnits {
			ous[i] = &OU{
				Id:   *ou.Id,
				Name: *ou.Name,
			}
		}
		return ous, nil
	}
	return nil, fmt.Errorf("failed to get OUs for parent %s, most likely due to rate limits", parentId)
}

// getRootID gets the ID of the root OU.
func getRootId(ctx context.Context, api ListRoots) (string, error) {
	// Retry 5 times if the API call fails due to rate limits
	for i := 0; i < 5; i++ {
		// Get the root OU id.
		rootOU, err := api.ListRoots(ctx, &organizations.ListRootsInput{})
		if err != nil {
			if strings.Contains(err.Error(), "exceeded maximum number of attempts") {
				time.Sleep(5 * time.Second)
				continue
			}
			return "", err
		}

		return *rootOU.Roots[0].Id, nil
	}
	return "", fmt.Errorf("failed to get root OU ID, most likely due to rate limits")
}

// GetAccountsFromOU gets a list of aws accounts from an OU name.
func getAccountsFromOU(ctx context.Context, svc *organizations.Client, ouId string, ouBlock string) ([]types.Account, error) {
	// Get the child accounts of the parameter OU.
	paginator := organizations.NewListAccountsForParentPaginator(svc, &organizations.ListAccountsForParentInput{
		ParentId: &ouId,
	})
	accounts, err := getAllAccountsFromOUID(ctx, paginator, ouId)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// getAllAccountsFromOUID gets all the accounts from an OU ID.
func getAllAccountsFromOUID(ctx context.Context, svc ListAccountsForParentPaginator, ou_id string) ([]types.Account, error) {
	// Get the child accounts of the given ou.
	account_list, err := svc.NextPage(ctx)
	if err != nil {
		return nil, err
	}

	// Get the next page of accounts if there is one.
	for svc.HasMorePages() {
		output, err := svc.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		// Append the accounts to the account list.
		account_list.Accounts = append(account_list.Accounts, output.Accounts...)
	}

	return account_list.Accounts, nil
}
