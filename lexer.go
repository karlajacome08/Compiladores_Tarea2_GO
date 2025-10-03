package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

type Token int

const (
	EOF Token = iota
	ILLEGAL
	IDENT
	INT
	SEMI // ;

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	ASSIGN // =
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	INT:     "INT",
	SEMI:    ";",
	ADD:     "+",
	SUB:     "-",
	MUL:     "*",
	DIV:     "/",
	ASSIGN:  "=",
}

func (t Token) String() string { return tokens[int(t)] }

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(r),
	}
}

func (l *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}
			panic(err)
		}

		l.pos.column++

		switch r {
		case '\n':
			l.resetPosition()
		case ';':
			return l.pos, SEMI, ";"
		case '+':
			return l.pos, ADD, "+"
		case '-':
			return l.pos, SUB, "-"
		case '*':
			return l.pos, MUL, "*"
		case '/':
			return l.pos, DIV, "/"
		case '=':
			return l.pos, ASSIGN, "="
		default:
			if unicode.IsSpace(r) {
				continue
			}
			if unicode.IsDigit(r) {
				startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return startPos, INT, lit
			}
			if unicode.IsLetter(r) || r == '_' {
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()
				return startPos, IDENT, lit
			}
			return l.pos, ILLEGAL, string(r)
		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}
	l.pos.column--
}

func (l *Lexer) lexInt() string {
	var lit []rune
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return string(lit)
			}
			panic(err)
		}
		l.pos.column++
		if unicode.IsDigit(r) {
			lit = append(lit, r)
			continue
		}
		l.backup()
		return string(lit)
	}
}

func (l *Lexer) lexIdent() string {
	var lit []rune
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return string(lit)
			}
			panic(err)
		}
		l.pos.column++
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			lit = append(lit, r)
			continue
		}
		l.backup()
		return string(lit)
	}
}

func main() {
	f, err := os.Open("input.test")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lexer := NewLexer(f)
	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}
		fmt.Printf("%d:%d\t%s\t%q\n", pos.line, pos.column, tok, lit)
	}
}
