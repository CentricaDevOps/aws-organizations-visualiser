package cli

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/CentricaDevOps/aws-organizations-visualiser/generation"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"
)

func TestPrintTreeRecursive(t *testing.T) {
	// Create a tree to test with.
	tree := tree{
		referencedNode: &generation.OU{
			Name: "TestOU",
		},
		prefix:   endFork,
		spaces:   []string{},
		children: []tree{},
	}

	// Print the tree with detailed output.
	expectedOutput := "└─── TestOU (0)\n"
	output := captureOutput(func() {
		printTreeRecursive(tree, true)
	})
	require.Equal(t, expectedOutput, output, "Tree was not printed correctly")

	// Print the tree without detailed output.
	expectedOutput = "└─── TestOU\n"
	output = captureOutput(func() {
		printTreeRecursive(tree, false)
	})
	require.Equal(t, expectedOutput, output, "Tree was not printed correctly")

}

func TestPrintTreeRecursiveDeep(t *testing.T) {
	tree := tree{
		referencedNode: &generation.OU{
			Name: "TestOU",
		},
		prefix: endForkChild,
		spaces: []string{},
		children: []tree{
			{
				referencedNode: &generation.OU{
					Name: "TestOU2",
				},
				prefix: endForkChild,
				spaces: []string{
					"  ",
				},
				children: []tree{
					{
						referencedNode: &generation.OU{
							Name: "TestOU3",
						},
						prefix: endFork,
						spaces: []string{
							"  ",
							"  ",
						},
						children: []tree{},
					},
				},
			},
		},
	}

	// Print the tree with detailed output.
	expectedOutput := "" +
		"└─┬─ TestOU (0)\n" +
		"  └─┬─ TestOU2 (0)\n" +
		"    └─── TestOU3 (0)\n"
	output := captureOutput(func() {
		printTreeRecursive(tree, true)
	})
	require.Equal(t, expectedOutput, output, "Tree was not printed correctly")

	// Print the tree without detailed output.
	expectedOutput = "" +
		"└─┬─ TestOU\n" +
		"  └─┬─ TestOU2\n" +
		"    └─── TestOU3\n"
	output = captureOutput(func() {
		printTreeRecursive(tree, false)
	})
	require.Equal(t, expectedOutput, output, "Tree was not printed correctly")

}

func TestPrintTreeRecursiveWide(t *testing.T) {
	tree := tree{
		referencedNode: &generation.OU{
			Name: "TestOU",
		},
		prefix: endForkChild,
		spaces: []string{},
		children: []tree{
			{
				referencedNode: &generation.OU{
					Name: "TestOU2",
				},
				prefix: forkChild,
				spaces: []string{
					"  ",
				},
				children: []tree{
					{
						referencedNode: &generation.OU{
							Name: "TestOU3",
						},
						prefix: endFork,
						spaces: []string{
							"  ",
							"│ ",
						},
						children: []tree{},
					},
				},
			},
			{
				referencedNode: &generation.OU{
					Name: "TestOU4",
				},
				prefix: endForkChild,
				spaces: []string{
					"  ",
				},
				children: []tree{
					{
						referencedNode: &generation.OU{
							Name: "TestOU5",
						},
						prefix: endFork,
						spaces: []string{
							"  ",
							"  ",
						},
						children: []tree{},
					},
				},
			},
		},
	}

	// Print the tree with detailed output.
	expectedOutput := "" +
		"└─┬─ TestOU (0)\n" +
		"  ├─┬─ TestOU2 (0)\n" +
		"  │ └─── TestOU3 (0)\n" +
		"  └─┬─ TestOU4 (0)\n" +
		"    └─── TestOU5 (0)\n"
	output := captureOutput(func() {
		printTreeRecursive(tree, true)
	})
	require.Equal(t, expectedOutput, output, "Tree was not printed correctly")

	// Print the tree without detailed output.
	expectedOutput = "" +
		"└─┬─ TestOU\n" +
		"  ├─┬─ TestOU2\n" +
		"  │ └─── TestOU3\n" +
		"  └─┬─ TestOU4\n" +
		"    └─── TestOU5\n"
	output = captureOutput(func() {
		printTreeRecursive(tree, false)
	})
	require.Equal(t, expectedOutput, output, "Tree was not printed correctly")

}

func TestSetupTreeRecursiveSimple(t *testing.T) {
	// Create a tree to test with.
	tree := tree{
		referencedNode: &generation.OU{
			Name: "TestOU",
		},
		prefix:   "",
		spaces:   []string{},
		children: []tree{},
	}

	// Return a simple tree
	returnedTree := setupTreeRecursive(tree)

	// Check that the tree was returned correctly.
	require.Equal(t, tree.prefix, returnedTree.prefix, "Tree prefix was not set correctly")
	require.Len(t, returnedTree.spaces, 0, "Tree spaces was not set correctly")
	require.Len(t, returnedTree.children, 0, "Tree children was not set correctly")
}

