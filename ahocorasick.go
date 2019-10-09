// Package ahocorasick provides string-searching using Aho-Corasick algorithm.
package ahocorasick

import "unicode/utf8"

// New returns new search tree for dictionary of rune slices.
func New(dictionary []string) Trie {
	// Building the trie is a 3-step process:
	// 1) Create and link tree nodes with parent-child relationships.
	// 2) Update each node's longest possible strict suffix.
	// 3) Update each node's directory suffix.

	if len(dictionary) == 0 {
		return Trie{}
	}

	t := Trie{root: new(node)}

	// Cache references of the nodes in each level.
	levelNodes := make(
		[]map[*node]struct{},
		func() int {
			max := 0
			for _, d := range dictionary {
				c := len(d)
				if c > max {
					max = c
				}
			}
			return max
		}(),
	)

	// Step 1) Create and link tree nodes with parent-child relationships.
	for idx, ent := range dictionary {
		n := t.root
		length := len(ent)
		for i, r := range ent {
			c, ok := n.children[r]
			if !ok {
				if n.children == nil {
					n.children = make(map[rune]*node)
				}
				c = &node{rune: r, parent: n}
				n.children[r] = c
			}

			if i+utf8.RuneLen(r) == length {
				// r is the last rune. This makes c a dictionary node.
				c.entry = &entry{index: idx, len: length}
			}

			if levelNodes[i] == nil {
				levelNodes[i] = make(map[*node]struct{})
			}
			levelNodes[i][c] = struct{}{}

			// The child becomes parent in the next iteration, so that the next
			// ent rune is associated with a child of c.
			n = c
		}
	}

	// Step 2) Update each node's longest possible strict suffix.
	//
	// To find the longest possible strict suffix of a node, traverse the longest
	// possible strict suffixes of the parent node, until the a node with child
	// for the same rune is found.
	//
	// Start from the first level after the root and traverse the tree one level
	// at the time, so that each node's parent has its longest possible strict
	// suffix link updated.
	//
	// This can be performed only after Step 1), since we need all nodes.
	for _, nodes := range levelNodes {
		for n := range nodes {
			if n.parent == t.root {
				n.suffixLink = t.root
			}
			for r, p := n.rune, n.parent; p.suffixLink != nil; p = p.suffixLink {
				if l, ok := p.suffixLink.children[r]; ok {
					n.suffixLink = l
					break
				}
			}
		}
	}

	// Step 3) Update each node's directory suffix.
	//
	// To find the dictionary suffix link of a node, traverse the longest possible
	// strict suffixes until a dictionary node is found.
	//
	// This can be performed only after Step 1) and Step 2), since we need all
	// nodes with their longest possible strict suffix link.
	for _, nodes := range levelNodes {
		for n := range nodes {
			for p := n.suffixLink; p != nil; p = p.suffixLink {
				if p.entry != nil {
					n.dictSuffixLink = p
					break
				}
			}
		}
	}

	return t
}

// Trie provides interface to the search tree.
type Trie struct {
	root *node
}

// Search returns the indexes of the dictionary entries and the positions of their
// occurrences in text.
func (t Trie) Search(text string) map[int][]int {
	if t.root == nil || len(text) == 0 {
		return nil
	}
	m := map[int][]int{}
	n := t.root
	for i, r := range text {
		n = n.find(r)
		if n == nil {
			n = t.root
		}
		for _, e := range n.entries() {
			if _, ok := m[e.index]; !ok {
				m[e.index] = []int{}
			}
			m[e.index] = append(m[e.index], i+utf8.RuneLen(r)-e.len)
		}
	}
	return m
}

// entry contains properties of dictionary entry.
type entry struct {
	// The index of the dictionary entry.
	index int

	// The number of bytes of the dictionary entry.
	len int
}

// node represents a node corresponding to a rune in the tree.
type node struct {
	// The rune to which this node is mapped to.
	rune rune

	// Dictionary entry details if the node is a dictionary node.
	entry *entry

	// Pointer to the parent node.
	parent *node

	// Pointers to the child nodes.
	children map[rune]*node

	// Pointer to the longest strict suffix of the node in the graph.
	//
	// The node can be computed in linear time by repeatedly traversing the
	// strict suffix nodes of a node's parent until the traversing node has a child
	// matching the character of the target node.
	suffixLink *node

	// Pointer to the dictionary suffix node.
	//
	// The dictionary suffix can be computed in linear time by repeatedly
	// traversing the longest strict suffix nodes until a dictionary node is found.
	dictSuffixLink *node
}

// entries returns the entries.
func (n *node) entries() []*entry {
	var out []*entry
	if n.entry != nil {
		out = append(out, n.entry)
	}
	for p := n.dictSuffixLink; p != nil; p = p.dictSuffixLink {
		out = append(out, p.entries()...)
	}
	return out
}

// find finds and returns the child corresponding to r and if that does not
// exist, tries to find it in the suffix's children, and if that does not exist,
// tries to find it in the suffix's suffix's children, and so on, until the root
// is reached or a matching child node is found.
func (n *node) find(r rune) *node {
	if c, ok := n.children[r]; ok {
		return c
	}
	if n.suffixLink != nil {
		return n.suffixLink.find(r)
	}
	return nil
}
