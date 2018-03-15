package trie

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Node interface {
	Add(text []rune)
	Regexp() string
}

func NewNode() Node {
	return fork(make(map[rune]*branch))
}

type branch struct {
	text []rune
	fork fork
}

func newBranch(text []rune, fork fork) *branch {
	return &branch{
		text,
		fork,
	}
}

func (b *branch) Add(text []rune) {
	l, lt := len(b.text), len(text)
	if l > lt {
		l = lt
	}
	var i int
	for ; i < l; i++ {
		if b.text[i] != text[i] {
			break
		}
	}
	commonPrefix := text[:i]
	oldSuffix := b.text[i:]
	oldFork := b.fork
	newSuffix := text[i:]
	b.text = commonPrefix

	if i == l && i == lt && b.fork == nil {
		return
	}

	if len(oldSuffix) != 0 {
		forked := newBranch(oldSuffix, oldFork)
		b.fork = newFork()
		b.fork[forked.text[0]] = forked
	}
	if b.fork == nil {
		b.fork = newFork()
		if len(oldSuffix) == 0 {
			b.fork.Add(oldSuffix)
		}
	}
	b.fork.Add(newSuffix)

	return
}

func (b *branch) Regexp() string {
	if b.text == nil {
		return ""
	}
	s := regexp.QuoteMeta(string(b.text))
	if b.fork != nil {
		return s + b.fork.Regexp()
	}
	return s
}

func newFork() fork {
	return fork(make(map[rune]*branch))
}

type fork map[rune]*branch

func (f fork) Add(text []rune) {
	if len(text) == 0 {
		f[termination] = newBranch(nil, nil)
		return
	}
	head := text[0]
	b, ok := f[head]
	if ok {
		b.Add(text)
		return
	}
	f[head] = newBranch(text, nil)
	return
}

func (f fork) Regexp() string {
	if len(f) == 0 {
		return ""
	}
	childs := make([]string, 0, len(f))
	for _, b := range f {
		childs = append(childs, b.Regexp())
	}
	return fmt.Sprintf("(%s)", strings.Join(childs, "|"))
}

var termination rune
