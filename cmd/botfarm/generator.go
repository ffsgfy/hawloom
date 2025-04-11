package main

import (
	"slices"
	"strings"

	"github.com/mb-14/gomarkov"
)

type TextGenerator struct {
	chain *gomarkov.Chain
	ngram gomarkov.NGram
}

func NewTextGenerator(chain *gomarkov.Chain) *TextGenerator {
	res := TextGenerator{
		chain: chain,
		ngram: make(gomarkov.NGram, chain.Order),
	}
	res.Reset()
	return &res
}

func (g *TextGenerator) Clone() *TextGenerator {
	return &TextGenerator{
		chain: g.chain,
		ngram: slices.Clone(g.ngram),
	}
}

func (g *TextGenerator) Reset() *TextGenerator {
	for i := range g.ngram {
		g.ngram[i] = "\n"
	}
	return g
}

func (g *TextGenerator) Seed(source []string) {
	for i := range g.ngram {
		s := "\n"
		if i < len(source) {
			s = source[len(source)-i-1]
		}
		g.ngram[len(g.ngram)-i-1] = s
	}
}

func (g *TextGenerator) Generate() string {
	if next, err := g.chain.Generate(g.ngram); err == nil {
		for i := range len(g.ngram) - 1 {
			g.ngram[i] = g.ngram[i+1]
		}
		g.ngram[len(g.ngram)-1] = next
		return next
	}
	g.Reset()
	return "\n"
}

func (g *TextGenerator) GenerateN(n int, collapseWs bool) string {
	builder := strings.Builder{}
	lastWs := false

	for n > 0 {
		s := g.Generate()
		ws := s == "\n" || s == " "
		if !ws || !lastWs || !collapseWs {
			builder.WriteString(s)
		}
		if !ws {
			n--
		}
		lastWs = ws
	}

	return builder.String()
}

func NewMarkovChain(source string, order int) *gomarkov.Chain {
	chain := gomarkov.NewChain(order)
	chain.Add(strings.Split(source, ""))
	return chain
}
