package parse

import (
	"github.com/romanoff/servicebuilder/app"
	"io"
)

func NewParser(name string, reader io.Reader) *Parser {
	parser := &Parser{lexer: NewLexer(name, reader), buffer: make([]Lexeme, 0, 0)}
	parser.scanner = parser.lexer.Scan()
	return parser
}

type Parser struct {
	lexer   *Lexer
	scanner chan Lexeme
	buffer  []Lexeme
	index   int
}

func (self *Parser) scan() *Lexeme {
	if self.index == len(self.buffer) {
		lexeme := <-self.scanner
		self.buffer = append(self.buffer, lexeme)
		self.index++
		return &lexeme
	}
	lexeme := self.buffer[self.index]
	self.index++
	return &lexeme
}

func (self *Parser) scanIgnoreWhitespace() *Lexeme {
	lexeme := self.scan()
	if lexeme.Token == WS {
		lexeme = self.scan()
	}
	return lexeme
}

func (self *Parser) unscan() *Lexeme {
	if self.index == 0 {
		return nil
	}
	self.index--
	lexeme := self.buffer[self.index]
	return &lexeme
}
func (self *Parser) Parse() (app.Application, error) {
	return app.Application{}, nil
}
