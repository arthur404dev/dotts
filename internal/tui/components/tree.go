package components

import (
	"github.com/arthur404dev/dotts/internal/tui/theme"
	"github.com/charmbracelet/lipgloss"
)

// TreeNode represents a node in a hierarchical tree structure.
type TreeNode struct {
	ID       string
	Label    string
	Icon     string
	Children []TreeNode
	Expanded bool
}

// Tree is a hierarchical tree view component for displaying nested data.
type Tree struct {
	theme      *theme.Theme
	nodes      []TreeNode
	indentSize int
	showGuides bool
}

// NewTree creates a new Tree component with the given theme.
func NewTree(t *theme.Theme) *Tree {
	return &Tree{
		theme:      t,
		nodes:      []TreeNode{},
		indentSize: 2,
		showGuides: true,
	}
}

// SetNodes sets the root nodes of the tree.
func (t *Tree) SetNodes(nodes []TreeNode) *Tree {
	t.nodes = nodes
	return t
}

// AddNode appends a root node to the tree.
func (t *Tree) AddNode(node TreeNode) *Tree {
	t.nodes = append(t.nodes, node)
	return t
}

// SetIndentSize sets the number of spaces used for indenting child nodes.
func (t *Tree) SetIndentSize(size int) *Tree {
	if size >= 0 {
		t.indentSize = size
	}
	return t
}

// SetShowGuides enables or disables the tree guide lines (├─, └─, │).
func (t *Tree) SetShowGuides(show bool) *Tree {
	t.showGuides = show
	return t
}

// View renders the tree as a formatted string.
func (tr *Tree) View() string {
	var lines []string
	for i, node := range tr.nodes {
		isLast := i == len(tr.nodes)-1
		lines = append(lines, tr.renderNode(node, "", isLast)...)
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (tr *Tree) renderNode(node TreeNode, prefix string, isLast bool) []string {
	t := tr.theme
	var lines []string

	var branch, continuation string
	if tr.showGuides {
		if isLast {
			branch = "└─"
			continuation = "  "
		} else {
			branch = "├─"
			continuation = "│ "
		}
	} else {
		branch = ""
		continuation = ""
	}

	guideStyle := lipgloss.NewStyle().Foreground(t.FgSubtle)
	labelStyle := lipgloss.NewStyle().Foreground(t.FgBase)

	var label string
	if node.Icon != "" {
		label = node.Icon + " " + node.Label
	} else {
		label = node.Label
	}

	if len(node.Children) > 0 {
		if node.Expanded {
			label = "▼ " + label
		} else {
			label = "▶ " + label
		}
	}

	line := prefix
	if tr.showGuides && prefix != "" || branch != "" {
		line += guideStyle.Render(branch) + " "
	}
	line += labelStyle.Render(label)

	lines = append(lines, line)

	if node.Expanded && len(node.Children) > 0 {
		childPrefix := prefix
		if tr.showGuides {
			childPrefix += guideStyle.Render(continuation)
			for i := 0; i < tr.indentSize; i++ {
				childPrefix += " "
			}
		} else {
			for i := 0; i < tr.indentSize+2; i++ {
				childPrefix += " "
			}
		}

		for i, child := range node.Children {
			childIsLast := i == len(node.Children)-1
			lines = append(lines, tr.renderNode(child, childPrefix, childIsLast)...)
		}
	}

	return lines
}

// SimpleTree creates a tree from a simple parent-children map.
func SimpleTree(t *theme.Theme, data map[string][]string) *Tree {
	tree := NewTree(t)
	for parent, children := range data {
		node := TreeNode{
			ID:       parent,
			Label:    parent,
			Expanded: true,
		}
		for _, child := range children {
			node.Children = append(node.Children, TreeNode{
				ID:    child,
				Label: child,
			})
		}
		tree.AddNode(node)
	}
	return tree
}
