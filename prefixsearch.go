package prefixsearch

import (
	"strings"
	"unicode"
)

type SearchTree struct {
	root *node
}

type node struct {
	value    interface{}
	childnum uint
	childs   map[rune]*node
}

func New() *SearchTree {
	return &SearchTree{
		root: &node{childs: map[rune]*node{}},
	}
}

func (tree *SearchTree) Add(key string, value interface{}) {
	current := tree.root

	needUpdate := (nil == tree.Search(key))

	for _, sym := range strings.ToLower(key) {
		if needUpdate {
			current.childnum++
		}
		next, ok := current.childs[sym]
		if !ok {
			newone := &node{childs: map[rune]*node{}}
			current.childs[sym] = newone
			next = newone
		}
		current = next
	}

	if needUpdate {
		current.childnum++
	}
	current.value = value
}

func (tree *SearchTree) AutoComplete(prefix string) []interface{} {
	current := tree.root
	for _, sym := range prefix {
		var ok bool
		current, ok = current.childs[unicode.ToLower(sym)]
		if !ok {
			return []interface{}{}
		}
	}

	result := make([]interface{}, 0, current.childnum)
	current.recurse(func(v interface{}) {
		if nil != v {
			result = append(result, v)
		}
	})
	return result
}

func (tree *SearchTree) Search(key string) interface{} {
	current := tree.root
	for _, sym := range key {
		var ok bool
		current, ok = current.childs[unicode.ToLower(sym)]
		if !ok {
			return nil
		}
	}
	return current.value
}

func (n *node) recurse(callback func(interface{})) {
	callback(n.value)
	for _, v := range n.childs {
		v.recurse(callback)
	}
}