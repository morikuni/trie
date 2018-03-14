package trie

import (
	"fmt"
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

	if l > lt {
		if b.fork == nil {
			b.fork = newFork()
		}
		b.fork.Add(oldSuffix)
	} else {
		if len(oldSuffix) != 0 {
			b.fork = newFork()
			b.fork.Add(oldSuffix)
			b.fork[oldSuffix[0]].fork = oldFork
		} else if b.fork == nil {
			b.fork = newFork()
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
	s := replacer.Replace(string(b.text))
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

var replacer = strings.NewReplacer(
	`\`, `\\`,
	`*`, `\*`,
	`+`, `\+`,
	`.`, `\.`,
	`?`, `\?`,
	`{`, `\{`,
	`}`, `\}`,
	`(`, `\(`,
	`)`, `\)`,
	`[`, `\[`,
	`]`, `\]`,
	`^`, `\^`,
	`$`, `\$`,
	`|`, `\|`,
)
