package lexer_test

import (
	"testing"

	require "github.com/alecthomas/assert/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type staticLexer struct {
	tokens []lexer.Token
}

func (s *staticLexer) Next() (lexer.Token, error) {
	if len(s.tokens) == 0 {
		return lexer.EOFToken(lexer.Position{}), nil
	}
	t := s.tokens[0]
	s.tokens = s.tokens[1:]
	return t, nil
}

func TestUpgrade(t *testing.T) {
	t0 := lexer.Token{Type: 1, Value: "moo"}
	ts := lexer.Token{Type: 3, Value: " "}
	t1 := lexer.Token{Type: 2, Value: "blah"}
	tokens := []lexer.Token{t0, ts, t1}
	l, err := lexer.Upgrade(&staticLexer{tokens: tokens}, 3)
	require.NoError(t, err)
	require.Equal(t, t0, l.Peek())
	require.Equal(t, t0, l.Peek())
	require.Equal(t, tokens, l.Range(0, 3))
}

func TestPeekAndNextAny(t *testing.T) {
	slexdef := lexer.MustSimple([]lexer.SimpleRule{
		{"Ident", `\w+`},
		{"Whitespace", `\s+`},
	})
	slex, err := slexdef.LexString("", `hello world last`)
	require.NoError(t, err)
	plex, err := lexer.Upgrade(slex, slexdef.Symbols()["Whitespace"])
	require.NoError(t, err)
	expected := []lexer.Token{
		{Type: -2, Value: "hello", Pos: lexer.Position{Line: 1, Column: 1}},
		{Type: -3, Value: " ", Pos: lexer.Position{Line: 1, Column: 6, Offset: 5}},
		{Type: -2, Value: "world", Pos: lexer.Position{Line: 1, Column: 7, Offset: 6}},
		{Type: -3, Value: " ", Pos: lexer.Position{Line: 1, Column: 12, Offset: 11}},
		{Type: -2, Value: "last", Pos: lexer.Position{Line: 1, Column: 13, Offset: 12}},
	}
	tok := plex.Next()
	require.Equal(t, expected[0], tok)
	require.Equal(t, expected[2], plex.Peek(), "should have skipped whitespace")
}
