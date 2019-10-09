package ahocorasick

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestNew(t *testing.T) {
	type node struct {
		path           []rune
		inDict         bool
		suffixLink     []rune
		dictSuffixLink []rune
	}
	tests := []struct {
		dict  []string
		nodes []node
	}{
		{
			nil,
			nil,
		},
		{
			[]string{"a", "ab", "bab", "bc", "bca", "c", "caa"},
			[]node{
				{[]rune{}, false, nil, nil},
				{[]rune{'a'}, true, []rune{}, nil},
				{[]rune{'a', 'b'}, true, []rune{'b'}, nil},
				{[]rune{'b'}, false, []rune{}, nil},
				{[]rune{'b', 'a'}, false, []rune{'a'}, []rune{'a'}},
				{[]rune{'b', 'a', 'b'}, true, []rune{'a', 'b'}, []rune{'a', 'b'}},
				{[]rune{'b', 'c'}, true, []rune{'c'}, []rune{'c'}},
				{[]rune{'b', 'c', 'a'}, true, []rune{'c', 'a'}, []rune{'a'}},
				{[]rune{'c'}, true, []rune{}, nil},
				{[]rune{'c', 'a'}, false, []rune{'a'}, []rune{'a'}},
				{[]rune{'c', 'a', 'a'}, true, []rune{'a'}, []rune{'a'}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.dict, ","), func(t *testing.T) {
			root := New(tt.dict).root
			if got, want := countNodes(root), len(tt.nodes); got != want {
				t.Errorf("got %d nodes in tree, want %d", got, want)
			}
			for _, want := range tt.nodes {
				t.Run(fmt.Sprintf("%q", want.path), func(t *testing.T) {
					n := findNode(root, want.path)
					if n == nil {
						t.Fatalf("got no node %q", want.path)
					}
					if gotInDict := n.entry != nil; gotInDict != want.inDict {
						t.Errorf("got node %q in dictionary = %t, want %t", want.path, gotInDict, want.inDict)
					}
					if got, want := n.suffixLink, findNode(root, want.suffixLink); got != want {
						t.Errorf("got suffix link node %q, want %q", nodePath(got), nodePath(want))
					}
					if got, want := n.dictSuffixLink, findNode(root, want.dictSuffixLink); got != want {
						t.Errorf("got directory suffix link node %q, want %q", nodePath(got), nodePath(want))
					}
				})
			}
			if t.Failed() {
				t.Log("================================================================================")
				t.Log("Debug tree:\n")
				buf := new(bytes.Buffer)
				fprintNode(buf, root, 0)
				t.Log(buf.String())
				t.Log("================================================================================")
			}
		})
	}
}

func TestTrie_Search(t *testing.T) {
	tests := []struct {
		dict []string
		text string
		want map[int][]int
	}{
		{
			nil,
			"search without dictionary",
			nil,
		},
		{
			[]string{"search", "without", "text"},
			"",
			nil,
		},
		{
			[]string{"a", "ab", "bab", "bc", "bca", "c", "caa"},
			"abccab",
			map[int][]int{
				0: {0, 4}, // Abccab and abccAb
				1: {0, 4}, // ABccab and abccAB
				3: {1},    // aBCcab
				5: {2, 3}, // abCcab and abcCab
			},
		},
		{
			[]string{"а", "аб", "баб", "бв", "бва", "в", "ваа"},
			"абвваб",
			map[int][]int{
				0: {0, len("абвв")},        // Абвваб and абввАб
				1: {0, len("абвв")},        // АБвваб and абввАБ
				3: {len("а")},              // аБВваб
				5: {len("аб"), len("абв")}, // абВваб and абвВаб
			},
		},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.dict, ","), func(t *testing.T) {
			trie := New(tt.dict)
			if got := trie.Search(tt.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Trie.Search(%q) = %v, want %v", tt.text, got, tt.want)

				buf := new(bytes.Buffer)
				t.Log("================================================================================")
				t.Logf("Debug text: %q (%d bytes)", tt.text, len(tt.text))
				fprintText(buf, tt.text)
				t.Log("\n" + buf.String())

				buf.Reset()
				t.Log("================================================================================")
				t.Log("Debug trie:\n")
				fprintNode(buf, trie.root, 0)
				t.Log(buf.String())
				t.Log("================================================================================")
			}
		})
	}
}

func countNodes(n *node) int {
	var c int
	if n == nil {
		return c
	}
	c++ // +1 for n
	for _, n = range n.children {
		c += countNodes(n)
	}
	return c
}

func findNode(root *node, path []rune) *node {
	if path == nil {
		return nil
	}
	n := root
	for _, r := range path {
		c, ok := n.children[r]
		if !ok {
			return nil
		}
		n = c
	}
	return n
}

// nodePath returns the path to n.
//
// It traverses the tree branch from n to the root to determine the steps and
// returns the steps in reversed order.
func nodePath(n *node) []rune {
	var path []rune
	if n == nil {
		return path
	}
	for ; n.parent != nil; n = n.parent {
		path = append(path, n.rune)
	}
	for left, right := 0, len(path)-1; left < right; left, right = left+1, right-1 {
		path[left], path[right] = path[right], path[left]
	}
	return path
}

func fprintNode(w io.Writer, n *node, indent int) {
	if n == nil {
		fmt.Fprintln(w, "<nil>")
		return
	}
	printEntry := func(n *node) string {
		if n.entry == nil {
			return "nil"
		}
		return fmt.Sprintf("index = %d, len = %d", n.entry.index, n.entry.len)
	}
	out := []string{
		fmt.Sprintf("Path: %q %p", nodePath(n), n),
		fmt.Sprintf("Dictionaty entry: %s", printEntry(n)),
		fmt.Sprintf("Suffix link: %q %p", nodePath(n.suffixLink), n.suffixLink),
		fmt.Sprintf("Dict suffix link: %q %p", nodePath(n.dictSuffixLink), n.dictSuffixLink),
		fmt.Sprintf("Children count: %d", len(n.children)),
	}
	for i, _ := range out {
		out[i] = strings.Repeat(".", indent*4) + out[i]
	}
	fmt.Fprintln(w, strings.Join(out, "\n"))

	for _, n := range n.children {
		fprintNode(w, n, indent+1)
	}
}

func fprintText(w io.Writer, text string) {
	for i, r := range text {
		fmt.Fprintf(w, "%3d %q (%d bytes)\n", i, r, utf8.RuneLen(r))
	}
}