func TestSetupTreeRecursiveComplicated(t *testing.T) {
	// Create a tree to test with.
	tree := tree{
		referencedNode: &generation.OU{
			Name: "TestOU",
			Children: []*generation.OU{
				{
					Name: "TestOU2",
					Children: []*generation.OU{
						{
							Name: "TestOU3",
						},
					},
				},
				{
					Name: "TestOU3",
					Children: []*generation.OU{
						{
							Name: "TestOU4",
						},
					},
				},
			},
		},
		prefix:   "",
		spaces:   []string{},
		children: []tree{},
	}

	returnedTree := setupTreeRecursive(tree)

	// Require the root to have 2 children.
	require.Len(t, returnedTree.children, 2, "Tree children was not set correctly")

	// Each of the children should have 1 child.
	require.Len(t, returnedTree.children[0].children, 1, "Tree children was not set correctly")
	require.Len(t, returnedTree.children[1].children, 1, "Tree children was not set correctly")

	// Each of the children's children should have no children.
	require.Len(t, returnedTree.children[0].children[0].children, 0, "Tree children was not set correctly")
	require.Len(t, returnedTree.children[1].children[0].children, 0, "Tree children was not set correctly")

	// Check that the Names are correct
	require.Equal(t, "TestOU", returnedTree.referencedNode.Name, "Tree referencedNode was not set correctly")
	require.Equal(t, "TestOU2", returnedTree.children[0].referencedNode.Name, "Tree referencedNode was not set correctly")
	require.Equal(t, "TestOU3", returnedTree.children[1].referencedNode.Name, "Tree referencedNode was not set correctly")
	require.Equal(t, "TestOU3", returnedTree.children[0].children[0].referencedNode.Name, "Tree referencedNode was not set correctly")
	require.Equal(t, "TestOU4", returnedTree.children[1].children[0].referencedNode.Name, "Tree referencedNode was not set correctly")

	// Check that the prefixes are correct
	require.Equal(t, "", returnedTree.prefix, "Tree prefix was not set correctly")
	require.Equal(t, forkChild, returnedTree.children[0].prefix, "Tree prefix was not set correctly")
	require.Equal(t, endForkChild, returnedTree.children[1].prefix, "Tree prefix was not set correctly")
	require.Equal(t, endFork, returnedTree.children[0].children[0].prefix, "Tree prefix was not set correctly")
	require.Equal(t, endFork, returnedTree.children[1].children[0].prefix, "Tree prefix was not set correctly")

	// Check that the spaces are correct
	require.Equal(t, []string{}, returnedTree.spaces, "Tree spaces was not set correctly")
	require.Equal(t, []string{emptySpacer}, returnedTree.children[0].spaces, "Tree spaces was not set correctly")
	require.Equal(t, []string{emptySpacer}, returnedTree.children[1].spaces, "Tree spaces was not set correctly")
	require.Equal(t, []string{emptySpacer, continueSpacer}, returnedTree.children[0].children[0].spaces, "Tree spaces was not set correctly")
	require.Equal(t, []string{emptySpacer, emptySpacer}, returnedTree.children[1].children[0].spaces, "Tree spaces was not set correctly")

}

func TestDisplayTree(t *testing.T) {
	ou := &generation.OU{
		Name: "TestOU",
		Children: []*generation.OU{
			{
				Name: "TestOU2",
				Children: []*generation.OU{
					{
						Name: "TestOU3",
					},
				},
			},
			{
				Name: "TestOU3",
				Children: []*generation.OU{
					{
						Name: "TestOU4",
					},
				},
			},
		},
		Accounts: []types.Account{
			{
				Name: aws.String("TestAccount"),
			},
			{
				Name: aws.String("TestAccount2"),
			},
		},
	}

	// Display the tree with detailed output.
	expectedOutput := "" +
		"└─┬─ TestOU (2)\n" +
		"  ├─┬─ TestOU2 (0)\n" +
		"  │ └─── TestOU3 (0)\n" +
		"  └─┬─ TestOU3 (0)\n" +
		"    └─── TestOU4 (0)\n"
	output := captureOutput(func() {
		displayTree(ou, true)
	})
	require.Equal(t, expectedOutput, output, "Tree was not displayed correctly")

	// Display the tree without detailed output.
	expectedOutput = "" +
		"└─┬─ TestOU\n" +
		"  ├─┬─ TestOU2\n" +
		"  │ └─── TestOU3\n" +
		"  └─┬─ TestOU3\n" +
		"    └─── TestOU4\n"
	output = captureOutput(func() {
		displayTree(ou, false)
	})
	require.Equal(t, expectedOutput, output, "Tree was not displayed correctly")
}

// captureOutput is a helper function to capture the output of a function.
// This is used to test the output of the display functions.
func captureOutput(f func()) string {
	// Store the old stdout and replace it with a pipe.
	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	// Run the function.
	f()
	// Close the pipe and restore stdout.
	w.Close()
	os.Stdout = original

	// Read the output from the pipe.
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		panic(err)
	}
	// Return the output.
	return buf.String()
}
