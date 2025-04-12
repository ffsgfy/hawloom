package main

import (
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/mb-14/gomarkov"
)

type rngWrapper struct {
	rand.Rand
}

func NewRNG() *rngWrapper {
	return &rngWrapper{
		Rand: *rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
	}
}

func (rng *rngWrapper) Intn(n int) int {
	return rng.IntN(n)
}

type TextGenerator struct {
	chain *gomarkov.Chain
	ngram gomarkov.NGram
	rng   gomarkov.PRNG
}

func NewTextGenerator(chain *gomarkov.Chain) *TextGenerator {
	res := TextGenerator{
		chain: chain,
		ngram: make(gomarkov.NGram, chain.Order),
		rng:   NewRNG(),
	}
	res.Reset()
	return &res
}

func (g *TextGenerator) Clone() *TextGenerator {
	return &TextGenerator{
		chain: g.chain,
		ngram: slices.Clone(g.ngram),
		rng:   NewRNG(),
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
	if next, err := g.chain.GenerateDeterministic(g.ngram, g.rng); err == nil {
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
