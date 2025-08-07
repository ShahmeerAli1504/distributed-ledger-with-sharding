package core

import (
	"fmt"
)

// Color represents the color of a Red-Black Tree node
type Color bool

const (
	Red   Color = true
	Black Color = false
)

// RBNode represents a node in the Red-Black Tree
type RBNode struct {
	Shard  *Shard  // Shard metadata
	Color  Color   // Red or Black
	Left   *RBNode // Left child (corrected from RBpfNode)
	Right  *RBNode // Right child
	Parent *RBNode // Parent node
}

// RBTree represents a Red-Black Tree for shard indexing
type RBTree struct {
	Nil  *RBNode // Sentinel node for nil leaves
	Root *RBNode // Root of the tree
}

// NewRBTree creates a new Red-Black Tree
func NewRBTree() *RBTree {
	nilNode := &RBNode{Color: Black}
	return &RBTree{
		Nil:  nilNode,
		Root: nilNode,
	}
}

// Insert adds a shard to the Red-Black Tree
func (t *RBTree) Insert(shard *Shard) {
	node := &RBNode{
		Shard:  shard,
		Color:  Red,
		Left:   t.Nil,
		Right:  t.Nil,
		Parent: t.Nil,
	}

	// Standard BST insertion
	current := t.Root
	parent := t.Nil
	for current != t.Nil {
		parent = current
		if shard.ID < current.Shard.ID {
			current = current.Left
		} else {
			current = current.Right
		}
	}

	node.Parent = parent
	if parent == t.Nil {
		t.Root = node
	} else if shard.ID < parent.Shard.ID {
		parent.Left = node
	} else {
		parent.Right = node
	}

	// Fix Red-Black Tree properties
	t.fixInsert(node)
}

// fixInsert balances the tree after insertion
func (t *RBTree) fixInsert(node *RBNode) {
	for node.Parent.Color == Red {
		if node.Parent == node.Parent.Parent.Left {
			uncle := node.Parent.Parent.Right
			if uncle.Color == Red {
				node.Parent.Color = Black
				uncle.Color = Black
				node.Parent.Parent.Color = Red
				node = node.Parent.Parent
			} else {
				if node == node.Parent.Right {
					node = node.Parent
					t.leftRotate(node)
				}
				node.Parent.Color = Black
				node.Parent.Parent.Color = Red
				t.rightRotate(node.Parent.Parent)
			}
		} else {
			uncle := node.Parent.Parent.Left
			if uncle.Color == Red {
				node.Parent.Color = Black
				uncle.Color = Black
				node.Parent.Parent.Color = Red
				node = node.Parent.Parent
			} else {
				if node == node.Parent.Left {
					node = node.Parent
					t.rightRotate(node)
				}
				node.Parent.Color = Black
				node.Parent.Parent.Color = Red
				t.leftRotate(node.Parent.Parent)
			}
		}
		if node == t.Root {
			break
		}
	}
	t.Root.Color = Black
}

// leftRotate performs a left rotation on a node
func (t *RBTree) leftRotate(x *RBNode) {
	y := x.Right
	x.Right = y.Left
	if y.Left != t.Nil {
		y.Left.Parent = x
	}
	y.Parent = x.Parent
	if x.Parent == t.Nil {
		t.Root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}
	y.Left = x
	x.Parent = y
}

// rightRotate performs a right rotation on a node
func (t *RBTree) rightRotate(x *RBNode) {
	y := x.Left
	x.Left = y.Right
	if y.Right != t.Nil {
		y.Right.Parent = x
	}
	y.Parent = x.Parent
	if x.Parent == t.Nil {
		t.Root = y
	} else if x == x.Parent.Right {
		x.Parent.Right = y
	} else {
		x.Parent.Left = y
	}
	y.Right = x
	x.Parent = y
}

// FindShard retrieves a shard by ID in O(log n) time
func (t *RBTree) FindShard(id int) (*Shard, bool) {
	current := t.Root
	for current != t.Nil {
		if id == current.Shard.ID {
			return current.Shard, true
		} else if id < current.Shard.ID {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	return nil, false
}

// GetAllShards collects all shards in the tree
func (t *RBTree) GetAllShards() []*Shard {
	var shards []*Shard
	t.inOrderTraversal(t.Root, &shards)
	return shards
}

// inOrderTraversal collects shards in-order
func (t *RBTree) inOrderTraversal(node *RBNode, shards *[]*Shard) {
	if node != t.Nil {
		t.inOrderTraversal(node.Left, shards)
		*shards = append(*shards, node.Shard)
		t.inOrderTraversal(node.Right, shards)
	}
}

// PrintTree displays the tree structure (for debugging)
func (t *RBTree) PrintTree() {
	fmt.Println("\n--- Red-Black Tree Shard Index ---")
	t.printNode(t.Root, 0)
}

func (t *RBTree) printNode(node *RBNode, level int) {
	if node == t.Nil {
		return
	}
	prefix := ""
	for i := 0; i < level; i++ {
		prefix += "  "
	}
	color := "Black"
	if node.Color == Red {
		color = "Red"
	}
	fmt.Printf("%sShard #%d (Color: %s, Blocks: %d, Merkle Root: %s)\n", prefix, node.Shard.ID, color, len(node.Shard.Blocks), node.Shard.GetRoot())
	t.printNode(node.Left, level+1)
	t.printNode(node.Right, level+1)
}
