package parse

import (
	"bufio"
	"bytes"
	"io"
)

type Lexeme struct {
	Token Token
	Value string
	Pos   Position
}

// Position of lexeme
type Position struct {
	Line, Column int
}

// Lexical parser for servicebuilder
type Lexer struct {
	Name     string        // Input filename. For error messages
	Reader   *bufio.Reader //  Content reader
	bus      chan Lexeme   // Lexemes bus is populated with lexemes as they are consumed
	Position *Position
}

// Return new lexer
func NewLexer(name string, reader io.Reader) *Lexer {
	return &Lexer{
		Name:     name,
		Reader:   bufio.NewReader(reader),
		bus:      make(chan Lexeme),
		Position: &Position{0, 0},
	}
}

// Read next rune
func (self *Lexer) read() rune {
	ch, _, err := self.Reader.ReadRune()
	if self.Position != nil {
		if ch == '\n' {
			self.Position.Line++
			self.Position.Column = 0
		} else {
			self.Position.Column++
		}
	}
	if err != nil {
		return eof
	}
	return ch
}

// Go to the previous rune
func (self *Lexer) unread() {
	_ = self.Reader.UnreadRune()
	if self.Position != nil {
		if self.Position.Column == 0 {
			if self.Position.Line > 0 {
				self.Position.Line--
			}
		} else {
			self.Position.Column--
		}
	}
}

func (self *Lexer) Scan() chan Lexeme {
	go func() {
		self.scan()
	}()
	return self.bus
}

var eof = rune(0)

func (self *Lexer) scanWhitespace() Lexeme {
	var buf bytes.Buffer
	buf.WriteRune(self.read())
	position := *self.Position
	for {
		if ch := self.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			self.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return Lexeme{WS, buf.String(), position}
}

func (self *Lexer) scanIdentifier() Lexeme {
	var buf bytes.Buffer
	buf.WriteRune(self.read())
	position := *self.Position
	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := self.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			self.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	return Lexeme{IDENT, buf.String(), position}
}

func (self *Lexer) scan() {
ScanLoop:
	for {
		ch := self.read()
		if isWhitespace(ch) {
			self.unread()
			self.bus <- self.scanWhitespace()
			continue
		} else if isLetter(ch) {
			self.unread()
			self.bus <- self.scanIdentifier()
			continue
		}
		switch ch {
		case eof:
			self.bus <- Lexeme{EOF, "", *self.Position}
			break ScanLoop
		case '{':
			self.bus <- Lexeme{LEFTBRACE, string(ch), *self.Position}
			continue
		case '}':
			self.bus <- Lexeme{RIGHTBRACE, string(ch), *self.Position}
			continue
		}
		self.bus <- Lexeme{ILLEGAL, string(ch), *self.Position}
	}
}
