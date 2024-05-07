package avl_tree

import (
	"cmp"
)

type AVLTree[Key cmp.Ordered, Value any] struct {
	root *node[Key, Value]
}

func (t *AVLTree[Key, Value]) Add(key Key, value Value) {
	t.root = t.root.add(key, value)
}

func (t *AVLTree[Key, Value]) Remove(key Key) {
	t.root = t.root.remove(key)
}

func (t *AVLTree[Key, Value]) Search(key Key) (Value, bool) {

	if value := t.root.search(key); value == nil {
		return *new(Value), false
	} else {
		return value.value, true
	}
}

func (t *AVLTree[Key, Value]) GetRootValue() Value {
	return t.root.value
}

type node[Key cmp.Ordered, Value any] struct {
	key   Key
	value Value

	height int
	left   *node[Key, Value]
	right  *node[Key, Value]
}

func (n *node[Key, Value]) add(key Key, value Value) *node[Key, Value] {
	if n == nil {
		return &node[Key, Value]{key, value, 1, nil, nil}
	}

	if key < n.key {
		n.left = n.left.add(key, value)
	} else if key > n.key {
		n.right = n.right.add(key, value)
	} else {
		n.value = value
	}
	return n.rebalanceTree()
}

func (n *node[Key, Value]) remove(key Key) *node[Key, Value] {
	if n == nil {
		return nil
	}

	if key < n.key {
		n.left = n.left.remove(key)
	} else if key > n.key {
		n.right = n.right.remove(key)
	} else {
		if n.left != nil && n.right != nil {
			rightMinNode := n.right.findMinimum()
			n.key = rightMinNode.key
			n.value = rightMinNode.value
			n.right = n.right.remove(rightMinNode.key)
		} else if n.left != nil {
			n = n.left
		} else if n.right != nil {
			n = n.right
		} else {
			n = nil
		}

		return n
	}

	return n.rebalanceTree()
}

func (n *node[Key, Value]) search(key Key) *node[Key, Value] {
	if n == nil {
		return nil
	}

	if key < n.key {
		return n.left.search(key)
	} else if key > n.key {
		return n.right.search(key)
	}

	return n
}

func (n *node[Key, Value]) getHeight() int {
	if n == nil {
		return 0
	}

	return n.height
}

func (n *node[Key, Value]) recalculateHeight() {
	n.height = 1 + max(n.left.getHeight(), n.right.getHeight())
}

func (n *node[Key, Value]) rebalanceTree() *node[Key, Value] {
	if n == nil {
		return n
	}

	n.recalculateHeight()

	balanceFactor := n.left.getHeight() - n.right.getHeight()
	if balanceFactor == -2 {
		if n.right.left.getHeight() > n.right.right.getHeight() {
			n.right = n.right.rotateRight()
		}

		return n.rotateLeft()
	} else if balanceFactor == 2 {
		if n.left.right.getHeight() > n.left.left.getHeight() {
			n.left = n.left.rotateLeft()
		}

		return n.rotateRight()
	}

	return n
}

func (n *node[Key, Value]) rotateLeft() *node[Key, Value] {
	newRoot := n.right
	n.right = newRoot.left
	newRoot.left = n

	n.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (n *node[Key, Value]) rotateRight() *node[Key, Value] {
	newRoot := n.left
	n.left = newRoot.right
	newRoot.right = n

	n.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (n *node[Key, Value]) findMinimum() *node[Key, Value] {
	if n.left != nil {
		return n.left.findMinimum()
	}

	return n
}
