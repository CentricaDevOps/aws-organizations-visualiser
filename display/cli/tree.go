package cli

import (
	"fmt"
	"strings"

	"github.com/CentricaDevOps/aws-organizations-visualiser/generation"
)

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

func setupTreeRecursive(parentTree tree) tree {
	// If there are no children, return true
	if len(parentTree.referencedNode.Children) == 0 {
		return tree{}
	}

	// If the parent is the last child, don't add a space
	spaceToAppend := continueSpacer
	if parentTree.prefix != fork {
		spaceToAppend = emptySpacer
	}

	// If its gotten here, the parent has children therefore needs to change its
	// prefix to a fork with a child
	if parentTree.prefix == fork {
		parentTree.prefix = forkChild
	} else if parentTree.prefix == endFork {
		parentTree.prefix = endForkChild
	}

	for _, child := range parentTree.referencedNode.Children {
		childDisplay := tree{
			referencedNode: child,
			prefix:         fork,
			spaces:         append(parentTree.spaces, spaceToAppend),
			children:       []tree{},
		}
		parentTree.children = append(
			parentTree.children,
			childDisplay,
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
