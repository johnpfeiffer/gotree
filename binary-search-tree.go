package main

import (
	"bytes"
	"fmt"
	"strconv"
)

// Node contains data (and usually a value or a pointer to a value) and pointers to the child nodes
type Node struct {
	left  *Node
	right *Node
	Data  int
}

// BinarySearchTree https://en.wikipedia.org/wiki/Binary_search_tree
type BinarySearchTree struct {
	Root *Node
}

// Display returns a string with the data from a pre-order traversal of all the nodes in the tree
func (tree *BinarySearchTree) Display() string {
	return TraverseInOrder(tree.Root)
	// return TraversePreOrder(tree.Root)
}

// TraversePreOrder shows the node data (in pre-order) and continues recursively https://en.wikipedia.org/wiki/Tree_traversal#Pre-order
func TraversePreOrder(n *Node) string {
	var s string
	var b bytes.Buffer
	if n == nil {
		return ""
	}
	b.WriteString(fmt.Sprintf("%d ", n.Data))
	if n.left != nil {
		s += TraversePreOrder(n.left)
	}
	if n.right != nil {
		s += TraversePreOrder(n.right)
	}
	return b.String() + s
}

// TraverseInOrder shows the node data (in-order) and continues recursively, in a BST this ouputs the data in sorted order
func TraverseInOrder(n *Node) string {
	var s string
	if n == nil {
		return ""
	}
	if n.left != nil {
		s = s + TraverseInOrder(n.left)
	}
	s = s + " " + strconv.Itoa(n.Data)
	if n.right != nil {
		s = s + TraverseInOrder(n.right)
	}
	return s
}

// TODO: post-order display

// Height as defined by https://en.wikipedia.org/wiki/Binary_tree
func (tree *BinarySearchTree) Height() int {
	if tree.Root.left == nil && tree.Root.right == nil {
		return 0
	}
	return subtreeHeight(tree.Root) - 1
}

func subtreeHeight(n *Node) int {
	leftMax := 0
	rightMax := 0
	if n.left == nil && n.right == nil {
		return 1
	}
	if n.left != nil {
		leftMax = subtreeHeight(n.left)
	}
	if n.right != nil {
		rightMax = subtreeHeight(n.right)
	}
	if leftMax > rightMax {
		return leftMax + 1
	}
	return rightMax + 1
}

// Find returns the first node that has a matching key
func (tree *BinarySearchTree) Find(target int) *Node {
	current := tree.Root
	for {
		switch {
		case current == nil:
			return nil
		case current.Data == target:
			return current
		case current.Data > target:
			current = current.left
		case current.Data < target:
			current = current.right
		}
	}
}

// FindParent returns the parent of the first matching node from the subtree
func FindParent(target int, start *Node) *Node {
	current := start
	parent := current
	for {
		switch {
		case current == nil:
			return nil
		case current.Data == target:
			return parent
		case current.Data < target:
			parent = current
			current = current.right
		case current.Data > target:
			parent = current
			current = current.left
		}
	}
}

// InsertValue adds data (with a new node) to the Binary Search Tree
func (tree *BinarySearchTree) InsertValue(target int) {
	if tree.Root == nil {
		tree.Root = &Node{Data: target}
		return
	}
	current := tree.Root
	for {
		if current.Data > target {
			if current.left == nil {
				current.left = &Node{Data: target}
				return
			}
			current = current.left
		} else {
			if current.right == nil {
				current.right = &Node{Data: target}
				return
			}
			current = current.right
		}
	}
}

// RemoveValue removes the first node with the matching data
func (tree *BinarySearchTree) RemoveValue(target int) {
	if tree.Root == nil {
		return
	}
	if tree.Root.Data == target {
		RemoveRoot(tree)
		return
	}
	parent := FindParent(target, tree.Root)
	if parent == nil {
		return
	}
	if parent.left != nil && parent.left.Data == target {
		RemoveNode(parent.left, parent)
		return
	}
	if parent.right != nil && parent.right.Data == target {
		RemoveNode(parent.right, parent)
	}
}

// FindRightMostParent returns the parent of the right most node in a (sub)tree
func FindRightMostParent(subtreeRoot *Node) *Node {
	// TODO: defensive handling of passing in nil
	parent := subtreeRoot
	for current := subtreeRoot; current.right != nil; current = current.right {
		parent = current
	}
	return parent
}

// FindLeftMostParent returns the parent of the left most node in a (sub)tree
func FindLeftMostParent(subtreeRoot *Node) *Node {
	// TODO: defensive handling of passing in nil
	parent := subtreeRoot
	for current := subtreeRoot; current.left != nil; current = current.left {
		parent = current
	}
	return parent
}

