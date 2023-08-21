package cli

import (
	"fmt"
	"strings"

	"github.com/CentricaDevOps/aws-organizations-visualiser/generation"
)

// tree is a struct that holds the information needed to display the tree in the
// CLI. It is used with the setupTreeRecursive and printTreeRecursive functions.
type tree struct {
	referencedNode *generation.OU
	prefix         string
	spaces         []string
	children       []tree
}

// vars for displaying the tree
var (
	// wordSpacer is the gap between the prefix and the name of the Node
	wordSpacer = " "

	// empty Spacer is the gap between the start and the "┬" in the forks
	emptySpacer = "  "

	// continueSpacer is the continuation of a higher level Node which isn't
	// linked the current Node
	continueSpacer = "│ "

	// fork is the start of a child Node with siblings but no children
	fork = "├───" + wordSpacer

	// forkChild is the start of a child Node with siblings and children
	forkChild = "├─┬─" + wordSpacer

	// endFork is the end of a child Node with siblings but no children
	endFork = "└───" + wordSpacer

	// endForkChild is the end of a child Node with siblings and children
	endForkChild = "└─┬─" + wordSpacer
)

// setupTreeRecursive is a recursive function that sets up the tree struct with
// the correct information for displaying the tree in the CLI.
// This is done by determining the prefix, spaces and children of each Node.
//
// This is helpful as the output is {spaces[...]}{prefix}{name} therefore a
// large portion of this function is determining how many spaces to add and
// what the prefix should be.
//
// For example, if a node is the last child of its parent, it should have the
// endFork prefix. If it is not the last child, it should have the fork prefix.
func setupTreeRecursive(parentTree tree) tree {
	// Base case
	if len(parentTree.referencedNode.Children) == 0 {
		return tree{}
	}

	// If the parent is the last child, don't add a space
	spaceToAppend := continueSpacer
	if parentTree.prefix != fork {
		spaceToAppend = emptySpacer
	}

	// Change the prefix to the correct one as it has children
	if parentTree.prefix == fork {
		parentTree.prefix = forkChild
	} else if parentTree.prefix == endFork {
		parentTree.prefix = endForkChild
	}

	// For each child, add it to the children slice
	for _, child := range parentTree.referencedNode.Children {
		parentTree.children = append(
			parentTree.children,
			tree{
				referencedNode: child,
				prefix:         fork,
				spaces:         append(parentTree.spaces, spaceToAppend),
				children:       []tree{},
			},
		)
	}

	// Fix the prefix for the last child
	parentTree.children[len(parentTree.children)-1].prefix = endFork

	// For each child, recursively call this function
	for i := range parentTree.children {
		childValue := setupTreeRecursive(parentTree.children[i])
		if childValue.referencedNode != nil {
			parentTree.children[i] = childValue
		}
	}

	return parentTree
}

// printTreeRecursive is a recursive function that prints the tree of OUs in the
// CLI using the information from the tree struct. The detailed bool is used to
// determine whether to print the number of accounts in each OU.
func printTreeRecursive(display tree, detailed bool) {
	info := ""
	if detailed {
		info = fmt.Sprintf(
			" (%d)",
			len(display.referencedNode.Accounts),
		)
	}
	fmt.Printf("%s%s%s%s\n",
		strings.Join(display.spaces[:], ""),
		display.prefix,
		display.referencedNode.Name,
		info,
	)
	for _, child := range display.children {
		printTreeRecursive(child, detailed)
	}
}

/*
	 displayTree is a function that displays the tree of OUs in the CLI using a
	 tree structure similar to the one below:
	 └─┬─ TestOU (2)
	   ├─┬─ TestOU2 (5)
	   │ └─── TestOU3 (2)
	   └─┬─ TestOU4 (12)
		 └─── TestOU5 (2)

	 or without the detailed output:
	 └─┬─ TestOU
	   ├─┬─ TestOU2
	   │ └─── TestOU3
	   └─┬─ TestOU4
	     └─── TestOU5
*/
func displayTree(ouTree *generation.OU, detailed bool) {
	parentDisplay := tree{
		referencedNode: ouTree,
		prefix:         endFork,
		spaces:         []string{},
		children:       []tree{},
	}
	display := setupTreeRecursive(parentDisplay)
	printTreeRecursive(display, detailed)
}
