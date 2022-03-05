package utils

import "strings"

type TrieTree struct {
	root *node
}

type node struct {
	isRoot   bool
	isLeaf   bool
	char     rune
	children map[rune]*node
}

func NewTrieTree() *TrieTree {
	return &TrieTree{
		root: newNode(0, true),
	}
}

func newNode(char rune, isRoot bool) *node {
	return &node{
		isRoot:   isRoot,
		char:     char,
		children: make(map[rune]*node),
	}
}

func (tree *TrieTree) AddWord(words ...string) {
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word != "" {
			tree.addWord(word)
		}
	}
}

func (tree *TrieTree) addWord(word string) {
	var current = tree.root
	var runes = []rune(word)
	for idx := 0; idx < len(runes); idx++ {
		r := runes[idx]
		if next, ok := current.children[r]; ok {
			current = next
		} else {
			child := newNode(r, false)
			current.children[r] = child
			current = child
		}
		if idx == len(runes)-1 {
			current.isLeaf = true
		}
	}
}

// Filter: replace the longest matched runes with mask
func (tree *TrieTree) Filter(text string, mask rune) string {
	runes := []rune(text)
	length := len(runes)
	current := tree.root
	leftMatched := 0
	found := false
	for idx := 0; idx < length; idx++ {
		current, found = current.children[runes[idx]]

		if !found || (!current.isLeaf && idx == length-1) {
			current = tree.root
			idx = leftMatched
			leftMatched++
			continue
		}

		if current.isLeaf && leftMatched <= idx {
			for i := leftMatched; i <= idx; i++ {
				runes[i] = mask
			}
		}
	}

	return string(runes)
}
