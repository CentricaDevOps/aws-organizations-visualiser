package generation

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
)

// --- ListRoots ---------------------------------------------------------------
// ListRoots is an interface for the ListRoots function from the organizations
// service in the AWS SDK that allows for mocking.
type ListRoots interface {
	ListRoots(
		ctx context.Context,
		params *organizations.ListRootsInput,
		optFns ...func(*organizations.Options),
	) (
		*organizations.ListRootsOutput,
		error,
	)
}

type ListRootsMock struct {
	ListRootsFunc func(
		ctx context.Context,
		params *organizations.ListRootsInput,
		optFns ...func(*organizations.Options),
	) (
		*organizations.ListRootsOutput,
		error,
	)
}

func (m *ListRootsMock) ListRoots(
	ctx context.Context,
	params *organizations.ListRootsInput,
	optFns ...func(*organizations.Options),
) (
	*organizations.ListRootsOutput,
	error,
) {
	return m.ListRootsFunc(ctx, params, optFns...)
}

// --- ListOrganizationalUnitsForParent ----------------------------------------
// ListOrganizationalUnitsForParent is an interface for the organizations
// ListOrganizationalUnitsForParent function in the AWS SDK that allows for
// mocking.
type ListOrganizationalUnitsForParent interface {
	ListOrganizationalUnitsForParent(
		ctx context.Context,
		params *organizations.ListOrganizationalUnitsForParentInput,
		optFns ...func(*organizations.Options),
	) (
		*organizations.ListOrganizationalUnitsForParentOutput,
		error,
	)
}

type ListOrganizationalUnitsForParentMock struct {
	ListOrganizationalUnitsForParentFunc func(
		ctx context.Context,
		params *organizations.ListOrganizationalUnitsForParentInput,
		optFns ...func(*organizations.Options),
	) (
		*organizations.ListOrganizationalUnitsForParentOutput,
		error,
	)
}

func (m *ListOrganizationalUnitsForParentMock) ListOrganizationalUnitsForParent(
	ctx context.Context,
	params *organizations.ListOrganizationalUnitsForParentInput,
	optFns ...func(*organizations.Options),
) (
	*organizations.ListOrganizationalUnitsForParentOutput,
	error,
) {
	return m.ListOrganizationalUnitsForParentFunc(ctx, params, optFns...)
}

// --- ListAccountsForParentPaginator ------------------------------------------
// ListAccountsForParentPaginator is an interface for the organizations
// ListAccountsForParentPaginator function in the AWS SDK that allows for
// mocking, this is a paginator function so it has a HasMorePages and NextPage
// function.
type ListAccountsForParentPaginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context,
		optFns ...func(*organizations.Options),
	) (
		*organizations.ListAccountsForParentOutput,
		error,
	)
}

type mockListAccountsForParentPager struct {
	PageNum int
	Pages   []organizations.ListAccountsForParentOutput
	Error   error
}

func (m *mockListAccountsForParentPager) HasMorePages() bool {
	return m.PageNum < len(m.Pages)
}

func (m *mockListAccountsForParentPager) NextPage(
	ctx context.Context,
	optFns ...func(*organizations.Options),
) (
	*organizations.ListAccountsForParentOutput,
	error,
) {
	if m.Error != nil {
		return nil, m.Error
	}
	if !m.HasMorePages() {
		return nil, fmt.Errorf("no more pages")
	}
	page := m.Pages[m.PageNum]
	m.PageNum++
	return &page, nil
}