// RemoveRoot handles the special case of removing the root node
func RemoveRoot(tree *BinarySearchTree) {
	if tree.Root == nil {
		return
	}
	switch {
	case tree.Root.right == nil && tree.Root.left == nil:
		tree.Root = nil // if pointers then tree.Root.Data = nil to prevent memory leaks
	case tree.Root.right != nil:
		if tree.Root.right.left == nil { // simplest hoisting case, e.g. root to leaf: 1 2(root) 5 6 becomes 1 5(root) 6
			originalRootLeft := tree.Root.left
			tree.Root = tree.Root.right
			tree.Root.left = originalRootLeft
		} else { // more complex hoisting case, e.g. 1 2(root) 5 4 3 , need to find the right subtree left most, then fix the right subtree
			parent := tree.Root.right
			for current := tree.Root.right; current.left != nil; current = current.left {
				parent = current
			}
			originalRootLeft := tree.Root.left
			originalRootRight := tree.Root.right
			replacementRoot := parent.left

			tree.Root = replacementRoot
			parent.left = replacementRoot.right
			replacementRoot.left = originalRootLeft
			replacementRoot.right = originalRootRight
		}
		return

	case tree.Root.left != nil:
		if tree.Root.left.right == nil { // simplest hoisting case, e.g. root to leaf: 0 1 2(root) 5 becomes 0 1(root) 5
			originalRootRight := tree.Root.right
			tree.Root = tree.Root.left
			tree.Root.right = originalRootRight
		} else { // more complex hoisting case, e.g. 1 0 2(root) 5 , need to find the left subtree right most, then fix the left subtree
			parent := tree.Root.left
			for current := tree.Root.left; current.right != nil; current = current.right {
				parent = current
			}
			originalRootLeft := tree.Root.left
			originalRootRight := tree.Root.right
			replacementRoot := parent.right

			tree.Root = replacementRoot
			parent.right = replacementRoot.left
			replacementRoot.left = originalRootLeft
			replacementRoot.right = originalRootRight
		}
		return

	default:
		fmt.Println("ERROR should never reach here")
	}
}

// replaceLeft is used to remove a node (when the right most node is found in a left subtree)
// case 1 (leaf): removeeParent is 0, removee is -1, replacementParent is -1, replacement is nil
// case 2: removeeParent is 2, removee is 0, replacementParent is 0, replacement is 1

/*
	     8
	   4
     1  6
   -1  3
	 0

LEAF (remove -1)  NO RIGHT (remove 1)   RIGHT (remove -2)    BOTH
  2   2   2      2   2     2              2      2     2       2       2       2       2
-1   0   0      1   1     1             -2     -2    -2      -2      -2      -2      -2
   -1  -1 1    0   0     0                0      0     0    -6 0   -6  0   -6  0   -6   0
                   -1  -2 -1                   -1    -1 1         -7        -3      -4
                                                                                      -3
*/
func removeLeft(removeeParent *Node) {
	switch {
	case removeeParent == nil:
		fmt.Println("ERROR should never reach here with removeeParent as nil")
		return
	case removeeParent.left == nil:
		return
	}
	removee := removeeParent.left
	if removee.right == nil && removee.left == nil {
		removeeParent.left = nil // if pointers then removee.Data = nil to prevent memory leaks
		return
	}
	// CONTINUE, there is a replacement node that must be in the left subtree
	if removee.right == nil { // the easy case, just hoist the left child up
		removeeParent.left = removee.left
		return
	}
	// CONTINUE, there is a replacement that must be in the right subtree
	if removee.left == nil { // the easy case, just hoist the right child up
		removeeParent.left = removee.right
		return
	}
	// CONTINUE, the replacement is either the leftmost of the right subtree or the rightmost of the left subtree
	// TODO: should I measure the depth of the left and right subtrees to push towards balance?
	// ARBITRARILY choosing the left subtree to pick a replacement (the right-most)
	replacementParent := FindRightMostParent(removee.left)
	removeeRight := removee.right
	// edge case where the replacement is the left child of the removee
	if replacementParent == removee.left { // the easy case, no competing right children so just hoist the left child up
		removeeParent.left = replacementParent
		replacementParent.right = removeeRight
		return
	}
	// CONTINUE, the replacementParent.right is different than the removee.left and is the right most in the subtree
	replacement := replacementParent.right
	fmt.Println("JOHN", replacementParent.Data, replacement.Data)
	if replacement.left == nil { // no merge required yet
		removeeParent.left = replacement
		replacementParent.right = removeeRight
		return
	}
	// TODO: more complex variations

}

// RemoveNode also handles the special case of removing a leaf node
func RemoveNode(node, parent *Node) {
	if node == nil || parent == nil {
		fmt.Println("ERROR should never reach here with node or parent as nil")
		return
	}
	switch {
	case parent.left == node:
		if node.left != nil {
			// fmt.Println("j", node.left.Data)
			parent.left = node.left // hoist the remaining child, it is ok if we re-assign nil
			node.left = nil         // reminder that for pointers node.Data = nil prevents memory leaks
		} else {
			// fmt.Println("JJ", node.right.Data)
			parent.left = node.right // simple logic as we do not care if we re-assign nil
			node.right = nil
		}
	case parent.right == node:
		if node.left != nil {
			parent.right = node.left
			node.left = nil
		} else {
			parent.right = node.right
			node.right = nil
		}
	default:
		fmt.Println("ERROR should never reach here with parent not matching the child node")
	}
}
